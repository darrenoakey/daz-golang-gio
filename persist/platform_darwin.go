//go:build darwin

package persist

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Cocoa

#import <Cocoa/Cocoa.h>

struct WindowFrame {
	double x;
	double y;
	double width;
	double height;
};

struct WindowFrame getWindowFrame(uintptr_t viewPtr) {
	NSView *view = (__bridge NSView *)(void *)viewPtr;
	NSWindow *window = [view window];
	struct WindowFrame wf = {0, 0, 0, 0};
	if (window == nil) return wf;
	NSRect frame = [window frame];
	wf.x = frame.origin.x;
	wf.y = frame.origin.y;
	wf.width = frame.size.width;
	wf.height = frame.size.height;
	return wf;
}

void setWindowFrame(uintptr_t viewPtr, double x, double y, double width, double height) {
	NSView *view = (__bridge NSView *)(void *)viewPtr;
	NSWindow *window = [view window];
	if (window == nil) return;
	NSRect frame = NSMakeRect(x, y, width, height);
	[window setFrame:frame display:YES animate:NO];
}
*/
import "C"

// GetWindowPosition reads the current window frame from the native macOS window.
func GetWindowPosition(view uintptr) (x, y, width, height int) {
	if view == 0 {
		return 0, 0, 0, 0
	}
	frame := C.getWindowFrame(C.uintptr_t(view))
	return int(frame.x), int(frame.y), int(frame.width), int(frame.height)
}

// SetWindowPosition moves and resizes the native macOS window.
func SetWindowPosition(view uintptr, x, y, width, height int) {
	if view == 0 {
		return
	}
	C.setWindowFrame(C.uintptr_t(view), C.double(x), C.double(y), C.double(width), C.double(height))
}

// PositionSupported returns true on macOS where native position access is available.
func PositionSupported() bool {
	return true
}
