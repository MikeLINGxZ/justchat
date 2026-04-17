package main

import (
	"embed"
	_ "embed"
	"log"
	"time"

	"github.com/wailsapp/wails/v3/pkg/application"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/logger"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/service"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	logger.NewStaticLogger("lemontea")
	serviceInstance := service.NewService()
	app := application.New(application.Options{
		Name:        "lemon_tea_desktop",
		Description: "A ai agent client",
		Services: []application.Service{
			application.NewService(serviceInstance),
		},
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
	})

	initialized, err := service.IsAppInitialized()
	if err != nil {
		log.Fatal(err)
	}

	if initialized {
		homeWindow := service.NewHomeWindow(app)
		serviceInstance.RegisterFileDropHandler(homeWindow)
	} else {
		service.NewOnboardingWindow(app)
	}

	go func() {
		for {
			now := time.Now().Format(time.RFC1123)
			app.Event.Emit("time", now)
			time.Sleep(time.Second)
		}
	}()

	err = app.Run()
	if err != nil {
		log.Fatal(err)
	}
}
