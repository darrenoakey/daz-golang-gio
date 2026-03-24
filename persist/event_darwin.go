//go:build darwin && !ios

package persist

import "gioui.org/app"

// handlePlatformEvent intercepts macOS AppKitViewEvent to capture the native
// view handle for window position tracking.
func (w *Window) handlePlatformEvent(e any) {
	if ve, ok := e.(app.AppKitViewEvent); ok && ve.Valid() {
		w.setView(ve.View)
	}
}
