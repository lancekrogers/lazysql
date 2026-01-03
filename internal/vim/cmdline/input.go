package cmdline

import (
	"context"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Executor interface {
	Execute(ctx context.Context, command string) error
}

type CommandLine struct {
	*tview.InputField
	app           *tview.Application
	history       *CommandHistory
	executor      Executor
	ctx           context.Context
	onCancel      func()
	onError       func(error)
	onShow        func()
	onHide        func()
	previousFocus tview.Primitive
}

func NewCommandLine(app *tview.Application, history *CommandHistory) *CommandLine {
	if history == nil {
		history = NewCommandHistory(100)
	}
	commandLine := &CommandLine{
		InputField: tview.NewInputField(),
		app:        app,
		history:    history,
		ctx:        context.Background(),
		onCancel:   func() {},
		onError:    func(error) {},
		onShow:     func() {},
		onHide:     func() {},
	}
	commandLine.SetLabel(":")
	commandLine.SetFieldWidth(0)
	commandLine.SetDoneFunc(commandLine.handleDone)
	commandLine.SetInputCapture(commandLine.handleInput)
	return commandLine
}

func (c *CommandLine) SetExecutor(executor Executor) {
	c.executor = executor
}

func (c *CommandLine) SetContext(ctx context.Context) {
	if ctx == nil {
		c.ctx = context.Background()
		return
	}
	c.ctx = ctx
}

func (c *CommandLine) SetOnCancel(onCancel func()) {
	if onCancel == nil {
		c.onCancel = func() {}
		return
	}
	c.onCancel = onCancel
}

func (c *CommandLine) SetOnError(onError func(error)) {
	if onError == nil {
		c.onError = func(error) {}
		return
	}
	c.onError = onError
}

func (c *CommandLine) SetOnShow(onShow func()) {
	if onShow == nil {
		c.onShow = func() {}
		return
	}
	c.onShow = onShow
}

func (c *CommandLine) SetOnHide(onHide func()) {
	if onHide == nil {
		c.onHide = func() {}
		return
	}
	c.onHide = onHide
}

func (c *CommandLine) Activate() {
	if c.app == nil {
		return
	}
	c.previousFocus = c.app.GetFocus()
	c.SetText("")
	c.history.ResetCursor()
	c.onShow()
	c.app.SetFocus(c)
}

func (c *CommandLine) Hide() {
	c.SetText("")
	c.onHide()
	if c.app != nil && c.previousFocus != nil {
		c.app.SetFocus(c.previousFocus)
	}
	c.previousFocus = nil
}

func (c *CommandLine) handleDone(key tcell.Key) {
	switch key {
	case tcell.KeyEnter:
		command := strings.TrimSpace(c.GetText())
		if command != "" {
			c.history.Add(command)
			if c.executor != nil {
				if err := c.executor.Execute(c.ctx, command); err != nil {
					c.onError(err)
				}
			}
		}
		c.Hide()
	case tcell.KeyEscape:
		c.onCancel()
		c.Hide()
	}
}

func (c *CommandLine) handleInput(event *tcell.EventKey) *tcell.EventKey {
	if event == nil {
		return event
	}
	switch event.Key() {
	case tcell.KeyUp:
		if cmd, ok := c.history.Previous(); ok {
			c.SetText(cmd)
		}
		return nil
	case tcell.KeyDown:
		if cmd, ok := c.history.Next(); ok {
			c.SetText(cmd)
		} else {
			c.SetText("")
		}
		return nil
	}
	return event
}
