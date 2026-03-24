//go:build !darwin || ios

package persist

// GetWindowFrame is a no-op on non-macOS platforms.
func GetWindowFrame(view uintptr) (x, y, width, height float64) {
	return 0, 0, 0, 0
}

// SetWindowFrame is a no-op on non-macOS platforms.
func SetWindowFrame(view uintptr, x, y, width, height float64) {}

// IsOnScreen always returns false on non-macOS platforms.
func IsOnScreen(x, y float64) bool {
	return false
}

// PositionSupported returns false on platforms without native position access.
func PositionSupported() bool {
	return false
}
