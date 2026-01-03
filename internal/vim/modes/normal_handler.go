package modes

import (
	"time"

	"github.com/gdamore/tcell/v2"
)

type Editor interface {
	MoveLeft()
	MoveRight()
	MoveUp()
	MoveDown()
	MoveWordForward()
	MoveWordBackward()
	MoveLineStart()
	MoveLineEnd()
	JumpToStart()
	JumpToEnd()
	InsertNewline()
}

type LeaderActivator interface {
	Activate()
}

type CommandLineActivator interface {
	Activate()
}

type NormalHandler struct {
	modeManager     *ModeManager
	leader          LeaderActivator
	commandLine     CommandLineActivator
	editor          Editor
	lastKey         rune
	lastKeyTime     time.Time
	sequenceTimeout time.Duration
}

func NewNormalHandler(mm *ModeManager, leader LeaderActivator, commandLine CommandLineActivator, editor Editor) *NormalHandler {
	return &NormalHandler{
		modeManager:     mm,
		leader:          leader,
		commandLine:     commandLine,
		editor:          editor,
		sequenceTimeout: 500 * time.Millisecond,
	}
}

func (h *NormalHandler) SetSequenceTimeout(timeout time.Duration) {
	h.sequenceTimeout = timeout
}

func (h *NormalHandler) HandleKey(event *tcell.EventKey) bool {
	if event == nil {
		return false
	}

	if event.Key() != tcell.KeyRune {
		return false
	}

	key := event.Rune()
	if h.handleSequence(key) {
		return true
	}

	switch key {
	case 'h':
		h.editor.MoveLeft()
		return true
	case 'j':
		h.editor.MoveDown()
		return true
	case 'k':
		h.editor.MoveUp()
		return true
	case 'l':
		h.editor.MoveRight()
		return true
	case 'w':
		h.editor.MoveWordForward()
		return true
	case 'b':
		h.editor.MoveWordBackward()
		return true
	case 'G':
		h.editor.JumpToEnd()
		return true
	}

	switch key {
	case 'i':
		h.modeManager.EnterInsert()
		return true
	case 'I':
		h.editor.MoveLineStart()
		h.modeManager.EnterInsert()
		return true
	case 'a':
		h.editor.MoveRight()
		h.modeManager.EnterInsert()
		return true
	case 'A':
		h.editor.MoveLineEnd()
		h.modeManager.EnterInsert()
		return true
	case 'o':
		h.editor.MoveLineEnd()
		h.editor.InsertNewline()
		h.modeManager.EnterInsert()
		return true
	case 'O':
		h.editor.MoveLineStart()
		h.editor.InsertNewline()
		h.editor.MoveUp()
		h.modeManager.EnterInsert()
		return true
	}

	if key == '\\' {
		if h.leader != nil {
			h.leader.Activate()
		}
		return true
	}

	if key == ':' {
		if h.commandLine != nil {
			h.commandLine.Activate()
		}
		return true
	}

	return false
}

func (h *NormalHandler) handleSequence(key rune) bool {
	if h.lastKey != 0 && time.Since(h.lastKeyTime) > h.sequenceTimeout {
		h.lastKey = 0
	}

	if key == 'g' {
		if h.lastKey == 'g' {
			h.lastKey = 0
			h.editor.JumpToStart()
			return true
		}
		h.lastKey = key
		h.lastKeyTime = time.Now()
		return true
	}

	h.lastKey = 0
	return false
}
