package plugin

import (
	"context"

	"github.com/wailsapp/wails/v3/pkg/application"
)

// ServiceStartup captures the Wails application handle once the service is registered.
func (p *Plugin) ServiceStartup(ctx context.Context, options application.ServiceOptions) error {
	_ = ctx
	_ = options
	p.wailsApp = application.Get()
	return nil
}
