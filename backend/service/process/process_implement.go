package process

import (
	"context"

	"github.com/wailsapp/wails/v3/pkg/application"
)

func (p *Process) ServiceStartup(ctx context.Context, options application.ServiceOptions) error {
	p.wailsApp = application.Get()
	return nil
}
