// Package table provides a column-based table layout widget for Gio.
// Supports fixed-width and flexible columns, alternating row backgrounds,
// row separators, and clickable headers.
package table

import (
	"image"
	"image/color"

	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

// Column defines a table column.
type Column struct {
	Label string
	Width unit.Dp        // Fixed width in Dp. 0 means flex (fills remaining space).
	Align text.Alignment // Text alignment within the column.
}

// Style controls the visual appearance of a table.
type Style struct {
	HeaderHeight unit.Dp
	RowHeight    unit.Dp
	HeaderBG     color.NRGBA
	RowAltBG     color.NRGBA // Alternating row background. Zero value = no alternation.
	SeparatorClr color.NRGBA // Row separator color. Zero value = no separators.
	TextSize     unit.Sp
	HeaderColor  color.NRGBA // Header label color.
	PadH         unit.Dp     // Horizontal padding inside each cell.
}

// DefaultDarkStyle returns a table style matching the dark neon theme.
func DefaultDarkStyle() Style {
	return Style{
		HeaderHeight: 32,
		RowHeight:    36,
		HeaderBG:     color.NRGBA{R: 0x18, G: 0x18, B: 0x18, A: 0xff},
		RowAltBG:     color.NRGBA{R: 0x14, G: 0x14, B: 0x14, A: 0xff},
		SeparatorClr: color.NRGBA{R: 0x22, G: 0x22, B: 0x22, A: 0xff},
		TextSize:     12,
		HeaderColor:  color.NRGBA{R: 0x60, G: 0x60, B: 0x70, A: 0xff},
		PadH:         16,
	}
}

// DefaultLightStyle returns a table style matching the Material Design 3 light theme.
func DefaultLightStyle() Style {
	return Style{
		HeaderHeight: 32,
		RowHeight:    36,
		HeaderBG:     color.NRGBA{R: 0xF3, G: 0xED, B: 0xF7, A: 0xff},
		RowAltBG:     color.NRGBA{R: 0xF7, G: 0xF2, B: 0xFA, A: 0xff},
		SeparatorClr: color.NRGBA{R: 0xCA, G: 0xC4, B: 0xD0, A: 0xff},
		TextSize:     12,
		HeaderColor:  color.NRGBA{R: 0x49, G: 0x45, B: 0x4F, A: 0xff},
		PadH:         16,
	}
}

// Cell describes the content and style of a single table cell.
type Cell struct {
	Text  string
	Color color.NRGBA
	Bold  bool
}

// Table holds the state for a column-based table.
type Table struct {
	Columns      []Column
	HeaderClicks []widget.Clickable // One per column (optional, for sortable headers).
	style        Style
}

// New creates a Table with the given columns and style.
func New(columns []Column, style Style) *Table {
	clicks := make([]widget.Clickable, len(columns))
	return &Table{
		Columns:      columns,
		HeaderClicks: clicks,
		style:        style,
	}
}

// columnWidths computes pixel widths for each column, flexing the first
// column with Width==0 to fill remaining space.
func (t *Table) columnWidths(gtx layout.Context) []int {
	totalW := gtx.Constraints.Max.X
	widths := make([]int, len(t.Columns))
	flexIdx := -1
	used := 0
	for i, col := range t.Columns {
		if col.Width == 0 {
			flexIdx = i
		} else {
			widths[i] = gtx.Dp(col.Width)
			used += widths[i]
		}
	}
	if flexIdx >= 0 {
		remaining := totalW - used
		if remaining < 0 {
			remaining = 0
		}
		widths[flexIdx] = remaining
	}
	return widths
}

// LayoutHeader draws the table header row. Returns which column index
// was clicked (-1 if none).
func (t *Table) LayoutHeader(gtx layout.Context, th *material.Theme) (dims layout.Dimensions, clicked int) {
	clicked = -1
	headerH := gtx.Dp(t.style.HeaderHeight)
	totalW := gtx.Constraints.Max.X
	widths := t.columnWidths(gtx)

	// Header background.
	paint.FillShape(gtx.Ops, t.style.HeaderBG, clip.Rect{Max: image.Pt(totalW, headerH)}.Op())

	x := 0
	for i, col := range t.Columns {
		colW := widths[i]
		offset := op.Offset(image.Pt(x, 0)).Push(gtx.Ops)
		gtxCol := gtx
		gtxCol.Constraints = layout.Exact(image.Pt(colW, headerH))

		// Check clicks.
		if i < len(t.HeaderClicks) {
			for t.HeaderClicks[i].Clicked(gtxCol) {
				clicked = i
			}
		}

		layout.Inset{
			Left: t.style.PadH, Right: t.style.PadH,
			Top: unit.Dp(8),
		}.Layout(gtxCol, func(gtx layout.Context) layout.Dimensions {
			l := material.Body2(th, col.Label)
			l.Font.Weight = font.Medium
			l.TextSize = unit.Sp(11)
			l.Color = t.style.HeaderColor
			l.Alignment = col.Align
			l.MaxLines = 1
			return l.Layout(gtx)
		})

		// Clickable area.
		if i < len(t.HeaderClicks) {
			t.HeaderClicks[i].Layout(gtxCol, func(gtx layout.Context) layout.Dimensions {
				return layout.Dimensions{Size: image.Pt(colW, headerH)}
			})
		}

		offset.Pop()
		x += colW
	}

	return layout.Dimensions{Size: image.Pt(totalW, headerH)}, clicked
}

// LayoutRow draws a single data row. The cells slice must match the number of columns.
func (t *Table) LayoutRow(gtx layout.Context, th *material.Theme, index int, cells []Cell) layout.Dimensions {
	return t.LayoutRowWithBG(gtx, th, index, cells, color.NRGBA{})
}

// LayoutRowWithBG draws a data row with an optional custom background color.
// If bg is zero-value, the default alternating pattern is used.
func (t *Table) LayoutRowWithBG(gtx layout.Context, th *material.Theme, index int, cells []Cell, bg color.NRGBA) layout.Dimensions {
	rowH := gtx.Dp(t.style.RowHeight)
	totalW := gtx.Constraints.Max.X
	widths := t.columnWidths(gtx)

	// Row background.
	if bg.A > 0 {
		paint.FillShape(gtx.Ops, bg, clip.Rect{Max: image.Pt(totalW, rowH)}.Op())
	} else if t.style.RowAltBG.A > 0 && index%2 == 0 {
		paint.FillShape(gtx.Ops, t.style.RowAltBG, clip.Rect{Max: image.Pt(totalW, rowH)}.Op())
	}

	x := 0
	for i, cell := range cells {
		if i >= len(widths) {
			break
		}
		colW := widths[i]
		offset := op.Offset(image.Pt(x, 0)).Push(gtx.Ops)
		gtxCol := gtx
		gtxCol.Constraints = layout.Exact(image.Pt(colW, rowH))

		layout.Inset{
			Left: t.style.PadH, Right: t.style.PadH,
			Top: unit.Dp(10),
		}.Layout(gtxCol, func(gtx layout.Context) layout.Dimensions {
			l := material.Body2(th, cell.Text)
			l.Color = cell.Color
			l.TextSize = t.style.TextSize
			l.MaxLines = 1
			if i < len(t.Columns) {
				l.Alignment = t.Columns[i].Align
			}
			if cell.Bold {
				l.Font.Weight = font.Medium
			}
			return l.Layout(gtx)
		})

		offset.Pop()
		x += colW
	}

	// Row separator.
	if t.style.SeparatorClr.A > 0 {
		pad := gtx.Dp(t.style.PadH)
		sepOff := op.Offset(image.Pt(pad, rowH-1)).Push(gtx.Ops)
		sepW := totalW - pad*2
		if sepW > 0 {
			paint.FillShape(gtx.Ops, t.style.SeparatorClr, clip.Rect{Max: image.Pt(sepW, 1)}.Op())
		}
		sepOff.Pop()
	}

	return layout.Dimensions{Size: image.Pt(totalW, rowH)}
}
