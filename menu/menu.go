// Package menu provides a reusable floating context menu component for Gio
// with hover highlighting, click-outside dismiss, and HiDPI support.
package menu

import (
	"image"
	"image/color"

	"gioui.org/font"
	"gioui.org/io/event"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget/material"
)

// Item defines a single context menu entry.
type Item struct {
	Label string
	Color color.NRGBA // per-item text color; zero value uses default
}

// Result is returned from Layout when an item is clicked.
type Result struct {
	Index int  // index of the clicked item
	OK    bool // true if an item was selected
}

// Dimensions in dp.
const (
	ItemHeight = 32
	Width      = 120
	PadTop     = 4
	PadBottom  = 4
)

// Default colors (dark theme).
var (
	DefaultBG      = color.NRGBA{R: 0x28, G: 0x28, B: 0x28, A: 0xf0}
	DefaultBorder  = color.NRGBA{R: 0x3a, G: 0x3a, B: 0x3a, A: 0xff}
	DefaultHoverBG = color.NRGBA{R: 0x24, G: 0x3e, B: 0x6c, A: 0xff}
	DefaultText    = color.NRGBA{R: 0xe8, G: 0xe8, B: 0xe8, A: 0xff}
	DefaultHover   = color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff}
)

// ContextMenu is a reusable floating context menu with hover highlighting.
//
// The dismiss overlay uses pointer.PassOp so events pass through to handlers
// underneath (e.g. row click handlers). A showFrame flag prevents stale events
// from dismissing the menu on the same frame Show() was called.
//
// All dimensions are specified in dp and converted to pixels at render time
// via gtx.Dp(), ensuring correct sizing on HiDPI displays.
type ContextMenu struct {
	visible   bool
	showFrame bool        // skip bg dismiss on the frame Show() was called
	pos       image.Point // absolute window position (pixels)
	items     []Item

	itemTags []*bool // stable pointer event tags, one per item
	bgTag    bool    // background dismiss tag

	hoverIdx int // currently hovered item index, -1 = none
}

// Show opens the context menu at the given window position with the given items.
func (m *ContextMenu) Show(pos image.Point, items []Item) {
	m.visible = true
	m.showFrame = true
	m.pos = pos
	m.items = items
	m.hoverIdx = -1
	m.ensureTags()
}

// Dismiss closes the context menu.
func (m *ContextMenu) Dismiss() {
	m.visible = false
	m.hoverIdx = -1
}

// Visible returns whether the context menu is currently shown.
func (m *ContextMenu) Visible() bool {
	return m.visible
}

// ensureTags allocates stable pointer tags for each item.
func (m *ContextMenu) ensureTags() {
	for len(m.itemTags) < len(m.items) {
		tag := new(bool)
		m.itemTags = append(m.itemTags, tag)
	}
}

// ClampPosition adjusts a menu position to stay within window bounds.
func ClampPosition(pos image.Point, menuW, menuH, winW, winH int) image.Point {
	if maxX := winW - menuW; pos.X > maxX {
		pos.X = maxX
	}
	if maxY := winH - menuH; pos.Y > maxY {
		pos.Y = maxY
	}
	if pos.X < 0 {
		pos.X = 0
	}
	if pos.Y < 0 {
		pos.Y = 0
	}
	return pos
}

// drainEvents discards any queued pointer events for all menu tags.
func (m *ContextMenu) drainEvents(gtx layout.Context) {
	for {
		if _, ok := gtx.Event(pointer.Filter{Target: &m.bgTag, Kinds: pointer.Press}); !ok {
			break
		}
	}
	for _, tag := range m.itemTags {
		for {
			if _, ok := gtx.Event(pointer.Filter{
				Target: tag,
				Kinds:  pointer.Press | pointer.Enter | pointer.Leave,
			}); !ok {
				break
			}
		}
	}
}

