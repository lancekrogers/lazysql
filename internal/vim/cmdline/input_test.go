package cmdline

import (
	"context"
	"testing"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type captureExecutor struct {
	commands []string
	ctx      context.Context
}

func (e *captureExecutor) Execute(ctx context.Context, command string) error {
	e.ctx = ctx
	e.commands = append(e.commands, command)
	return nil
}

func TestCommandLineExecuteAndFocus(t *testing.T) {
	app := tview.NewApplication()
	dummy := tview.NewBox()
	app.SetRoot(dummy, true)
	app.SetFocus(dummy)

	history := NewCommandHistory(10)
	commandLine := NewCommandLine(app, history)
	exec := &captureExecutor{}
	commandLine.SetExecutor(exec)

	commandLine.Activate()
	if app.GetFocus() != commandLine {
		t.Fatal("expected command line to be focused after activate")
	}

	commandLine.SetText("  w foo.sql  ")
	commandLine.handleDone(tcell.KeyEnter)

	if app.GetFocus() != dummy {
		t.Fatal("expected focus to return to previous primitive after execute")
	}
	if len(exec.commands) != 1 || exec.commands[0] != "w foo.sql" {
		t.Fatalf("unexpected commands: %v", exec.commands)
	}
	if len(history.commands) != 1 || history.commands[0] != "w foo.sql" {
		t.Fatalf("unexpected history: %v", history.commands)
	}
}

func TestCommandLineHistoryKeys(t *testing.T) {
	app := tview.NewApplication()
	commandLine := NewCommandLine(app, NewCommandHistory(10))
	commandLine.history.Add("one")
	commandLine.history.Add("two")

	commandLine.Activate()
	commandLine.handleInput(tcell.NewEventKey(tcell.KeyUp, 0, tcell.ModNone))
	if commandLine.GetText() != "two" {
		t.Fatalf("expected latest command, got %q", commandLine.GetText())
	}

	commandLine.handleInput(tcell.NewEventKey(tcell.KeyUp, 0, tcell.ModNone))
	if commandLine.GetText() != "one" {
		t.Fatalf("expected previous command, got %q", commandLine.GetText())
	}

	commandLine.handleInput(tcell.NewEventKey(tcell.KeyDown, 0, tcell.ModNone))
	if commandLine.GetText() != "two" {
		t.Fatalf("expected next command, got %q", commandLine.GetText())
	}
}

func TestCommandLineCancel(t *testing.T) {
	app := tview.NewApplication()
	commandLine := NewCommandLine(app, NewCommandHistory(10))
	canceled := false
	commandLine.SetOnCancel(func() {
		canceled = true
	})

	commandLine.handleDone(tcell.KeyEscape)
	if !canceled {
		t.Fatal("expected cancel callback to fire")
	}
}
