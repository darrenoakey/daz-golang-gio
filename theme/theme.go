// Package theme provides predefined color palettes for Gio apps.
// Use Dark or Light as a starting point, then override individual colors.
package theme

import "image/color"

// Palette holds named colors for a Gio app's UI.
type Palette struct {
	// Backgrounds
	BG             color.NRGBA // Main window background
	Surface        color.NRGBA // Card/panel surfaces
	SurfaceAlt     color.NRGBA // Alternating row background
	HeaderBG       color.NRGBA // Table header background
	StatusBarBG    color.NRGBA // Status bar background
	SeparatorColor color.NRGBA // Row/section dividers

	// Text
	TextPrimary   color.NRGBA // High-emphasis text
	TextSecondary color.NRGBA // Medium-emphasis text
	TextMuted     color.NRGBA // Low-emphasis text (headers, labels)

	// Accents
	AccentBlue   color.NRGBA
	AccentGreen  color.NRGBA
	AccentOrange color.NRGBA
	AccentRed    color.NRGBA
	AccentPurple color.NRGBA
	AccentCyan   color.NRGBA
}

// Dark returns a dark theme palette matching the neon-terminal style
// used by activity monitor and spark view.
func Dark() Palette {
	return Palette{
		BG:             hex(0x0f0f0f),
		Surface:        hex(0x1a1a1a),
		SurfaceAlt:     hex(0x141414),
		HeaderBG:       hex(0x181818),
		StatusBarBG:    hex(0x121218),
		SeparatorColor: hex(0x2a2a2a),

		TextPrimary:   hex(0xe8e8e8),
		TextSecondary: hex(0xa8a8a8),
		TextMuted:     hex(0x606070),

		AccentBlue:   hex(0x5c9cff),
		AccentGreen:  hex(0x5cb85c),
		AccentOrange: hex(0xffb84d),
		AccentRed:    hex(0xff5c5c),
		AccentPurple: hex(0xb47aff),
		AccentCyan:   hex(0x00d4ff),
	}
}

// Light returns a Material Design 3 inspired light theme palette
// matching the auto-ps style.
func Light() Palette {
	return Palette{
		BG:             hex(0xFEF7FF),
		Surface:        hex(0xFFFFFF),
		SurfaceAlt:     hex(0xF7F2FA),
		HeaderBG:       hex(0xF3EDF7),
		StatusBarBG:    hex(0xF3EDF7),
		SeparatorColor: hex(0xCAC4D0),

		TextPrimary:   hex(0x1D1B20),
		TextSecondary: hex(0x49454F),
		TextMuted:     hex(0x79747E),

		AccentBlue:   hex(0x0061A4),
		AccentGreen:  hex(0x386A20),
		AccentOrange: hex(0x7D5700),
		AccentRed:    hex(0xB3261E),
		AccentPurple: hex(0x6750A4),
		AccentCyan:   hex(0x006A6A),
	}
}

// Hex creates a color from a 24-bit hex value (e.g., 0xFF5C5C).
func Hex(c uint32) color.NRGBA {
	return hex(c)
}

func hex(c uint32) color.NRGBA {
	return color.NRGBA{R: uint8(c >> 16), G: uint8(c >> 8), B: uint8(c), A: 0xFF}
}
