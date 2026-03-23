package persist

import (
	"log"
	"sync"
	"time"

	"gioui.org/app"
	"gioui.org/unit"
)

const saveInterval = 30 * time.Second
const saveDebounceDuration = 500 * time.Millisecond

// Window wraps app.Window with automatic position and size persistence.
// Replace new(app.Window) with persist.NewWindow("appname") for one-line persistence.
type Window struct {
	*app.Window
	name         string
	state        State
	view         uintptr
	dirty        bool
	mu           sync.Mutex
	debounce     *time.Timer
	periodicDone chan struct{}
	restored     bool
}

// NewWindow creates a Gio window that automatically persists its position and size.
// The name identifies this window's saved state (stored at ~/.config/daz-golang-gio/{name}.json).
// Pass any additional app.Option values after the name.
func NewWindow(name string, opts ...app.Option) *Window {
	saved, err := LoadState(name)
	if err != nil {
		log.Printf("persist: failed to load state for %q: %v", name, err)
	}

	var applyOpts []app.Option
	if saved.Width > 0 && saved.Height > 0 && saved.Mode != "Maximized" {
		applyOpts = append(applyOpts, app.Size(unit.Dp(saved.Width), unit.Dp(saved.Height)))
	}
	applyOpts = append(applyOpts, opts...)

	w := &Window{
		Window:       new(app.Window),
		name:         name,
		state:        saved,
		periodicDone: make(chan struct{}),
	}
	w.Window.Option(applyOpts...)

	go w.periodicSave()

	return w
}

// Event returns the next window event, intercepting config and destroy events
// to track window geometry. Use this in your event loop instead of w.Window.Event().
func (w *Window) Event() any {
	e := w.Window.Event()
	w.handleEvent(e)
	return e
}

// Close flushes any pending state to disk and stops the periodic saver.
// Call this when the window is destroyed (e.g., defer w.Close()).
func (w *Window) Close() {
	close(w.periodicDone)
	w.mu.Lock()
	if w.debounce != nil {
		w.debounce.Stop()
	}
	w.mu.Unlock()
	w.flushState()
}

func (w *Window) handleEvent(e any) {
	switch e := e.(type) {
	case app.ConfigEvent:
		w.mu.Lock()
		cfg := e.Config
		if cfg.Size.X > 0 && cfg.Size.Y > 0 {
			w.state.Width = cfg.Size.X
			w.state.Height = cfg.Size.Y
		}
		w.state.Mode = cfg.Mode.String()
		w.dirty = true
		w.scheduleSave()
		w.mu.Unlock()

	case app.DestroyEvent:
		w.capturePosition()
		w.flushState()

	default:
		w.handlePlatformEvent(e)
	}
}

func (w *Window) capturePosition() {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.view == 0 || !PositionSupported() {
		return
	}
	x, y, width, height := GetWindowPosition(w.view)
	if width > 0 && height > 0 {
		w.state.X = x
		w.state.Y = y
		w.state.Width = width
		w.state.Height = height
		w.dirty = true
	}
}

func (w *Window) scheduleSave() {
	if w.debounce != nil {
		w.debounce.Stop()
	}
	w.debounce = time.AfterFunc(saveDebounceDuration, func() {
		w.capturePosition()
		w.flushState()
	})
}

func (w *Window) flushState() {
	w.mu.Lock()
	state := w.state
	dirty := w.dirty
	w.dirty = false
	w.mu.Unlock()

	if !dirty {
		return
	}
	if err := SaveState(w.name, state); err != nil {
		log.Printf("persist: failed to save state for %q: %v", w.name, err)
	}
}

func (w *Window) periodicSave() {
	ticker := time.NewTicker(saveInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			w.capturePosition()
			w.flushState()
		case <-w.periodicDone:
			return
		}
	}
}

// SetView stores the native view handle for position tracking.
// Called by platform-specific event handlers.
func (w *Window) SetView(view uintptr) {
	w.mu.Lock()
	w.view = view
	w.mu.Unlock()
}

// RestorePosition applies saved position if not already restored.
// Called by platform-specific event handlers when the native view is available.
func (w *Window) RestorePosition(view uintptr) {
	w.mu.Lock()
	if w.restored {
		w.mu.Unlock()
		return
	}
	w.restored = true
	state := w.state
	w.mu.Unlock()

	if state.X != 0 || state.Y != 0 {
		SetWindowPosition(view, state.X, state.Y, state.Width, state.Height)
	}
}
