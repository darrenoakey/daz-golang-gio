package persist

import "testing"

// TestShouldRedrawIdleTickDoesNotRedraw pins the fix for the memory leak where
// tracker() invalidated the Gio window unconditionally every 100ms for the
// entire process lifetime, forcing continuous GPU redraws that leaked
// Metal-backed drawable memory over long uptimes (observed growing to tens
// of GB). An idle poll — same frame as last saved — must never redraw.
func TestShouldRedrawIdleTickDoesNotRedraw(t *testing.T) {
	frame := State{X: 10, Y: 20, Width: 800, Height: 600}
	if shouldRedraw(frame, frame) {
		t.Error("shouldRedraw(same, same) = true, want false: idle ticks must not force a redraw")
	}
}

func TestShouldRedrawOnFrameChange(t *testing.T) {
	last := State{X: 10, Y: 20, Width: 800, Height: 600}
	current := State{X: 15, Y: 20, Width: 800, Height: 600}
	if !shouldRedraw(current, last) {
		t.Error("shouldRedraw(moved, last) = false, want true: a real move/resize must redraw")
	}
}

func TestShouldRedrawInvalidFrameNeverRedraws(t *testing.T) {
	zero := State{}
	last := State{X: 10, Y: 20, Width: 800, Height: 600}
	if shouldRedraw(zero, last) {
		t.Error("shouldRedraw(invalid, last) = true, want false: an unset/zero-area frame must not redraw")
	}
}
