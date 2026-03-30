package service

import (
	"github.com/wailsapp/wails/v3/pkg/application"
)

func NewHomeWindow(app *application.App) *application.WebviewWindow {
	return app.Window.NewWithOptions(application.WebviewWindowOptions{
		Name:  WindowNameHome,
		Title: "Home",
		Mac: application.MacWindow{
			InvisibleTitleBarHeight: 50,
			Backdrop:                application.MacBackdropTranslucent,
			TitleBar:                application.MacTitleBarDefault,
		},
		BackgroundColour: application.NewRGB(27, 38, 54),
		URL:              "/",
		Width:            1300,
		Height:           860,
		MinWidth:         350,
		MinHeight:        550,
	})
}

func NewOnboardingWindow(app *application.App) *application.WebviewWindow {
	return app.Window.NewWithOptions(application.WebviewWindowOptions{
		Name:             WindowNameOnboarding,
		Title:            "Welcome",
		URL:              "/?entry=onboarding",
		Width:            980,
		Height:           780,
		MinWidth:         900,
		MinHeight:        690,
		DisableResize:    true,
		BackgroundColour: application.NewRGB(244, 239, 228),
		Mac: application.MacWindow{
			InvisibleTitleBarHeight: 50,
			Backdrop:                application.MacBackdropTranslucent,
			TitleBar:                application.MacTitleBarDefault,
		},
	})
}

func (s *Service) ensureHomeWindow() *application.WebviewWindow {
	if s.app == nil {
		return nil
	}
	if homeWindow, ok := s.app.Window.GetByName(WindowNameHome); ok {
		if webviewWindow, ok := homeWindow.(*application.WebviewWindow); ok {
			return webviewWindow
		}
		return nil
	}
	return NewHomeWindow(s.app)
}

func (s *Service) showHomeWindow() {
	homeWindow := s.ensureHomeWindow()
	if homeWindow == nil {
		return
	}
	homeWindow.Center()
	homeWindow.Show()
	homeWindow.Focus()
}
