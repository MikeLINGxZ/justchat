//go:build !windows

package plugin

import (
	"context"
	"fmt"
	"os"
	"strings"

	pkgcli "gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/cli"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/ierror"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/service/plugin/plugin_dto"
)

// loginSessionIface is the narrow interface that plugin_login methods use to interact
// with a running login session. *pkgcli.LoginSession satisfies this interface on
// non-Windows platforms; fakeLoginSession satisfies it in tests.
type loginSessionIface interface {
	Output() <-chan []byte
	Write(data []byte) error
	Resize(rows, cols uint16) error
	Cancel() error
	Wait() (int, error)
}

// loginOutputPayload is emitted on the "cli.login.output" event for each PTY chunk.
type loginOutputPayload struct {
	ID   string `json:"id"`
	Data string `json:"data"`
}

// loginDonePayload is emitted on the "cli.login.done" event when the session exits.
type loginDonePayload struct {
	ID       string `json:"id"`
	ExitCode int    `json:"exit_code"`
	Error    string `json:"error"`
}

// LoginCli starts an interactive login flow for the named CLI plugin and returns immediately.
// The flow runs each entry in manifest.ResolveLoginSteps() sequentially against the same PTY
// dialog — output streams to the "cli.login.output" event throughout; "cli.login.done" is
// emitted once when the whole flow ends (last step success, first step failure, or cancellation).
func (p *Plugin) LoginCli(ctx context.Context, input plugin_dto.LoginCliInput) (*plugin_dto.LoginCliOutput, error) {
	config, err := p.loadConfig()
	if err != nil {
		return nil, ierror.Error(ierror.ErrSettingsLoadConfig, err)
	}
	item, ok := findExtension(config.Extensions, input.ID)
	if !ok || item.Kind != "cli" || !item.Enabled {
		return nil, ierror.Error(ierror.ErrCliLoginNotFound, os.ErrNotExist)
	}
	if item.RootDir == "" {
		return nil, ierror.Error(ierror.ErrCliLoginNotFound, os.ErrNotExist)
	}
	if _, statErr := os.Stat(item.RootDir); statErr != nil {
		return nil, ierror.Error(ierror.ErrCliLoginNotFound, statErr)
	}

	manifest, err := pkgcli.LoadManifest(item.ConfigFilePath)
	if err != nil {
		return nil, ierror.Error(ierror.ErrCliLoginStartFailed, err)
	}
	steps := manifest.ResolveLoginSteps()
	if len(steps) == 0 {
		return nil, ierror.Error(ierror.ErrCliLoginNoCommand, fmt.Errorf("manifest has no login_command/login_steps"))
	}

	name, err := cliNameFromID(item.ID)
	if err != nil {
		return nil, err
	}

	id := item.ID
	p.loginMu.Lock()
	if p.loginSessions[id] != nil {
		p.loginMu.Unlock()
		return nil, ierror.Error(ierror.ErrCliLoginSessionConflict, fmt.Errorf("session already active for %s", id))
	}

	// Start the first step synchronously so we can fail-fast on common errors
	// (executable missing, sentinel "no login_command" from the manager, etc.)
	// without making the user wait for an event roundtrip just to see the dialog flash.
	firstSession, startErr := p.startCliLoginCommand(ctx, name, steps[0])
	if startErr != nil {
		p.loginMu.Unlock()
		if strings.Contains(startErr.Error(), "no login_command") {
			return nil, ierror.Error(ierror.ErrCliLoginNoCommand, startErr)
		}
		return nil, ierror.Error(ierror.ErrCliLoginStartFailed, startErr)
	}
	p.loginSessions[id] = firstSession
	p.loginMu.Unlock()

	go p.runLoginFlow(ctx, id, name, steps, firstSession)
	return &plugin_dto.LoginCliOutput{}, nil
}

