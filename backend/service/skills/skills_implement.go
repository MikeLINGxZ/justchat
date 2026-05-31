package skills

import (
	"context"

	"github.com/wailsapp/wails/v3/pkg/application"
)

// ServiceStartup captures the Wails application handle once the service is registered.
func (s *Skills) ServiceStartup(ctx context.Context, options application.ServiceOptions) error {
	_ = ctx
	_ = options
	s.wailsApp = application.Get()
	return nil
}
