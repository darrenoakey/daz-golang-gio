//go:build darwin

package persist

import "gioui.org/app"

// handlePlatformEvent intercepts macOS AppKitViewEvent to capture the native
// view handle for window position tracking.
func (w *Window) handlePlatformEvent(e any) {
	if ve, ok := e.(app.AppKitViewEvent); ok {
		w.SetView(ve.View)
		w.RestorePosition(ve.View)
	}
}
