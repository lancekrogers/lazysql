package whichkey

import (
	"testing"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func TestOverlayVisibility(t *testing.T) {
	overlay := NewOverlay(nil)
	if overlay.IsVisible() {
		t.Fatal("expected overlay to start hidden")
	}

	overlay.Show()
	if !overlay.IsVisible() {
		t.Fatal("expected overlay to be visible after Show")
	}

	overlay.Hide()
	if overlay.IsVisible() {
		t.Fatal("expected overlay to be hidden after Hide")
	}
}

func TestOverlayPositioning(t *testing.T) {
	overlay := NewOverlay(nil)
	overlay.SetContent([]KeyHint{{Key: "f", Description: "find"}})
	container := rect{x: 0, y: 0, width: 80, height: 24}
	width, height := overlay.calculateSize(container.width, container.height)

	overlay.SetPosition(PositionCenter)
	x, y := overlay.positionRect(container, width, height)
	if x != (container.width-width)/2 || y != (container.height-height)/2 {
		t.Fatalf("unexpected center position: got (%d,%d)", x, y)
	}

	overlay.SetPosition(PositionBottomRight)
	x, y = overlay.positionRect(container, width, height)
	if x != container.width-width-1 || y != container.height-height-1 {
		t.Fatalf("unexpected bottom-right position: got (%d,%d)", x, y)
	}
}

func TestOverlayDraw(t *testing.T) {
	screen := tcell.NewSimulationScreen("UTF-8")
	if err := screen.Init(); err != nil {
		t.Fatalf("init screen: %v", err)
	}
	screen.SetSize(80, 24)

	overlay := NewOverlay(nil)
	overlay.SetRect(0, 0, 80, 24)
	overlay.SetTitle(`\`)
	overlay.SetContent([]KeyHint{
		{Key: "f", Description: "find", IsGroup: true},
		{Key: "t", Description: "table"},
	})
	overlay.Show()
	overlay.Draw(screen)
}

func TestOverlayInputHandling(t *testing.T) {
	tree := NewKeyTree()
	tree.AddGroup([]rune{'f'}, "find")
	executed := false
	tree.AddCommand([]rune{'f', 't'}, "table", func() {
		executed = true
	})

	overlay := NewOverlay(nil)
	overlay.SetKeyTree(tree)
	overlay.Show()

	handler := overlay.InputHandler()
	handler(tcell.NewEventKey(tcell.KeyRune, 'f', tcell.ModNone), nil)
	if tree.Current == tree.Root {
		t.Fatal("expected to move into group on 'f'")
	}

	handler(tcell.NewEventKey(tcell.KeyRune, 't', tcell.ModNone), nil)
	if !executed {
		t.Fatal("expected action to execute on leaf")
	}
	if overlay.IsVisible() {
		t.Fatal("expected overlay to hide after execution")
	}

	overlay.Show()
	handler(tcell.NewEventKey(tcell.KeyRune, 'f', tcell.ModNone), nil)
	handler(tcell.NewEventKey(tcell.KeyEscape, 0, tcell.ModNone), nil)
	if overlay.IsVisible() && tree.Current != tree.Root {
		t.Fatal("expected escape to move back or hide overlay")
	}
}

func TestOverlayEmptyTreeSync(t *testing.T) {
	tree := NewKeyTree()
	overlay := NewOverlay(nil)
	overlay.SetKeyTree(tree)
	overlay.Show()

	if len(overlay.content) == 0 || overlay.content[0].Description != "No mappings" {
		t.Fatal("expected empty tree to render no mappings message")
	}
	if overlay.title != `\ > ` {
		t.Fatalf("unexpected breadcrumb for root: %q", overlay.title)
	}
}

func TestOverlayFocusReturn(t *testing.T) {
	overlay := NewOverlay(nil)
	returnFocus := tview.NewBox()

	var gotFocus tview.Primitive
	overlay.SetFocusSetter(func(p tview.Primitive) {
		gotFocus = p
	})
	overlay.SetReturnFocus(returnFocus)

	overlay.Show()
	if gotFocus != overlay {
		t.Fatal("expected overlay to receive focus on show")
	}

	overlay.Hide()
	if gotFocus != returnFocus {
		t.Fatal("expected return focus after hide")
	}
}

func TestOverlaySyncNoTree(t *testing.T) {
	overlay := NewOverlay(nil)
	overlay.SetLeaderPrefix(';')
	overlay.syncFromTree()

	if len(overlay.content) == 0 || overlay.content[0].Description != "No mappings" {
		t.Fatal("expected no tree to render no mappings message")
	}
	if overlay.title != ";" {
		t.Fatalf("expected title to use leader prefix, got %q", overlay.title)
	}
}

func TestOverlayPositionClamp(t *testing.T) {
	overlay := NewOverlay(nil)
	overlay.SetContent([]KeyHint{{Key: "longkey", Description: "long description"}})

	container := rect{x: 0, y: 0, width: 10, height: 5}
	width, height := overlay.calculateSize(container.width, container.height)
	if width > container.width || height > container.height {
		t.Fatal("expected size to fit within container")
	}

	overlay.SetPosition(PositionBottomRight)
	x, y := overlay.positionRect(container, width, height)
	if x < container.x || y < container.y {
		t.Fatalf("expected position within container, got (%d,%d)", x, y)
	}
}

func TestOverlayRapidToggle(t *testing.T) {
	overlay := NewOverlay(nil)
	for i := 0; i < 10; i++ {
		overlay.Show()
		overlay.Hide()
	}
	if overlay.IsVisible() {
		t.Fatal("expected overlay to remain hidden after rapid toggle")
	}
}
