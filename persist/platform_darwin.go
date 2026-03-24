//go:build darwin && !ios

package persist

/*
#cgo CFLAGS: -x objective-c -fobjc-arc
#cgo LDFLAGS: -framework AppKit

#import <AppKit/AppKit.h>
#include <stdint.h>
#include <dispatch/dispatch.h>

// getWindowFrame reads the window frame on the main thread via dispatch_sync.
static void getWindowFrame(uintptr_t viewRef, CGFloat *ox, CGFloat *oy, CGFloat *ow, CGFloat *oh) {
	__block CGFloat bx = 0, by = 0, bw = 0, bh = 0;
	void (^work)(void) = ^{
		@autoreleasepool {
			NSView *view = (__bridge NSView *)(void *)viewRef;
			NSWindow *window = view.window;
			if (!window) return;
			NSRect frame = [window frame];
			bx = frame.origin.x;
			by = frame.origin.y;
			bw = frame.size.width;
			bh = frame.size.height;
		}
	};
	if ([NSThread isMainThread]) {
		work();
	} else {
		dispatch_sync(dispatch_get_main_queue(), work);
	}
	*ox = bx; *oy = by; *ow = bw; *oh = bh;
}

// setWindowFrame sets the window frame on the main thread via dispatch_async.
static void setWindowFrame(uintptr_t viewRef, CGFloat x, CGFloat y, CGFloat w, CGFloat h) {
	void (^work)(void) = ^{
		@autoreleasepool {
			NSView *view = (__bridge NSView *)(void *)viewRef;
			NSWindow *window = view.window;
			if (!window) return;
			NSRect r = NSMakeRect(x, y, w, h);
			[window setFrame:r display:YES];
		}
	};
	if ([NSThread isMainThread]) {
		work();
	} else {
		dispatch_async(dispatch_get_main_queue(), work);
	}
}

// isPointOnScreen checks if (x, y) falls within any connected screen.
// NSScreen.screens is thread-safe for reading — no main thread dispatch needed.
static int isPointOnScreen(CGFloat x, CGFloat y) {
	@autoreleasepool {
		for (NSScreen *screen in [NSScreen screens]) {
			if (NSPointInRect(NSMakePoint(x, y), screen.frame)) {
				return 1;
			}
		}
		return 0;
	}
}
*/
import "C"

// GetWindowFrame reads the current window frame from the native macOS window.
// Safe to call from any goroutine (dispatches to main thread).
func GetWindowFrame(view uintptr) (x, y, width, height float64) {
	if view == 0 {
		return 0, 0, 0, 0
	}
	var cx, cy, cw, ch C.CGFloat
	C.getWindowFrame(C.uintptr_t(view), &cx, &cy, &cw, &ch)
	return float64(cx), float64(cy), float64(cw), float64(ch)
}

// SetWindowFrame moves and resizes the native macOS window.
// Safe to call from any goroutine (dispatches to main thread).
func SetWindowFrame(view uintptr, x, y, width, height float64) {
	if view == 0 {
		return
	}
	C.setWindowFrame(C.uintptr_t(view), C.CGFloat(x), C.CGFloat(y), C.CGFloat(width), C.CGFloat(height))
}

// IsOnScreen returns true if the given point is within any connected display.
// Safe to call from any goroutine.
func IsOnScreen(x, y float64) bool {
	return C.isPointOnScreen(C.CGFloat(x), C.CGFloat(y)) == 1
}

// PositionSupported returns true on macOS where native position access is available.
func PositionSupported() bool {
	return true
}