// Layout renders the context menu overlay and returns any triggered action.
// Call this AFTER laying out the main content so the menu renders on top,
// but row/button handlers underneath still receive events via PassOp.
func (m *ContextMenu) Layout(gtx layout.Context, th *material.Theme) Result {
	if !m.visible {
		m.drainEvents(gtx)
		return Result{}
	}

	// On the frame Show() was called, drain stale bg events to prevent
	// immediate dismiss from events queued before Show().
	if m.showFrame {
		m.showFrame = false
		m.drainEvents(gtx)
	}

	// Convert dp constants to pixels for current display density
	itemH := gtx.Dp(ItemHeight)
	w := gtx.Dp(Width)
	padTop := gtx.Dp(PadTop)
	padBot := gtx.Dp(PadBottom)
	totalH := itemH*len(m.items) + padTop + padBot

	// Compute display position: center first item at cursor, shift left slightly
	displayPos := ClampPosition(
		image.Pt(m.pos.X-gtx.Dp(8), m.pos.Y-itemH/2-padTop),
		w, totalH, gtx.Constraints.Max.X, gtx.Constraints.Max.Y,
	)

	// Background dismiss area: full window, with pass-through so
	// handlers underneath still receive pointer events.
	passStack := pointer.PassOp{}.Push(gtx.Ops)
	bgArea := clip.Rect{Max: gtx.Constraints.Max}.Push(gtx.Ops)
	event.Op(gtx.Ops, &m.bgTag)
	bgArea.Pop()
	passStack.Pop()

	// Check for dismiss click (press outside menu bounds)
	for {
		ev, ok := gtx.Event(pointer.Filter{
			Target: &m.bgTag,
			Kinds:  pointer.Press,
		})
		if !ok {
			break
		}
		if e, ok := ev.(pointer.Event); ok {
			menuRect := image.Rect(displayPos.X, displayPos.Y,
				displayPos.X+w, displayPos.Y+totalH)
			pt := image.Pt(int(e.Position.X), int(e.Position.Y))
			if !pt.In(menuRect) {
				m.Dismiss()
				return Result{}
			}
		}
	}

	pos := displayPos

	// Draw menu at position
	menuOffset := op.Offset(pos).Push(gtx.Ops)

	// Border
	borderRect := clip.Rect{Max: image.Pt(w, totalH)}.Push(gtx.Ops)
	paint.FillShape(gtx.Ops, DefaultBorder, clip.Rect{Max: image.Pt(w, totalH)}.Op())
	borderRect.Pop()

	// Inner background
	innerRect := clip.Rect{Min: image.Pt(1, 1), Max: image.Pt(w-1, totalH-1)}.Push(gtx.Ops)
	paint.FillShape(gtx.Ops, DefaultBG, clip.Rect{Max: image.Pt(w-2, totalH-2)}.Op())
	innerRect.Pop()

	// Draw items
	var result Result
	y := padTop
	labelPadX := gtx.Dp(12)
	labelPadY := gtx.Dp(7)
	m.ensureTags()
	for i, item := range m.items {
		tag := m.itemTags[i]
		itemOff := op.Offset(image.Pt(1, y)).Push(gtx.Ops)
		iw := w - 2

		// Event area for click + hover
		itemArea := clip.Rect{Max: image.Pt(iw, itemH)}.Push(gtx.Ops)
		event.Op(gtx.Ops, tag)
		itemArea.Pop()

		// Process events (click + hover)
		for {
			ev, ok := gtx.Event(pointer.Filter{
				Target: tag,
				Kinds:  pointer.Press | pointer.Enter | pointer.Leave,
			})
			if !ok {
				break
			}
			if e, ok := ev.(pointer.Event); ok {
				switch e.Kind {
				case pointer.Enter:
					m.hoverIdx = i
				case pointer.Leave:
					if m.hoverIdx == i {
						m.hoverIdx = -1
					}
				case pointer.Press:
					result = Result{
						Index: i,
						OK:    true,
					}
					m.Dismiss()
				}
			}
		}

		// Hover background
		if m.hoverIdx == i {
			paint.FillShape(gtx.Ops, DefaultHoverBG,
				clip.Rect{Max: image.Pt(iw, itemH)}.Op())
		}

		// Label
		labelOff := op.Offset(image.Pt(labelPadX, labelPadY)).Push(gtx.Ops)

		labelColor := DefaultText
		if item.Color.A > 0 {
			labelColor = item.Color
		}
		if m.hoverIdx == i {
			labelColor = DefaultHover
		}

		l := material.Body2(th, item.Label)
		l.Color = labelColor
		l.TextSize = unit.Sp(13)
		l.Font.Weight = font.Normal
		l.Alignment = text.Start

		gtxLabel := gtx
		gtxLabel.Constraints = layout.Exact(image.Pt(iw-labelPadX*2, itemH-labelPadY*2))
		l.Layout(gtxLabel)
		labelOff.Pop()

		itemOff.Pop()
		y += itemH
	}

	menuOffset.Pop()

	return result
}