// runLoginFlow drives the multi-step login sequence in its own goroutine.
// The first session is supplied by LoginCli (already stored in the map); subsequent
// steps are started here.  After the loop, it emits cli.login.done exactly once and
// drops the map entry so a new login may begin.
func (p *Plugin) runLoginFlow(ctx context.Context, id, name string, steps [][]string, first loginSessionIface) {
	defer func() {
		p.loginMu.Lock()
		delete(p.loginSessions, id)
		delete(p.loginCancelled, id)
		p.loginMu.Unlock()
	}()

	var lastExit int
	var lastErr error
	current := first
	for i, argv := range steps {
		if i > 0 {
			p.loginMu.Lock()
			_, cancelled := p.loginCancelled[id]
			p.loginMu.Unlock()
			if cancelled {
				lastExit = 130
				lastErr = fmt.Errorf("cli login cancelled by user")
				break
			}

			p.emitLoginBanner(id, i, len(steps), argv)

			sess, err := p.startCliLoginCommand(ctx, name, argv)
			if err != nil {
				lastExit = -1
				lastErr = fmt.Errorf("start step %d (%s): %w", i+1, strings.Join(argv, " "), err)
				break
			}
			p.loginMu.Lock()
			p.loginSessions[id] = sess
			p.loginMu.Unlock()
			current = sess
		}

		for chunk := range current.Output() {
			p.emitLoginOutput(id, chunk)
		}
		code, werr := current.Wait()
		lastExit = code
		lastErr = werr

		if code != 0 || werr != nil {
			break
		}
	}

	errStr := ""
	if lastErr != nil {
		errStr = lastErr.Error()
	}
	if p.wailsApp != nil {
		p.wailsApp.Event.Emit("cli.login.done", loginDonePayload{
			ID:       id,
			ExitCode: lastExit,
			Error:    errStr,
		})
	}
}

// emitLoginBanner writes a synthetic step-transition banner into the same output channel
// the PTY uses, so xterm displays it inline with command output.
func (p *Plugin) emitLoginBanner(id string, stepNum, total int, argv []string) {
	if p.wailsApp == nil {
		return
	}
	banner := fmt.Sprintf("\r\n\x1b[36m=== step %d/%d: %s ===\x1b[0m\r\n", stepNum+1, total, strings.Join(argv, " "))
	p.wailsApp.Event.Emit("cli.login.output", loginOutputPayload{ID: id, Data: banner})
}

// emitLoginOutput forwards one PTY chunk to the frontend via cli.login.output.
func (p *Plugin) emitLoginOutput(id string, chunk []byte) {
	if p.wailsApp == nil {
		return
	}
	p.wailsApp.Event.Emit("cli.login.output", loginOutputPayload{ID: id, Data: string(chunk)})
}

// SendLoginStdin forwards raw bytes to the running login session's PTY stdin.
func (p *Plugin) SendLoginStdin(ctx context.Context, input plugin_dto.SendLoginStdinInput) error {
	_ = ctx
	p.loginMu.Lock()
	session := p.loginSessions[input.ID]
	p.loginMu.Unlock()
	if session == nil {
		return ierror.Error(ierror.ErrCliLoginNotFound, os.ErrNotExist)
	}
	if err := session.Write([]byte(input.Data)); err != nil {
		return fmt.Errorf("login stdin write: %w", err)
	}
	return nil
}

// ResizeLoginCli resizes the PTY window of a running login session.
func (p *Plugin) ResizeLoginCli(ctx context.Context, input plugin_dto.ResizeLoginCliInput) error {
	_ = ctx
	p.loginMu.Lock()
	session := p.loginSessions[input.ID]
	p.loginMu.Unlock()
	if session == nil {
		// The frontend may emit a trailing resize while the session is naturally
		// exiting; treat that race as a harmless no-op instead of surfacing an error.
		return nil
	}
	if err := session.Resize(input.Rows, input.Cols); err != nil {
		if strings.Contains(err.Error(), "already closed") {
			return nil
		}
		return fmt.Errorf("login resize: %w", err)
	}
	return nil
}

// CancelLoginCli terminates a running login flow. Sets the orchestrator's abort
// flag so the next step (if any) is skipped, and signals the current step's session
// so its PTY dies and the orchestrator advances out of its output drain.
func (p *Plugin) CancelLoginCli(ctx context.Context, input plugin_dto.CancelLoginCliInput) error {
	_ = ctx
	p.loginMu.Lock()
	session := p.loginSessions[input.ID]
	if _, running := p.loginSessions[input.ID]; running || session != nil {
		p.loginCancelled[input.ID] = struct{}{}
	}
	p.loginMu.Unlock()
	if session == nil {
		// Dialog teardown can race with the session finishing on its own.
		return nil
	}
	if err := session.Cancel(); err != nil {
		if strings.Contains(err.Error(), "already closed") {
			return nil
		}
		return err
	}
	return nil
}
