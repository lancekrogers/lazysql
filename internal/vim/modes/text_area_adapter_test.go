package modes

import (
	"strings"
	"testing"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func setupTextArea(t *testing.T, text string) (*tview.TextArea, tcell.Screen) {
	t.Helper()

	area := tview.NewTextArea()
	area.SetText(text, false)
	area.SetRect(0, 0, 40, 10)

	screen := tcell.NewSimulationScreen("UTF-8")
	if err := screen.Init(); err != nil {
		t.Fatalf("init screen: %v", err)
	}
	screen.SetSize(40, 10)
	area.Draw(screen)

	return area, screen
}

func TestTextAreaAdapterMovement(t *testing.T) {
	area, _ := setupTextArea(t, "hello world")
	adapter := NewTextAreaAdapter(area)

	_, colBefore := adapter.cursorPos()
	adapter.MoveRight()
	_, colAfter := adapter.cursorPos()
	if colAfter <= colBefore {
		t.Fatal("expected MoveRight to advance cursor")
	}

	adapter.MoveWordForward()
	_, colWord := adapter.cursorPos()
	if colWord <= colAfter {
		t.Fatal("expected MoveWordForward to advance cursor")
	}
}

func TestTextAreaAdapterJump(t *testing.T) {
	area, _ := setupTextArea(t, "aa\nbb\ncc")
	adapter := NewTextAreaAdapter(area)

	adapter.JumpToEnd()
	rowEnd, _ := adapter.cursorPos()
	if rowEnd < 2 {
		t.Fatalf("expected JumpToEnd to reach bottom row, got %d", rowEnd)
	}

	adapter.JumpToStart()
	rowStart, _ := adapter.cursorPos()
	if rowStart != 0 {
		t.Fatalf("expected JumpToStart to reach top row, got %d", rowStart)
	}
}

func TestTextAreaAdapterInsertAndDelete(t *testing.T) {
	area, _ := setupTextArea(t, "")
	adapter := NewTextAreaAdapter(area)

	adapter.InsertRune('a')
	if area.GetText() != "a" {
		t.Fatalf("expected text to be 'a', got %q", area.GetText())
	}

	adapter.DeleteBackward()
	if area.GetText() != "" {
		t.Fatalf("expected text to be empty, got %q", area.GetText())
	}

	adapter.InsertNewline()
	if !strings.Contains(area.GetText(), "\n") {
		t.Fatal("expected InsertNewline to add newline")
	}
}
