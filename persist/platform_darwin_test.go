//go:build darwin && !ios

package persist

import "testing"

func TestIsOnScreen(t *testing.T) {
	// Primary display always has origin (0, 0) on macOS.
	if !IsOnScreen(100, 100) {
		t.Error("IsOnScreen(100, 100) = false, want true")
	}
	if IsOnScreen(-99999, -99999) {
		t.Error("IsOnScreen(-99999, -99999) = true, want false")
	}
}

func TestGetWindowFrameZeroHandle(t *testing.T) {
	x, y, w, h := GetWindowFrame(0)
	if x != 0 || y != 0 || w != 0 || h != 0 {
		t.Errorf("GetWindowFrame(0) = (%g, %g, %g, %g), want zeros", x, y, w, h)
	}
}
