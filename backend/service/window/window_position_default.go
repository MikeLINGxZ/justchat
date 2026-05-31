//go:build !darwin

package window

import "github.com/wailsapp/wails/v3/pkg/application"

func centerWindowOnScreen(window application.Window, screen *application.Screen) bool {
	_ = window
	_ = screen
	return false
}
