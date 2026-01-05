package ui

import (
	"testing"

	"github.com/gdamore/tcell/v2"
)

func TestShellPaneVisibility(t *testing.T) {
	pane := NewShellPane()
	if pane.IsVisible() {
		t.Fatal("expected shell pane hidden by default")
	}

	pane.Show()
	if !pane.IsVisible() {
		t.Fatal("expected shell pane visible after show")
	}

	pane.Hide()
	if pane.IsVisible() {
		t.Fatal("expected shell pane hidden after hide")
	}
}

func TestShellPaneHistoryNavigation(t *testing.T) {
	pane := NewShellPane()
	pane.SetInputText("select 1")
	pane.handleInputDone(tcell.KeyEnter)
	pane.SetInputText("select 2")
	pane.handleInputDone(tcell.KeyEnter)

	pane.handleInputCapture(tcell.NewEventKey(tcell.KeyUp, 0, tcell.ModNone))
	if pane.input.GetText() != "select 2" {
		t.Fatalf("expected latest history entry, got %q", pane.input.GetText())
	}

	pane.handleInputCapture(tcell.NewEventKey(tcell.KeyUp, 0, tcell.ModNone))
	if pane.input.GetText() != "select 1" {
		t.Fatalf("expected earlier history entry, got %q", pane.input.GetText())
	}

	pane.handleInputCapture(tcell.NewEventKey(tcell.KeyDown, 0, tcell.ModNone))
	if pane.input.GetText() != "select 2" {
		t.Fatalf("expected next history entry, got %q", pane.input.GetText())
	}
}

func TestShellPaneExecuteCallback(t *testing.T) {
	pane := NewShellPane()
	var executed string
	pane.SetOnExecute(func(query string) {
		executed = query
	})

	pane.SetInputText("select 1")
	pane.handleInputDone(tcell.KeyEnter)

	if executed != "select 1" {
		t.Fatalf("expected execute callback, got %q", executed)
	}
}

func TestShellPaneHeight(t *testing.T) {
	pane := NewShellPane()
	pane.SetHeight(12)
	if pane.Height() != 12 {
		t.Fatalf("expected height 12, got %d", pane.Height())
	}

	pane.SetHeight(0)
	if pane.Height() != 12 {
		t.Fatalf("expected height to remain 12, got %d", pane.Height())
	}
}
