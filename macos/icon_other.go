//go:build !darwin || ios

// Package macos provides macOS-specific utilities for Gio apps.
package macos

// SetDockIcon is a no-op on non-macOS platforms.
func SetDockIcon(data []byte) {}
