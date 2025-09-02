package main

import (
	"context"
	"embed"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/core/cloud"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Create an instance of the app structure

	auth := cloud.NewAuth()

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "lemon_tea_desktop_temp",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup: func(ctx context.Context) {
			auth.Startup(ctx)
		},
		Bind: []interface{}{
			auth,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
