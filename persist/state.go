// Package persist provides automatic window position and size persistence for Gio apps.
package persist

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// State holds the saved window geometry.
type State struct {
	Width  int    `json:"width"`
	Height int    `json:"height"`
	X      int    `json:"x"`
	Y      int    `json:"y"`
	Mode   string `json:"mode"`
}

// ConfigDir returns the directory where state files are stored.
// Uses ~/.config/daz-golang-gio/ on all platforms.
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
		return State{}, err
	}
	var s State
	if err := json.Unmarshal(data, &s); err != nil {
		return State{}, err
	}
	return s, nil
}

// SaveState writes window state to disk, creating the config directory if needed.
func SaveState(name string, s State) error {
	dir := ConfigDir()
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	data, err := json.Marshal(s)
	if err != nil {
		return err
	}
	return os.WriteFile(StatePath(name), data, 0644)
}
