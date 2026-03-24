package menu

import (
	"image"
	"image/color"
	"testing"
)

var testItems = []Item{
	{Label: "Cut"},
	{Label: "Copy"},
	{Label: "Paste"},
	{Label: "Delete", Color: color.NRGBA{R: 0xff, G: 0x5c, B: 0x5c, A: 0xff}},
}

// --- Lifecycle: Show / Dismiss ---

func TestContextMenu_ShowMakesVisible(t *testing.T) {
	var m ContextMenu
	if m.Visible() {
		t.Fatal("new menu should not be visible")
	}
	m.cursorPos = image.Pt(100, 200)
	m.Show(testItems)
	if !m.Visible() {
		t.Error("menu should be visible after Show")
	}
}

func TestContextMenu_ShowUsesCursorPos(t *testing.T) {
	var m ContextMenu
	m.cursorPos = image.Pt(100, 200)
	m.Show(testItems)
	if m.pos != (image.Point{X: 100, Y: 200}) {
		t.Errorf("pos = %v, want (100,200)", m.pos)
	}
}

func TestContextMenu_ShowStoresItems(t *testing.T) {
	var m ContextMenu
	m.Show(testItems)
	if len(m.items) != len(testItems) {
		t.Fatalf("items count = %d, want %d", len(m.items), len(testItems))
	}
	for i, item := range m.items {
		if item.Label != testItems[i].Label {
			t.Errorf("item[%d].Label = %q, want %q", i, item.Label, testItems[i].Label)
		}
	}
}

func TestContextMenu_ShowResetsHover(t *testing.T) {
	var m ContextMenu
	m.hoverIdx = 2
	m.Show(testItems)
	if m.hoverIdx != -1 {
		t.Errorf("hoverIdx = %d, want -1 after Show", m.hoverIdx)
	}
}

func TestContextMenu_ShowSetsShowFrame(t *testing.T) {
	var m ContextMenu
	m.Show(testItems)
	if !m.showFrame {
		t.Error("showFrame should be true after Show")
	}
}

func TestContextMenu_DismissClearsVisible(t *testing.T) {
	var m ContextMenu
	m.cursorPos = image.Pt(100, 200)
	m.Show(testItems)
	m.Dismiss()
	if m.Visible() {
		t.Error("menu should not be visible after Dismiss")
	}
}

func TestContextMenu_DismissResetsHover(t *testing.T) {
	var m ContextMenu
	m.Show(testItems)
	m.hoverIdx = 2
	m.Dismiss()
	if m.hoverIdx != -1 {
		t.Errorf("hoverIdx = %d, want -1 after Dismiss", m.hoverIdx)
	}
}

func TestContextMenu_ShowUpdatesPosition(t *testing.T) {
	var m ContextMenu
	m.cursorPos = image.Pt(100, 100)
	m.Show(testItems[:2])
	m.cursorPos = image.Pt(200, 200)
	m.Show(testItems)
	if m.pos != (image.Point{X: 200, Y: 200}) {
		t.Errorf("pos = %v, want (200,200)", m.pos)
	}
	if len(m.items) != len(testItems) {
		t.Errorf("items count = %d, want %d after second Show", len(m.items), len(testItems))
	}
}

// --- ShowFrame bug fix: Show-Dismiss-Show cycle ---

func TestContextMenu_ShowDismissShowCycle(t *testing.T) {
	var m ContextMenu

	m.cursorPos = image.Pt(100, 100)
	m.Show(testItems)
	if !m.Visible() || !m.showFrame {
		t.Fatal("first Show failed")
	}

	m.showFrame = false // simulate frame processing
	m.Dismiss()

	m.cursorPos = image.Pt(200, 200)
	m.Show(testItems)
	if !m.Visible() {
		t.Error("should be visible after second Show")
	}
	if !m.showFrame {
		t.Error("showFrame must be true after second Show (prevents stale bg dismiss)")
	}
}

// --- Position clamping ---

func TestClampPosition_WithinBounds(t *testing.T) {
	pos := ClampPosition(image.Pt(50, 50), 120, 100, 800, 600)
	if pos != (image.Point{X: 50, Y: 50}) {
		t.Errorf("pos = %v, want (50,50)", pos)
	}
}

func TestClampPosition_RightOverflow(t *testing.T) {
	pos := ClampPosition(image.Pt(750, 50), 120, 100, 800, 600)
	if pos.X != 680 {
		t.Errorf("pos.X = %d, want 680", pos.X)
	}
}

func TestClampPosition_BottomOverflow(t *testing.T) {
	pos := ClampPosition(image.Pt(50, 550), 120, 100, 800, 600)
	if pos.Y != 500 {
		t.Errorf("pos.Y = %d, want 500", pos.Y)
	}
}

