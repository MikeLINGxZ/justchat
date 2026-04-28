package service

import (
	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"
)

const EventFilesDropped = "files-dropped"

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
		EnableFileDrop:   true,
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

func (s *Service) RegisterFileDropHandler(window *application.WebviewWindow) {
	window.OnWindowEvent(events.Common.WindowFilesDropped, func(event *application.WindowEvent) {
		paths := event.Context().DroppedFiles()
		if len(paths) == 0 {
			return
		}
		fileInfos, err := s.fileInfo(paths)
		if err != nil {
			return
		}
		s.app.Event.Emit(EventFilesDropped, fileInfos)
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
	window := NewHomeWindow(s.app)
	s.RegisterFileDropHandler(window)
	return window
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

func (s *Service) centerWindowOnHomeScreen(window application.Window) {
	if window == nil {
		return
	}

	screen := s.homeWindowScreen()
	if screen == nil {
		if currentScreen, err := window.GetScreen(); err == nil {
			screen = currentScreen
		}
	}
	if screen == nil && s.app != nil {
		screen = s.app.Screen.GetPrimary()
	}
	if screen == nil {
		window.Center()
		return
	}

	if !centerWindowOnScreen(window, screen) {
		window.Center()
	}
}

func (s *Service) homeWindowScreen() *application.Screen {
	if s.app == nil {
		return nil
	}
	homeWindow, ok := s.app.Window.GetByName(WindowNameHome)
	if !ok {
		return nil
	}
	screen, err := homeWindow.GetScreen()
	if err != nil {
		return nil
	}
	return screen
}
