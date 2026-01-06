package cmdline

import (
	"context"
	"errors"
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

type errorExecutor struct {
	err     error
	lastCtx context.Context
}

type ctxKey string

const commandLineCtxKey ctxKey = "ctx-key"

func (e *errorExecutor) Execute(ctx context.Context, _ string) error {
	e.lastCtx = ctx
	return e.err
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

func TestCommandLineCallbacksAndContext(t *testing.T) {
	app := tview.NewApplication()
	dummy := tview.NewBox()
	app.SetRoot(dummy, true)
	app.SetFocus(dummy)

	commandLine := NewCommandLine(app, NewCommandHistory(10))
	exec := &errorExecutor{err: errors.New("boom")}
	commandLine.SetExecutor(exec)

	commandLine.SetOnShow(nil)
	commandLine.SetOnHide(nil)
	commandLine.SetOnError(nil)
	commandLine.SetOnCancel(nil)

	showCount := 0
	hideCount := 0
	errorCount := 0
	commandLine.SetOnShow(func() {
		showCount++
	})
	commandLine.SetOnHide(func() {
		hideCount++
	})
	commandLine.SetOnError(func(err error) {
		if err != nil {
			errorCount++
		}
	})

	ctx := context.WithValue(context.Background(), commandLineCtxKey, "ctx-value")
	commandLine.SetContext(ctx)
	commandLine.Activate()
	commandLine.SetText("w test.sql")
	commandLine.handleDone(tcell.KeyEnter)

	if showCount != 1 || hideCount != 1 {
		t.Fatalf("expected show/hide to fire once, got show=%d hide=%d", showCount, hideCount)
	}
	if errorCount != 1 {
		t.Fatalf("expected onError to fire once, got %d", errorCount)
	}
	if exec.lastCtx == nil || exec.lastCtx.Value(commandLineCtxKey) != "ctx-value" {
		t.Fatal("expected context to be passed to executor")
	}

	exec.err = nil
	commandLine.SetContext(context.Background())
	commandLine.Activate()
	commandLine.SetText("w other.sql")
	commandLine.handleDone(tcell.KeyEnter)
	if exec.lastCtx == nil || exec.lastCtx.Value(commandLineCtxKey) != nil {
		t.Fatal("expected context to reset to background")
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

func TestCommandLineInputDefault(t *testing.T) {
	commandLine := NewCommandLine(tview.NewApplication(), NewCommandHistory(10))
	if event := commandLine.handleInput(nil); event != nil {
		t.Fatal("expected nil event to remain nil")
	}

	event := tcell.NewEventKey(tcell.KeyRune, 'x', tcell.ModNone)
	if got := commandLine.handleInput(event); got != event {
		t.Fatal("expected non-navigation key to pass through")
	}
}
