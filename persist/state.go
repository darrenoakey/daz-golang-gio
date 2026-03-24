// Package persist provides automatic window position and size persistence for Gio apps.
// All platform-specific calls are thread-safe (dispatched to the main thread on macOS).
package persist

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// State holds the saved window geometry in native screen coordinates.
type State struct {
	X      float64 `json:"x"`
	Y      float64 `json:"y"`
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
}

// Valid reports whether the state has a positive area.
func (s State) Valid() bool {
	return s.Width > 0 && s.Height > 0
}

// Equal reports whether two states are identical.
func (s State) Equal(other State) bool {
	return s == other
}

// ConfigDir returns the directory where state files are stored.
func ConfigDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		home = "."
	}
	return filepath.Join(home, ".config", "daz-golang-gio")
}

// StatePath returns the full path to a state file for the given app name.
func StatePath(name string) string {
	return filepath.Join(ConfigDir(), name+".json")
}

// LoadState reads saved window state from disk.
// Returns a zero State and nil error if the file does not exist.
func LoadState(name string) (State, error) {
	path := StatePath(name)
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return State{}, nil
		}
		return State{}, fmt.Errorf("read state %s: %w", path, err)
	}
	var s State
	if err := json.Unmarshal(data, &s); err != nil {
		return State{}, fmt.Errorf("parse state %s: %w", path, err)
	}
	return s, nil
}

// SaveState writes window state to disk atomically (tmp + rename).
func SaveState(name string, s State) error {
	dir := ConfigDir()
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("create config dir: %w", err)
	}
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal state: %w", err)
	}
	path := StatePath(name)
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return fmt.Errorf("write state tmp: %w", err)
	}
	if err := os.Rename(tmp, path); err != nil {
		return fmt.Errorf("rename state: %w", err)
	}
	return nil
}
