package table

import (
	"testing"

	"gioui.org/text"
	"gioui.org/unit"
)

func TestNewCreatesClickables(t *testing.T) {
	cols := []Column{
		{Label: "Name", Width: 0, Align: text.Start},
		{Label: "PID", Width: 80, Align: text.End},
		{Label: "CPU", Width: 60, Align: text.End},
	}
	tbl := New(cols, DefaultDarkStyle())
	if len(tbl.HeaderClicks) != 3 {
		t.Errorf("HeaderClicks len = %d, want 3", len(tbl.HeaderClicks))
	}
	if len(tbl.Columns) != 3 {
		t.Errorf("Columns len = %d, want 3", len(tbl.Columns))
	}
}

func TestDefaultDarkStyleValues(t *testing.T) {
	s := DefaultDarkStyle()
	if s.HeaderHeight != unit.Dp(32) {
		t.Errorf("HeaderHeight = %v, want 32", s.HeaderHeight)
	}
	if s.RowHeight != unit.Dp(36) {
		t.Errorf("RowHeight = %v, want 36", s.RowHeight)
	}
	if s.HeaderBG.A != 0xff {
		t.Error("HeaderBG should be opaque")
	}
}

func TestDefaultLightStyleValues(t *testing.T) {
	s := DefaultLightStyle()
	if s.HeaderBG.R < 0xF0 {
		t.Errorf("Light HeaderBG.R = %x, expected bright", s.HeaderBG.R)
	}
}
