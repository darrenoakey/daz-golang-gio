//go:build !darwin

package persist

// GetWindowPosition is a no-op on non-macOS platforms.
// Gio does not expose window position portably.
func GetWindowPosition(view uintptr) (x, y, width, height int) {
	return 0, 0, 0, 0
}

// SetWindowPosition is a no-op on non-macOS platforms.
func SetWindowPosition(view uintptr, x, y, width, height int) {}

// PositionSupported returns false on platforms without native position access.
func PositionSupported() bool {
	return false
}
