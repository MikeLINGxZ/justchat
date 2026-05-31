package notification

import (
	"context"

	"github.com/wailsapp/wails/v3/pkg/application"
)

// ServiceStartup captures the Wails app reference for global event emission.
func (n *Notification) ServiceStartup(ctx context.Context, options application.ServiceOptions) error {
	n.wailsApp = application.Get()
	return nil
}
