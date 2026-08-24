package persist

import (
	"log"
	"sync"
	"time"

	"gioui.org/app"
	"gioui.org/unit"
)

// Window wraps app.Window with automatic position and size persistence.
// Replace new(app.Window) with persist.NewWindow("appname") for one-line persistence.
type Window struct {
	*app.Window
	name  string
	saved State

	mu       sync.Mutex
	view     uintptr
	last     State
	restored bool
	done     chan struct{}
}

// NewWindow creates a Gio window that automatically persists its position and size.
// The name identifies this window's saved state (stored at ~/.config/daz-golang-gio/{name}.json).
// Pass any additional app.Option values after the name.
func NewWindow(name string, opts ...app.Option) *Window {
	saved, err := LoadState(name)
	if err != nil {
		log.Printf("persist: load state %q: %v", name, err)
	}

	w := &Window{
		Window: new(app.Window),
		name:   name,
		saved:  saved,
		done:   make(chan struct{}),
	}

	// Apply default size; position is restored later via native API.
	allOpts := []app.Option{app.Size(unit.Dp(800), unit.Dp(600))}
	allOpts = append(allOpts, opts...)
	w.Option(allOpts...)

	go w.tracker()

	return w
}

// Event returns the next window event, intercepting platform events
// to capture the native view handle. Use this instead of w.Window.Event().
func (w *Window) Event() any {
	e := w.Window.Event()
	w.handleEvent(e)
	return e
}

// Close stops the background tracker and does a final save.
// Call this when the window is destroyed.
func (w *Window) Close() {
	select {
	case <-w.done:
	default:
		close(w.done)
	}
}

// Frame returns the last known window frame. Safe to call from the event loop.
func (w *Window) Frame() State {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.last
}

func (w *Window) handleEvent(e any) {
	switch e.(type) {
	case app.DestroyEvent:
		w.Close()
	default:
		w.handlePlatformEvent(e)
	}
}

func (w *Window) setView(view uintptr) {
	w.mu.Lock()
	w.view = view
	w.mu.Unlock()
}

// tracker runs in a background goroutine, polling the native window frame
// and saving changes. All CGo and file I/O happens here, never in the event handler.
//
// It intentionally does NOT invalidate the Gio window on every tick. Gio
// redraws on macOS are driven directly by Invalidate (there is no
// CVDisplayLink to pace them — see gioui.org/app's os_macos.go), so an
// unconditional Invalidate here would force a full GPU redraw 10x/sec for
// the entire process lifetime, regardless of whether the window ever moved.
// Left running for hours that leaks Metal-backed drawable memory (observed
// growing to tens of GB in activity and agentd-gauge, both of which embed
// this package). Only invalidate when the polled frame actually changed.
func (w *Window) tracker() {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	var lastSaved State

	for {
		select {
		case <-w.done:
			w.finalSave()
			return
		case <-ticker.C:
			w.mu.Lock()
			view := w.view
			restored := w.restored
			w.mu.Unlock()

			if view == 0 {
				continue
			}

			// Restore saved position once.
			if !restored {
				w.restorePosition(view)
				w.mu.Lock()
				w.restored = true
				w.mu.Unlock()
			}

			// Read current frame.
			x, y, width, height := GetWindowFrame(view)
			current := State{X: x, Y: y, Width: width, Height: height}

			w.mu.Lock()
			w.last = current
			w.mu.Unlock()

			// Save and redraw only if the frame actually changed (e.g. the
			// user dragged or resized the window). Idle ticks must never
			// invalidate — see the tracker doc comment above.
			if shouldRedraw(current, lastSaved) {
				lastSaved = current
				if err := SaveState(w.name, current); err != nil {
					log.Printf("persist: save %q: %v", w.name, err)
				}
				w.Invalidate()
			}
		}
	}
}

// shouldRedraw reports whether a newly polled window frame differs from the
// last-saved frame and therefore warrants a save + Invalidate. Idle ticks
// (nothing moved, or the frame is not yet valid) must return false.
func shouldRedraw(current, lastSaved State) bool {
	return current.Valid() && !current.Equal(lastSaved)
}

func (w *Window) restorePosition(view uintptr) {
	if !w.saved.Valid() {
		return
	}
	if !PositionSupported() {
		return
	}
	if IsOnScreen(w.saved.X, w.saved.Y) {
		log.Printf("persist: restoring %q to (%.0f, %.0f) %.0fx%.0f",
			w.name, w.saved.X, w.saved.Y, w.saved.Width, w.saved.Height)
		SetWindowFrame(view, w.saved.X, w.saved.Y, w.saved.Width, w.saved.Height)
	} else {
		log.Printf("persist: saved position (%.0f, %.0f) is off-screen, skipping", w.saved.X, w.saved.Y)
	}
}

func (w *Window) finalSave() {
	w.mu.Lock()
	view := w.view
	w.mu.Unlock()
	if view == 0 || !PositionSupported() {
		return
	}
	x, y, width, height := GetWindowFrame(view)
	current := State{X: x, Y: y, Width: width, Height: height}
	if current.Valid() {
		if err := SaveState(w.name, current); err != nil {
			log.Printf("persist: final save %q: %v", w.name, err)
		}
	}
}
