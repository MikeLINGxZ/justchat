package service

import (
	"github.com/wailsapp/wails/v3/pkg/application"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/i18n"
)

func (s *Service) OpenSettingsAboutWindow() {
	settingsWindow, ok := s.app.Window.GetByName(WindowNameSettings)
	if ok {
		settingsWindow.SetURL("/?entry=settings&tab=about")
		settingsWindow.Focus()
		s.centerWindowOnHomeScreen(settingsWindow)
		settingsWindow.Show()
		return
	}
	settingsWindow = s.app.Window.NewWithOptions(application.WebviewWindowOptions{
		Name:  WindowNameSettings,
		Title: i18n.TCurrent("app.window.settings_title", nil),
		Mac: application.MacWindow{
			InvisibleTitleBarHeight: 50,
			Backdrop:                application.MacBackdropTranslucent,
			TitleBar:                application.MacTitleBarDefault,
		},
		BackgroundColour: application.NewRGB(27, 38, 54),
		URL:              "/?entry=settings&tab=about",
		Width:            1200,
		Height:           800,
		MinWidth:         350,
		MinHeight:        550,
	})
	settingsWindow.Focus()
	s.centerWindowOnHomeScreen(settingsWindow)
	settingsWindow.Show()
}

func (s *Service) OpenSettingsWindow() {
	settingsWindow, ok := s.app.Window.GetByName(WindowNameSettings)
	if ok {
		settingsWindow.Focus()
		s.centerWindowOnHomeScreen(settingsWindow)
		settingsWindow.Show()
		return
	}
	settingsWindow = s.app.Window.NewWithOptions(application.WebviewWindowOptions{
		Name:  WindowNameSettings,
		Title: i18n.TCurrent("app.window.settings_title", nil),
		Mac: application.MacWindow{
			InvisibleTitleBarHeight: 50,
			Backdrop:                application.MacBackdropTranslucent,
			TitleBar:                application.MacTitleBarDefault,
		},
		BackgroundColour: application.NewRGB(27, 38, 54),
		URL:              "/?entry=settings",
		Width:            1200,
		Height:           800,
		MinWidth:         350,
		MinHeight:        550,
	})
	settingsWindow.Focus()
	s.centerWindowOnHomeScreen(settingsWindow)
	settingsWindow.Show()
}
