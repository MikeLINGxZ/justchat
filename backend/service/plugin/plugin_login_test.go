//go:build !windows

package plugin

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	pkgcli "gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/cli"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/ierror"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/service/plugin/plugin_dto"
)

// setManifestLoginSteps overwrites the installed manifest's LoginSteps for tests
// that need the service-level "no login_command" guard to pass.
func setManifestLoginSteps(t *testing.T, configPath string, steps [][]string) {
	t.Helper()
	manifest, err := pkgcli.LoadManifest(configPath)
	if err != nil {
		t.Fatalf("load manifest: %v", err)
	}
	manifest.LoginSteps = steps
	if err := pkgcli.SaveManifest(configPath, manifest); err != nil {
		t.Fatalf("save manifest: %v", err)
	}
}

// fakeLoginSession is a controllable stand-in for loginSessionIface used in tests.
type fakeLoginSession struct {
	output   chan []byte
	done     chan struct{}
	written  []string
	mu       sync.Mutex
	closed   bool
	exitCode int
	waitErr  error
}

func newFakeLoginSession() *fakeLoginSession {
	return &fakeLoginSession{
		output: make(chan []byte, 16),
		done:   make(chan struct{}),
	}
}

// finish signals session completion; unblocks Wait and closes Output.
func (f *fakeLoginSession) finish(code int, err error) {
	f.mu.Lock()
	if f.closed {
		f.mu.Unlock()
		return
	}
	f.exitCode = code
	f.waitErr = err
	f.closed = true
	f.mu.Unlock()
	close(f.output)
	close(f.done)
}

func (f *fakeLoginSession) Output() <-chan []byte { return f.output }

func (f *fakeLoginSession) Write(data []byte) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.closed {
		return errors.New("cli: login session already closed")
	}
	f.written = append(f.written, string(data))
	return nil
}

func (f *fakeLoginSession) Resize(rows, cols uint16) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.closed {
		return errors.New("cli: login session already closed")
	}
	return nil
}

func (f *fakeLoginSession) Cancel() error {
	f.mu.Lock()
	closed := f.closed
	f.mu.Unlock()
	if closed {
		return errors.New("cli: login session already closed")
	}
	f.finish(0, nil)
	return nil
}

func (f *fakeLoginSession) Wait() (int, error) {
	<-f.done
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.exitCode, f.waitErr
}

// newTestPluginWithFakeSession returns a Plugin with a fake CLI extension installed (manifest
// has a one-step login_command so the service's "no login_command" check passes) and
// startCliLoginCommand replaced by one that returns the provided fake (or startErr if non-nil).
func newTestPluginWithFakeSession(t *testing.T, fake loginSessionIface, startErr error) (*Plugin, string) {
	t.Helper()
	p := NewPlugin()
	_ = withFakeCliManager(t, p)

	ctx := context.Background()
	out, err := p.InstallCliFromNpm(ctx, plugin_dto.InstallCliFromNpmInput{NpmPackage: "lark-cli", Name: "lark-cli"})
	if err != nil {
		t.Fatalf("install lark-cli: %v", err)
	}
	setManifestLoginSteps(t, out.Extension.ConfigFilePath, [][]string{{"auth", "login"}})

	p.startCliLoginCommand = func(_ context.Context, _ string, _ []string) (loginSessionIface, error) {
		if startErr != nil {
			return nil, startErr
		}
		return fake, nil
	}
	return p, out.Extension.ID
}

