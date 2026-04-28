//go:build !darwin

package service

import "github.com/wailsapp/wails/v3/pkg/application"

func centerWindowOnScreen(window application.Window, screen *application.Screen) bool {
	if window == nil || screen == nil {
		return false
	}

	width, height := window.Size()
	if width <= 0 || height <= 0 {
		bounds := window.Bounds()
		width = bounds.Width
		height = bounds.Height
	}
	if width <= 0 || height <= 0 {
		return false
	}

	workArea := screen.WorkArea
	if workArea.Width <= 0 || workArea.Height <= 0 {
		workArea = screen.Bounds
	}

	x := workArea.X + (workArea.Width-width)/2
	y := workArea.Y + (workArea.Height-height)/2
	if x < workArea.X {
		x = workArea.X
	}
	if y < workArea.Y {
		y = workArea.Y
	}
	window.SetPosition(x, y)
	return true
}
