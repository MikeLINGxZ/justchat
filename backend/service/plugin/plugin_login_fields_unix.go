//go:build !windows

package plugin

import (
	"context"
)

// loginFields holds the platform-specific fields for CLI login session management.
// It is embedded in Plugin and is only active on non-Windows platforms.
type loginFields struct {
	// loginSessions tracks the current step's session per extension ID. During a
	// multi-step flow the orchestrator updates this entry as it advances.
	loginSessions map[string]loginSessionIface
	// loginCancelled records that the user has asked to abort the orchestrator's
	// remaining steps. The current session is also signaled via its own Cancel.
	loginCancelled map[string]struct{}
	// startCliLoginCommand starts one login step's argv vector. In production it
	// delegates to cliManager.StartLoginCommand; tests replace it with a fake.
	startCliLoginCommand func(ctx context.Context, name string, argv []string) (loginSessionIface, error)
}

// initLoginFields initialises the login-specific fields of a newly constructed Plugin.
func initLoginFields(p *Plugin) {
	p.loginSessions = map[string]loginSessionIface{}
	p.loginCancelled = map[string]struct{}{}
	p.startCliLoginCommand = func(ctx context.Context, name string, argv []string) (loginSessionIface, error) {
		mgr, err := p.resolveCliManager()
		if err != nil {
			return nil, err
		}
		return mgr.StartLoginCommand(ctx, name, argv)
	}
}
