package theme

import "testing"

func TestDarkPalette(t *testing.T) {
	p := Dark()
	if p.BG.A != 0xff {
		t.Error("Dark BG alpha should be 0xff")
	}
	if p.TextPrimary.R != 0xe8 {
		t.Errorf("Dark TextPrimary.R = %x, want 0xe8", p.TextPrimary.R)
	}
	if p.AccentBlue == p.AccentRed {
		t.Error("AccentBlue and AccentRed should differ")
	}
}

func TestLightPalette(t *testing.T) {
	p := Light()
	if p.BG.R < 0xF0 {
		t.Errorf("Light BG.R = %x, expected a bright value", p.BG.R)
	}
	if p.TextPrimary.R > 0x30 {
		t.Errorf("Light TextPrimary.R = %x, expected a dark value", p.TextPrimary.R)
	}
}

func TestHex(t *testing.T) {
	c := Hex(0xFF5C5C)
	if c.R != 0xFF || c.G != 0x5C || c.B != 0x5C || c.A != 0xFF {
		t.Errorf("Hex(0xFF5C5C) = %+v, want {255 92 92 255}", c)
	}
}

func TestHexBlack(t *testing.T) {
	c := Hex(0x000000)
	if c.R != 0 || c.G != 0 || c.B != 0 || c.A != 0xFF {
		t.Errorf("Hex(0x000000) = %+v, want {0 0 0 255}", c)
	}
}
