//go:build darwin && !ios

// Package macos provides macOS-specific utilities for Gio apps.
package macos

/*
#cgo CFLAGS: -x objective-c -fobjc-arc
#cgo LDFLAGS: -framework AppKit

#import <AppKit/AppKit.h>
#include <stdint.h>
#include <dispatch/dispatch.h>

static void setDockIconFromData(const void *data, int len) {
	void (^work)(void) = ^{
		@autoreleasepool {
			NSData *imgData = [NSData dataWithBytes:data length:len];
			NSImage *img = [[NSImage alloc] initWithData:imgData];
			if (img) {
				[[NSApplication sharedApplication] setApplicationIconImage:img];
			}
		}
	};
	if ([NSThread isMainThread]) {
		work();
	} else {
		dispatch_async(dispatch_get_main_queue(), work);
	}
}
*/
import "C"
import "unsafe"

// SetDockIcon sets the macOS dock icon from PNG/JPEG image data.
// Pass the contents of an embedded image file (e.g., //go:embed icon.png).
// Safe to call from any goroutine (dispatches to main thread).
// No-op if data is empty.
func SetDockIcon(data []byte) {
	if len(data) == 0 {
		return
	}
	C.setDockIconFromData(unsafe.Pointer(&data[0]), C.int(len(data)))
}
