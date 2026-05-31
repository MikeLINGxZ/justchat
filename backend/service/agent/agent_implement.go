package agent

import (
	"context"

	"github.com/wailsapp/wails/v3/pkg/application"
)

// ServiceStartup is called by Wails when the application starts.
func (a *Agent) ServiceStartup(ctx context.Context, options application.ServiceOptions) error {
	a.manager.SetApp(application.Get())
	return nil
}
