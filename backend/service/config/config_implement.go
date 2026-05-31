package config

import (
	"context"

	"github.com/wailsapp/wails/v3/pkg/application"
)

func (c *Config) ServiceStartup(ctx context.Context, options application.ServiceOptions) error {
	c.wailsApp = application.Get()
	return nil
}
