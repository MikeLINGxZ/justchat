package window

import (
	"context"

	"github.com/wailsapp/wails/v3/pkg/application"
)

func (p *Window) ServiceStartup(ctx context.Context, options application.ServiceOptions) error {
	p.wailsApp = application.Get()
	return nil
}
