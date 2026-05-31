package window

import (
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/id/window_id"

	"github.com/wailsapp/wails/v3/pkg/application"
)

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

	if !centerWindowOnScreen(window, screen) {
		window.Center()
	}
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
