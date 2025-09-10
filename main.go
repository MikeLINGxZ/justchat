package main

import (
	"context"
	"embed"
	"fmt"
	"os"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/service"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {

	app := service.NewService()

	err := wails.Run(&options.App{
		Title:     "lemon_tea_desktop",
		Width:     1024,
		Height:    768,
		MinWidth:  350,
		MinHeight: 550,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup: func(ctx context.Context) {
			err := app.Startup(ctx)
			if err != nil {
				_, _ = runtime.MessageDialog(ctx, runtime.MessageDialogOptions{
					Type:    runtime.ErrorDialog,
					Title:   "程序错误",
					Message: fmt.Sprintf("%v", err.Error()),
				})
				os.Exit(-1)
			}
		},
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
