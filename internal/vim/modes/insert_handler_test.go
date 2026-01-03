package modes

import (
	"testing"

	"github.com/gdamore/tcell/v2"
)

type mockTextEditor struct {
	runes     []rune
	newlines  int
	backspace int
	delete    int
}

func (m *mockTextEditor) InsertRune(r rune) { m.runes = append(m.runes, r) }
func (m *mockTextEditor) InsertNewline()    { m.newlines++ }
func (m *mockTextEditor) DeleteBackward()   { m.backspace++ }
func (m *mockTextEditor) DeleteForward()    { m.delete++ }

func TestInsertHandlerTextInput(t *testing.T) {
	editor := &mockTextEditor{}
	manager := NewManager()
	manager.EnterInsert()
	handler := NewInsertHandler(manager, editor)

	if !handler.HandleKey(tcell.NewEventKey(tcell.KeyRune, 'a', tcell.ModNone)) {
		t.Fatal("expected rune input to be handled")
	}
	if len(editor.runes) != 1 || editor.runes[0] != 'a' {
		t.Fatal("expected rune to be inserted")
	}

	handler.HandleKey(tcell.NewEventKey(tcell.KeyTab, 0, tcell.ModNone))
	if len(editor.runes) != 2 || editor.runes[1] != '\t' {
		t.Fatal("expected tab to be inserted")
	}
}

func TestInsertHandlerEditingKeys(t *testing.T) {
	editor := &mockTextEditor{}
	manager := NewManager()
	manager.EnterInsert()
	handler := NewInsertHandler(manager, editor)

	handler.HandleKey(tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone))
	handler.HandleKey(tcell.NewEventKey(tcell.KeyBackspace, 0, tcell.ModNone))
	handler.HandleKey(tcell.NewEventKey(tcell.KeyDelete, 0, tcell.ModNone))

	if editor.newlines != 1 {
		t.Fatalf("expected newline count 1, got %d", editor.newlines)
	}
	if editor.backspace != 1 {
		t.Fatalf("expected backspace count 1, got %d", editor.backspace)
	}
	if editor.delete != 1 {
		t.Fatalf("expected delete count 1, got %d", editor.delete)
	}
}

func TestInsertHandlerEscape(t *testing.T) {
	editor := &mockTextEditor{}
	manager := NewManager()
	manager.EnterInsert()
	handler := NewInsertHandler(manager, editor)

	handler.HandleKey(tcell.NewEventKey(tcell.KeyEscape, 0, tcell.ModNone))
	if manager.Mode() != ModeNormal {
		t.Fatal("expected escape to return to normal mode")
	}

	manager.EnterInsert()
	handler.HandleKey(tcell.NewEventKey(tcell.KeyCtrlC, 0, tcell.ModNone))
	if manager.Mode() != ModeNormal {
		t.Fatal("expected ctrl+c to return to normal mode")
	}
}

func TestInsertHandlerIgnoresUnknownKey(t *testing.T) {
	editor := &mockTextEditor{}
	handler := NewInsertHandler(NewManager(), editor)
	if handler.HandleKey(tcell.NewEventKey(tcell.KeyUp, 0, tcell.ModNone)) {
		t.Fatal("expected unknown key to be unhandled")
	}
}
