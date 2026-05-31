//go:build windows

package plugin

import (
	"context"
	"errors"

	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/ierror"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/service/plugin/plugin_dto"
)

var errCliLoginUnsupportedOnWindows = errors.New("cli login is not supported on Windows")

// LoginCli starts an interactive login flow for the named CLI plugin and returns immediately.
// The flow runs each entry in manifest.ResolveLoginSteps() sequentially against the same PTY
// dialog — output streams to the "cli.login.output" event throughout; "cli.login.done" is
// emitted once when the whole flow ends (last step success, first step failure, or cancellation).
func (p *Plugin) LoginCli(ctx context.Context, input plugin_dto.LoginCliInput) (*plugin_dto.LoginCliOutput, error) {
	_ = ctx
	_ = input
	return nil, ierror.Error(ierror.ErrRuntimeUnsupportedOS, errCliLoginUnsupportedOnWindows)
}

// SendLoginStdin forwards raw bytes to the running login session's PTY stdin.
func (p *Plugin) SendLoginStdin(ctx context.Context, input plugin_dto.SendLoginStdinInput) error {
	_ = ctx
	_ = input
	return ierror.Error(ierror.ErrRuntimeUnsupportedOS, errCliLoginUnsupportedOnWindows)
}

// ResizeLoginCli resizes the PTY window of a running login session.
func (p *Plugin) ResizeLoginCli(ctx context.Context, input plugin_dto.ResizeLoginCliInput) error {
	_ = ctx
	_ = input
	return ierror.Error(ierror.ErrRuntimeUnsupportedOS, errCliLoginUnsupportedOnWindows)
}

// CancelLoginCli terminates a running login flow. Sets the orchestrator's abort
// flag so the next step (if any) is skipped, and signals the current step's session
// so its PTY dies and the orchestrator advances out of its output drain.
func (p *Plugin) CancelLoginCli(ctx context.Context, input plugin_dto.CancelLoginCliInput) error {
	_ = ctx
	_ = input
	return nil
}
