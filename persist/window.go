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
	w.Window.Option(allOpts...)

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

			// Save if changed.
			if current.Valid() && !current.Equal(lastSaved) {
				lastSaved = current
				if err := SaveState(w.name, current); err != nil {
					log.Printf("persist: save %q: %v", w.name, err)
				}
			}

			w.Window.Invalidate()
		}
	}
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
