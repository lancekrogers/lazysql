package ui

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type ShellPane struct {
	*tview.Flex
	input      *tview.InputField
	output     *tview.TextView
	history    []string
	historyIdx int
	visible    bool
	height     int
	onExecute  func(string)
}

func NewShellPane() *ShellPane {
	pane := &ShellPane{
		Flex:       tview.NewFlex().SetDirection(tview.FlexRow),
		input:      tview.NewInputField(),
		output:     tview.NewTextView(),
		history:    make([]string, 0, 100),
		visible:    false,
		height:     10,
		historyIdx: 0,
	}

	pane.output.SetDynamicColors(true)
	pane.output.SetScrollable(true)
	pane.output.SetWordWrap(true)

	pane.input.SetLabel("SQL> ")
	pane.input.SetDoneFunc(pane.handleInputDone)
	pane.input.SetInputCapture(pane.handleInputCapture)

	pane.AddItem(pane.output, 0, 1, false)
	pane.AddItem(pane.input, 1, 0, true)

	return pane
}

func (p *ShellPane) SetOnExecute(fn func(string)) {
	p.onExecute = fn
}

func (p *ShellPane) Input() *tview.InputField {
	return p.input
}

func (p *ShellPane) Output() *tview.TextView {
	return p.output
}

func (p *ShellPane) Height() int {
	return p.height
}

func (p *ShellPane) SetHeight(height int) {
	if height <= 0 {
		return
	}
	p.height = height
}

func (p *ShellPane) IsVisible() bool {
	return p.visible
}

func (p *ShellPane) Show() {
	p.visible = true
	p.SetBorder(true)
	p.SetTitle(" SQL Shell ")
}

func (p *ShellPane) Hide() {
	p.visible = false
	p.SetBorder(false)
	p.SetTitle("")
}

func (p *ShellPane) Toggle() {
	if p.visible {
		p.Hide()
	} else {
		p.Show()
	}
}

func (p *ShellPane) AppendOutput(message string) {
	if p.output == nil {
		return
	}
	if _, err := fmt.Fprintln(p.output, message); err != nil {
		return
	}
}

func (p *ShellPane) ClearOutput() {
	if p.output == nil {
		return
	}
	p.output.Clear()
}

func (p *ShellPane) History() []string {
	history := make([]string, len(p.history))
	copy(history, p.history)
	return history
}

func (p *ShellPane) SetInputText(text string) {
	if p.input == nil {
		return
	}
	p.input.SetText(text)
}

func (p *ShellPane) handleInputDone(key tcell.Key) {
	if key != tcell.KeyEnter {
		return
	}
	if p.input == nil {
		return
	}
	query := p.input.GetText()
	if query == "" {
		return
	}
	p.addToHistory(query)
	if p.onExecute != nil {
		p.onExecute(query)
	}
	p.input.SetText("")
}

func (p *ShellPane) handleInputCapture(event *tcell.EventKey) *tcell.EventKey {
	if event == nil {
		return nil
	}
	switch event.Key() {
	case tcell.KeyUp:
		p.navigateHistory(-1)
		return nil
	case tcell.KeyDown:
		p.navigateHistory(1)
		return nil
	}
	return event
}

func (p *ShellPane) addToHistory(query string) {
	p.history = append(p.history, query)
	p.historyIdx = len(p.history)
}

func (p *ShellPane) navigateHistory(delta int) {
	if p.input == nil {
		return
	}
	newIdx := p.historyIdx + delta
	if newIdx < 0 || newIdx > len(p.history) {
		return
	}
	p.historyIdx = newIdx
	if p.historyIdx < len(p.history) {
		p.input.SetText(p.history[p.historyIdx])
	} else {
		p.input.SetText("")
	}
}