func TestClampPosition_NegativeCoords(t *testing.T) {
	pos := ClampPosition(image.Pt(-10, -20), 120, 100, 800, 600)
	if pos.X != 0 || pos.Y != 0 {
		t.Errorf("pos = %v, want (0,0)", pos)
	}
}

func TestClampPosition_CornerOverflow(t *testing.T) {
	pos := ClampPosition(image.Pt(750, 550), 120, 100, 800, 600)
	if pos.X != 680 || pos.Y != 500 {
		t.Errorf("pos = %v, want (680,500)", pos)
	}
}

func TestClampPosition_ExactFit(t *testing.T) {
	pos := ClampPosition(image.Pt(0, 0), 800, 600, 800, 600)
	if pos.X != 0 || pos.Y != 0 {
		t.Errorf("pos = %v, want (0,0)", pos)
	}
}

func TestClampPosition_LargerThanWindow(t *testing.T) {
	pos := ClampPosition(image.Pt(100, 100), 900, 700, 800, 600)
	if pos.X != 0 || pos.Y != 0 {
		t.Errorf("pos = %v, want (0,0) when menu exceeds window", pos)
	}
}

// --- Item colors ---

func TestItem_DefaultColor(t *testing.T) {
	item := Item{Label: "Normal"}
	if item.Color.A != 0 {
		t.Error("zero-value Item.Color should have alpha 0 (use default)")
	}
}

func TestItem_CustomColor(t *testing.T) {
	red := color.NRGBA{R: 0xff, A: 0xff}
	item := Item{Label: "Delete", Color: red}
	if item.Color != red {
		t.Errorf("Item.Color = %v, want %v", item.Color, red)
	}
}

// --- Tag management ---

func TestContextMenu_EnsureTagsAllocates(t *testing.T) {
	var m ContextMenu
	m.items = testItems
	m.ensureTags()
	if len(m.itemTags) < len(m.items) {
		t.Errorf("itemTags len = %d, want >= %d", len(m.itemTags), len(m.items))
	}
}

func TestContextMenu_EnsureTagsStable(t *testing.T) {
	var m ContextMenu
	m.items = testItems
	m.ensureTags()
	first := make([]*bool, len(m.itemTags))
	copy(first, m.itemTags)

	m.ensureTags()
	for i, tag := range m.itemTags[:len(first)] {
		if tag != first[i] {
			t.Errorf("tag[%d] pointer changed after second ensureTags", i)
		}
	}
}

func TestContextMenu_EnsureTagsUnique(t *testing.T) {
	var m ContextMenu
	m.items = testItems
	m.ensureTags()
	seen := make(map[*bool]bool)
	for _, tag := range m.itemTags {
		if seen[tag] {
			t.Error("duplicate tag pointer found")
		}
		seen[tag] = true
	}
}

// --- Menu dimensions ---

func TestMenuDimensions(t *testing.T) {
	items := testItems
	expectedH := ItemHeight*len(items) + PadTop + PadBottom
	if expectedH != 136 { // 32*4 + 4 + 4
		t.Errorf("total menu height = %d dp, want 136", expectedH)
	}
	if Width != 120 {
		t.Errorf("menu width = %d dp, want 120", Width)
	}
}

// --- Hover state ---

func TestContextMenu_HoverInitialState(t *testing.T) {
	var m ContextMenu
	m.Show(testItems)
	if m.hoverIdx != -1 {
		t.Errorf("hoverIdx = %d, want -1 (no hover initially)", m.hoverIdx)
	}
}

func TestContextMenu_HoverClearedOnReshow(t *testing.T) {
	var m ContextMenu
	m.Show(testItems)
	m.hoverIdx = 2
	m.cursorPos = image.Pt(50, 50)
	m.Show(testItems)
	if m.hoverIdx != -1 {
		t.Errorf("hoverIdx = %d, want -1 after re-Show", m.hoverIdx)
	}
}

func TestContextMenu_HoverClearedOnDismiss(t *testing.T) {
	var m ContextMenu
	m.Show(testItems)
	m.hoverIdx = 1
	m.Dismiss()
	if m.hoverIdx != -1 {
		t.Errorf("hoverIdx = %d, want -1 after Dismiss", m.hoverIdx)
	}
}

// --- Cursor tracking ---

func TestContextMenu_CursorPosUsedByShow(t *testing.T) {
	var m ContextMenu
	m.cursorPos = image.Pt(300, 400)
	m.Show(testItems)
	if m.pos != (image.Point{X: 300, Y: 400}) {
		t.Errorf("pos = %v, want (300,400) from cursorPos", m.pos)
	}
}

func TestContextMenu_CursorPosZeroDefault(t *testing.T) {
	var m ContextMenu
	m.Show(testItems)
	if m.pos != (image.Point{}) {
		t.Errorf("pos = %v, want (0,0) when cursorPos not set", m.pos)
	}
}
