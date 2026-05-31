package memory

import (
	"context"

	"github.com/wailsapp/wails/v3/pkg/application"
)

// ServiceStartup is called by Wails when the memory service starts.
func (m *Memory) ServiceStartup(ctx context.Context, options application.ServiceOptions) error {
	return nil
}