// TestLoginCliConflictsOnSecondCall verifies that a second LoginCli call while the first
// session is still active returns ErrCliLoginSessionConflict, and that once the first
// session finishes the map slot is freed and a new call can proceed.
func TestLoginCliConflictsOnSecondCall(t *testing.T) {
	fake := newFakeLoginSession()
	p, id := newTestPluginWithFakeSession(t, fake, nil)
	ctx := context.Background()

	// First call must succeed.
	if _, err := p.LoginCli(ctx, plugin_dto.LoginCliInput{ID: id}); err != nil {
		t.Fatalf("first LoginCli: %v", err)
	}

	// Second call must return conflict.
	_, err := p.LoginCli(ctx, plugin_dto.LoginCliInput{ID: id})
	if !errors.Is(err, ierror.Error(ierror.ErrCliLoginSessionConflict, errors.New("x"))) {
		t.Fatalf("expected ErrCliLoginSessionConflict, got: %v", err)
	}

	// Finish the first session so the goroutine drains and cleans up the map.
	fake.finish(0, nil)

	// Poll until the cleanup goroutine has removed the entry.
	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		fake2 := newFakeLoginSession()
		p.startCliLoginCommand = func(_ context.Context, _ string, _ []string) (loginSessionIface, error) {
			return fake2, nil
		}
		_, thirdErr := p.LoginCli(ctx, plugin_dto.LoginCliInput{ID: id})
		if thirdErr == nil {
			fake2.finish(0, nil)
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatal("third LoginCli did not succeed after first session ended")
}

// TestSendLoginStdinForwardsToSession verifies that SendLoginStdin writes the supplied
// data to the underlying session's PTY stdin.
func TestSendLoginStdinForwardsToSession(t *testing.T) {
	fake := newFakeLoginSession()
	p, id := newTestPluginWithFakeSession(t, fake, nil)
	ctx := context.Background()

	if _, err := p.LoginCli(ctx, plugin_dto.LoginCliInput{ID: id}); err != nil {
		t.Fatalf("LoginCli: %v", err)
	}

	if err := p.SendLoginStdin(ctx, plugin_dto.SendLoginStdinInput{ID: id, Data: "hi\n"}); err != nil {
		t.Fatalf("SendLoginStdin: %v", err)
	}

	fake.mu.Lock()
	written := append([]string(nil), fake.written...)
	fake.mu.Unlock()

	if len(written) != 1 || written[0] != "hi\n" {
		t.Fatalf("expected written [\"hi\\n\"], got %v", written)
	}

	fake.finish(0, nil)
}

// TestLoginCliFailsWhenManifestHasNoLoginCommand verifies that when the manifest has neither
// login_command nor login_steps, LoginCli short-circuits with ErrCliLoginNoCommand and never
// reaches the startCliLoginCommand injection point.
func TestLoginCliFailsWhenManifestHasNoLoginCommand(t *testing.T) {
	p := NewPlugin()
	_ = withFakeCliManager(t, p)

	ctx := context.Background()
	out, err := p.InstallCliFromNpm(ctx, plugin_dto.InstallCliFromNpmInput{NpmPackage: "lark-cli", Name: "lark-cli"})
	if err != nil {
		t.Fatalf("install lark-cli: %v", err)
	}

	called := false
	p.startCliLoginCommand = func(_ context.Context, _ string, _ []string) (loginSessionIface, error) {
		called = true
		return nil, nil
	}

	_, err = p.LoginCli(ctx, plugin_dto.LoginCliInput{ID: out.Extension.ID})
	if !errors.Is(err, ierror.Error(ierror.ErrCliLoginNoCommand, errors.New("x"))) {
		t.Fatalf("expected ErrCliLoginNoCommand, got: %v", err)
	}
	if called {
		t.Fatal("startCliLoginCommand should not have been called when manifest has no login_command")
	}
}

// TestResizeLoginCliIgnoresMissingSession verifies that resize calls become a no-op
// once the login session has already exited and been cleaned up.
func TestResizeLoginCliIgnoresMissingSession(t *testing.T) {
	fake := newFakeLoginSession()
	p, id := newTestPluginWithFakeSession(t, fake, nil)
	ctx := context.Background()

	if _, err := p.LoginCli(ctx, plugin_dto.LoginCliInput{ID: id}); err != nil {
		t.Fatalf("LoginCli: %v", err)
	}

	fake.finish(0, nil)

	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		p.loginMu.Lock()
		_, exists := p.loginSessions[id]
		p.loginMu.Unlock()
		if !exists {
			if err := p.ResizeLoginCli(ctx, plugin_dto.ResizeLoginCliInput{ID: id, Rows: 24, Cols: 80}); err != nil {
				t.Fatalf("ResizeLoginCli after cleanup: %v", err)
			}
			return
		}
		time.Sleep(10 * time.Millisecond)
	}

	t.Fatal("login session was not cleaned up in time")
}

