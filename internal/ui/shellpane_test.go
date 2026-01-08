package ui

import (
	"strings"
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

func TestShellPanePreferredHeight(t *testing.T) {
	pane := NewShellPane()
	if got := pane.PreferredHeight(100); got != 30 {
		t.Fatalf("expected preferred height 30, got %d", got)
	}
	if got := pane.PreferredHeight(5); got != 3 {
		t.Fatalf("expected preferred height 3 for small height, got %d", got)
	}
	if got := pane.PreferredHeight(0); got != 10 {
		t.Fatalf("expected fallback height 10, got %d", got)
	}

	pane.SetHeight(12)
	if got := pane.PreferredHeight(100); got != 12 {
		t.Fatalf("expected preferred height to respect explicit height, got %d", got)
	}
}

func TestShellPaneToggle(t *testing.T) {
	pane := NewShellPane()
	if pane.IsVisible() {
		t.Fatal("expected shell pane hidden by default")
	}

	pane.Toggle()
	if !pane.IsVisible() {
		t.Fatal("expected shell pane visible after toggle")
	}

	pane.Toggle()
	if pane.IsVisible() {
		t.Fatal("expected shell pane hidden after second toggle")
	}
}

func TestShellPaneOutputHelpers(t *testing.T) {
	pane := NewShellPane()
	if pane.Input() == nil {
		t.Fatal("expected input field")
	}
	if pane.Output() == nil {
		t.Fatal("expected output view")
	}

	pane.AppendOutput("hello world")
	if text := pane.Output().GetText(false); !strings.Contains(text, "hello world") {
		t.Fatalf("expected output to contain message, got %q", text)
	}

	pane.ClearOutput()
	if text := pane.Output().GetText(false); text != "" {
		t.Fatalf("expected output cleared, got %q", text)
	}
}

func TestShellPaneHistorySnapshot(t *testing.T) {
	pane := NewShellPane()
	pane.SetInputText("select 1")
	pane.handleInputDone(tcell.KeyEnter)

	history := pane.History()
	if len(history) != 1 || history[0] != "select 1" {
		t.Fatalf("expected history entry, got %#v", history)
	}

	history[0] = "mutated"
	updated := pane.History()
	if updated[0] == "mutated" {
		t.Fatal("expected history to be a copy")
	}
}

func TestShellPaneInputDoneNoop(t *testing.T) {
	pane := NewShellPane()
	pane.SetInputText("select 1")
	pane.handleInputDone(tcell.KeyEscape)
	if len(pane.History()) != 0 {
		t.Fatalf("expected no history for non-enter key, got %#v", pane.History())
	}

	pane.SetInputText("")
	pane.handleInputDone(tcell.KeyEnter)
	if len(pane.History()) != 0 {
		t.Fatalf("expected no history for empty query, got %#v", pane.History())
	}
}

func TestShellPaneInputCaptureNoop(t *testing.T) {
	pane := NewShellPane()
	if pane.handleInputCapture(nil) != nil {
		t.Fatal("expected nil when input event is nil")
	}

	event := tcell.NewEventKey(tcell.KeyRune, 'x', tcell.ModNone)
	if pane.handleInputCapture(event) != event {
		t.Fatal("expected non-navigation keys to pass through")
	}
}

func TestShellPaneNavigateHistoryBounds(t *testing.T) {
	pane := NewShellPane()
	pane.SetInputText("select 1")
	pane.handleInputDone(tcell.KeyEnter)

	pane.navigateHistory(-2)
	if pane.Input().GetText() != "" {
		t.Fatalf("expected out-of-bounds history to keep input empty, got %q", pane.Input().GetText())
	}

	pane.navigateHistory(2)
	if pane.Input().GetText() != "" {
		t.Fatalf("expected out-of-bounds history to keep input empty, got %q", pane.Input().GetText())
	}
}

func TestShellPaneNilGuards(t *testing.T) {
	t.Helper()
	pane := NewShellPane()
	pane.input = nil
	pane.output = nil

	pane.SetInputText("noop")
	pane.AppendOutput("noop")
	pane.ClearOutput()
	pane.handleInputDone(tcell.KeyEnter)
	pane.navigateHistory(1)
}
