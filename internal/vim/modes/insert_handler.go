package modes

import "github.com/gdamore/tcell/v2"

type TextEditor interface {
	InsertRune(r rune)
	InsertNewline()
	DeleteBackward()
	DeleteForward()
}

type InsertHandler struct {
	modeManager *ModeManager
	editor      TextEditor
}

func NewInsertHandler(mm *ModeManager, editor TextEditor) *InsertHandler {
	return &InsertHandler{
		modeManager: mm,
		editor:      editor,
	}
}

func (h *InsertHandler) HandleKey(event *tcell.EventKey) bool {
	if event == nil {
		return false
	}

	switch event.Key() {
	case tcell.KeyEscape, tcell.KeyCtrlC:
		h.modeManager.EnterNormal()
		return true
	case tcell.KeyEnter:
		h.editor.InsertNewline()
		return true
	case tcell.KeyBackspace, tcell.KeyBackspace2:
		h.editor.DeleteBackward()
		return true
	case tcell.KeyDelete:
		h.editor.DeleteForward()
		return true
	case tcell.KeyTab:
		h.editor.InsertRune('\t')
		return true
	case tcell.KeyRune:
		if r := event.Rune(); r != 0 {
			h.editor.InsertRune(r)
			return true
		}
	}

	return false
}
