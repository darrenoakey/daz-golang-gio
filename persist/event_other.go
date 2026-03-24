//go:build !darwin || ios

package persist

// handlePlatformEvent is a no-op on non-macOS platforms.
func (w *Window) handlePlatformEvent(e any) {}