// TestLoginCliRunsStepsSequentially verifies that login_steps with multiple entries are
// invoked in order against the same dialog: step 2 starts only after step 1 exits zero.
func TestLoginCliRunsStepsSequentially(t *testing.T) {
	p := NewPlugin()
	_ = withFakeCliManager(t, p)

	ctx := context.Background()
	out, err := p.InstallCliFromNpm(ctx, plugin_dto.InstallCliFromNpmInput{NpmPackage: "lark-cli", Name: "lark-cli"})
	if err != nil {
		t.Fatalf("install: %v", err)
	}
	setManifestLoginSteps(t, out.Extension.ConfigFilePath, [][]string{
		{"config", "init", "--new"},
		{"auth", "login"},
	})

	fake1 := newFakeLoginSession()
	fake2 := newFakeLoginSession()
	queue := []loginSessionIface{fake1, fake2}
	var observedArgv [][]string
	var mu sync.Mutex

	p.startCliLoginCommand = func(_ context.Context, _ string, argv []string) (loginSessionIface, error) {
		mu.Lock()
		defer mu.Unlock()
		observedArgv = append(observedArgv, append([]string(nil), argv...))
		if len(queue) == 0 {
			return nil, errors.New("test: no more fake sessions in queue")
		}
		next := queue[0]
		queue = queue[1:]
		return next, nil
	}

	if _, err := p.LoginCli(ctx, plugin_dto.LoginCliInput{ID: out.Extension.ID}); err != nil {
		t.Fatalf("LoginCli: %v", err)
	}

	// Step 1 finishes successfully; orchestrator should advance to step 2.
	fake1.finish(0, nil)

	// Wait until step 2's session is recorded in the map.
	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		p.loginMu.Lock()
		current := p.loginSessions[out.Extension.ID]
		p.loginMu.Unlock()
		if current == fake2 {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
	p.loginMu.Lock()
	if p.loginSessions[out.Extension.ID] != fake2 {
		p.loginMu.Unlock()
		t.Fatalf("expected fake2 to be active after step 1 finished")
	}
	p.loginMu.Unlock()

	fake2.finish(0, nil)

	// Wait for cleanup.
	deadline = time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		p.loginMu.Lock()
		_, still := p.loginSessions[out.Extension.ID]
		p.loginMu.Unlock()
		if !still {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}

	mu.Lock()
	defer mu.Unlock()
	if len(observedArgv) != 2 {
		t.Fatalf("expected 2 invocations, got %d: %v", len(observedArgv), observedArgv)
	}
	if observedArgv[0][0] != "config" || observedArgv[1][0] != "auth" {
		t.Fatalf("unexpected argv order: %v", observedArgv)
	}
}

// TestLoginCliStopsOnFirstFailingStep verifies that if step N exits non-zero,
// remaining steps are skipped and a single cli.login.done is implicit (we only
// assert by checking the second fake was never called).
func TestLoginCliStopsOnFirstFailingStep(t *testing.T) {
	p := NewPlugin()
	_ = withFakeCliManager(t, p)

	ctx := context.Background()
	out, err := p.InstallCliFromNpm(ctx, plugin_dto.InstallCliFromNpmInput{NpmPackage: "lark-cli", Name: "lark-cli"})
	if err != nil {
		t.Fatalf("install: %v", err)
	}
	setManifestLoginSteps(t, out.Extension.ConfigFilePath, [][]string{
		{"step1"},
		{"step2"},
	})

	fake1 := newFakeLoginSession()
	var calls int
	p.startCliLoginCommand = func(_ context.Context, _ string, _ []string) (loginSessionIface, error) {
		calls++
		if calls == 1 {
			return fake1, nil
		}
		return nil, errors.New("test: step 2 must not be called")
	}

	if _, err := p.LoginCli(ctx, plugin_dto.LoginCliInput{ID: out.Extension.ID}); err != nil {
		t.Fatalf("LoginCli: %v", err)
	}

	// Fail step 1.
	fake1.finish(2, errors.New("step1 failed"))

	// Wait for orchestrator to drain and clean up.
	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		p.loginMu.Lock()
		_, still := p.loginSessions[out.Extension.ID]
		p.loginMu.Unlock()
		if !still {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}

	if calls != 1 {
		t.Fatalf("expected only step 1 to be invoked, got %d invocations", calls)
	}
}

// TestCancelLoginCliIgnoresMissingSession verifies that cancel remains safe during
// dialog teardown even if the backend session has already finished naturally.
func TestCancelLoginCliIgnoresMissingSession(t *testing.T) {
	fake := newFakeLoginSession()
	p, id := newTestPluginWithFakeSession(t, fake, nil)
	ctx := context.Background()

	if _, err := p.LoginCli(ctx, plugin_dto.LoginCliInput{ID: id}); err != nil {
		t.Fatalf("LoginCli: %v", err)
	}

	fake.finish(0, nil)

	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		p.loginMu.Lock()
		_, exists := p.loginSessions[id]
		p.loginMu.Unlock()
		if !exists {
			if err := p.CancelLoginCli(ctx, plugin_dto.CancelLoginCliInput{ID: id}); err != nil {
				t.Fatalf("CancelLoginCli after cleanup: %v", err)
			}
			return
		}
		time.Sleep(10 * time.Millisecond)
	}

	t.Fatal("login session was not cleaned up in time")
}
