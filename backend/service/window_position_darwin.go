//go:build darwin

package service

/*
#cgo CFLAGS: -x objective-c -fblocks
#cgo LDFLAGS: -framework Cocoa
#import <Cocoa/Cocoa.h>
#import <dispatch/dispatch.h>

static void centerWindowOnScreenDarwin(void* nsWindow, double x, double y, double width, double height) {
	NSWindow* window = (NSWindow*)nsWindow;
	void (^block)(void) = ^{
		NSRect frame = [window frame];
		CGFloat windowX = x + (width - frame.size.width) / 2.0;
		CGFloat windowY = y + (height - frame.size.height) / 2.0;
		if (windowX < x) {
			windowX = x;
		}
		if (windowY < y) {
			windowY = y;
		}
		[window setFrame:NSMakeRect(windowX, windowY, frame.size.width, frame.size.height) display:YES animate:NO];
	};
	if ([NSThread isMainThread]) {
		block();
	} else {
		dispatch_sync(dispatch_get_main_queue(), block);
	}
}
*/
import "C"

import "github.com/wailsapp/wails/v3/pkg/application"

func centerWindowOnScreen(window application.Window, screen *application.Screen) bool {
	if window == nil || screen == nil {
		return false
	}

	workArea := screen.WorkArea
	if workArea.Width <= 0 || workArea.Height <= 0 {
		workArea = screen.Bounds
	}
	if workArea.Width <= 0 || workArea.Height <= 0 {
		return false
	}

	C.centerWindowOnScreenDarwin(
		window.NativeWindow(),
		C.double(workArea.X),
		C.double(workArea.Y),
		C.double(workArea.Width),
		C.double(workArea.Height),
	)
	return true
}
