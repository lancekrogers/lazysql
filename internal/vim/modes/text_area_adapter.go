package modes

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// TextAreaAdapter adapts tview.TextArea to Vim editor interfaces.
type TextAreaAdapter struct {
	area *tview.TextArea
}

func NewTextAreaAdapter(area *tview.TextArea) *TextAreaAdapter {
	return &TextAreaAdapter{area: area}
}

func (a *TextAreaAdapter) MoveLeft() {
	a.sendKey(tcell.KeyLeft, 0)
}

func (a *TextAreaAdapter) MoveRight() {
	a.sendKey(tcell.KeyRight, 0)
}

func (a *TextAreaAdapter) MoveUp() {
	a.sendKey(tcell.KeyUp, 0)
}

func (a *TextAreaAdapter) MoveDown() {
	a.sendKey(tcell.KeyDown, 0)
}

func (a *TextAreaAdapter) MoveWordForward() {
	a.sendKey(tcell.KeyRight, tcell.ModCtrl)
}

func (a *TextAreaAdapter) MoveWordBackward() {
	a.sendKey(tcell.KeyLeft, tcell.ModCtrl)
}

func (a *TextAreaAdapter) MoveLineStart() {
	a.sendKey(tcell.KeyHome, 0)
}

func (a *TextAreaAdapter) MoveLineEnd() {
	a.sendKey(tcell.KeyEnd, 0)
}

func (a *TextAreaAdapter) JumpToStart() {
	if a.area == nil {
		return
	}
	for {
		rowBefore, _ := a.cursorPos()
		a.sendKey(tcell.KeyUp, 0)
		rowAfter, _ := a.cursorPos()
		if rowAfter == rowBefore {
			break
		}
	}
	a.MoveLineStart()
}

func (a *TextAreaAdapter) JumpToEnd() {
	if a.area == nil {
		return
	}
	for {
		rowBefore, _ := a.cursorPos()
		a.sendKey(tcell.KeyDown, 0)
		rowAfter, _ := a.cursorPos()
		if rowAfter == rowBefore {
			break
		}
	}
	a.MoveLineEnd()
}

func (a *TextAreaAdapter) InsertNewline() {
	a.sendKey(tcell.KeyEnter, 0)
}

func (a *TextAreaAdapter) InsertRune(r rune) {
	if a.area == nil {
		return
	}
	handler := a.area.InputHandler()
	handler(tcell.NewEventKey(tcell.KeyRune, r, tcell.ModNone), nil)
}

func (a *TextAreaAdapter) DeleteBackward() {
	a.sendKey(tcell.KeyBackspace, 0)
}

func (a *TextAreaAdapter) DeleteForward() {
	a.sendKey(tcell.KeyDelete, 0)
}

func (a *TextAreaAdapter) sendKey(key tcell.Key, mod tcell.ModMask) {
	if a.area == nil {
		return
	}
	handler := a.area.InputHandler()
	handler(tcell.NewEventKey(key, 0, mod), nil)
}

func (a *TextAreaAdapter) cursorPos() (int, int) {
	_, _, row, col := a.area.GetCursor()
	return row, col
}
