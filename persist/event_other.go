//go:build !darwin

package persist

// handlePlatformEvent is a no-op on non-macOS platforms.
// Window position tracking is not available without platform-specific APIs.
func (w *Window) handlePlatformEvent(e any) {}
