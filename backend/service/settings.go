package service

import "github.com/wailsapp/wails/v3/pkg/application"

func (s *Service) OpenSettingsWindow() {
	webviewWindow := s.app.Window.NewWithOptions(application.WebviewWindowOptions{
		Name:  "Settings",
		Title: "Settings",
		Mac: application.MacWindow{
			InvisibleTitleBarHeight: 50,
			Backdrop:                application.MacBackdropTranslucent,
			TitleBar:                application.MacTitleBarDefault,
		},
		BackgroundColour: application.NewRGB(27, 38, 54),
		URL:              "/settings",
		Width:            800,
		Height:           700,
		MinWidth:         350,
		MinHeight:        550,
		AlwaysOnTop:      true,
	})
	webviewWindow.Show()
}
