package window

import (
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/id/window_id"

	"github.com/wailsapp/wails/v3/pkg/application"
)

func (p *Window) childWindowOptions(options application.WebviewWindowOptions) application.WebviewWindowOptions {
	options.InitialPosition = application.WindowCentered
	options.Screen = p.homeWindowScreen()
	if options.Screen == nil && p.wailsApp != nil {
		options.Screen = p.wailsApp.Screen.GetPrimary()
	}
	return options
}

func (p *Window) showCenteredOnHomeScreen(window application.Window) {
	if window == nil {
		return
	}
	p.centerWindowOnHomeScreen(window)
	window.Show()
	p.centerWindowOnHomeScreen(window)
	window.Focus()
}

// centerWindowOnHomeScreen create a window at the main window screen
func (p *Window) centerWindowOnHomeScreen(window application.Window) {
	if window == nil {
		return
	}

	screen := p.homeWindowScreen()
	if screen == nil {
		if currentScreen, err := window.GetScreen(); err == nil {
			screen = currentScreen
		}
	}
	if screen == nil && p.wailsApp != nil {
		screen = p.wailsApp.Screen.GetPrimary()
	}
	if screen == nil {
		window.Center()
		return
	}

	window.SetScreen(screen)
}

// homeWindowScreen get main window screen
func (p *Window) homeWindowScreen() *application.Screen {
	if p.wailsApp == nil {
		return nil
	}
	homeWindow, ok := p.wailsApp.Window.GetByName(window_id.Home)
	if !ok {
		return nil
	}
	screen, err := homeWindow.GetScreen()
	if err != nil {
		return nil
	}
	return screen
}
