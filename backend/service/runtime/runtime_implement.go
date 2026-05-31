package runtime

import (
	"context"

	"github.com/wailsapp/wails/v3/pkg/application"
)

// ServiceStartup captures the Wails application handle once the service is registered.
func (s *Runtime) ServiceStartup(ctx context.Context, options application.ServiceOptions) error {
	s.wailsApp = application.Get()
	return nil
}
