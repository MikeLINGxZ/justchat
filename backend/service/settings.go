package service

import "github.com/wailsapp/wails/v3/pkg/application"

func (s *Service) OpenSettingsWindow() {
	settingsWindow, ok := s.app.Window.GetByName(WindowNameSettings)
	if ok {
		settingsWindow.Focus()
		settingsWindow.Center()
		settingsWindow.Show()
		return
	}
	settingsWindow = s.app.Window.NewWithOptions(application.WebviewWindowOptions{
		Name:  WindowNameSettings,
		Title: "Settings",
		Mac: application.MacWindow{
			InvisibleTitleBarHeight: 50,
			Backdrop:                application.MacBackdropTranslucent,
			TitleBar:                application.MacTitleBarDefault,
		},
		BackgroundColour: application.NewRGB(27, 38, 54),
		URL:              "/settings",
		Width:            1000,
		Height:           700,
		MinWidth:         350,
		MinHeight:        550,
	})
	settingsWindow.Focus()
	settingsWindow.Center()
	settingsWindow.Show()
}
