package modes

import (
	"testing"
	"time"

	"github.com/gdamore/tcell/v2"
)

type mockEditor struct {
	left         int
	right        int
	up           int
	down         int
	wordForward  int
	wordBackward int
	lineStart    int
	lineEnd      int
	jumpStart    int
	jumpEnd      int
	newline      int
}

func (m *mockEditor) MoveLeft()         { m.left++ }
func (m *mockEditor) MoveRight()        { m.right++ }
func (m *mockEditor) MoveUp()           { m.up++ }
func (m *mockEditor) MoveDown()         { m.down++ }
func (m *mockEditor) MoveWordForward()  { m.wordForward++ }
func (m *mockEditor) MoveWordBackward() { m.wordBackward++ }
func (m *mockEditor) MoveLineStart()    { m.lineStart++ }
func (m *mockEditor) MoveLineEnd()      { m.lineEnd++ }
func (m *mockEditor) JumpToStart()      { m.jumpStart++ }
func (m *mockEditor) JumpToEnd()        { m.jumpEnd++ }
func (m *mockEditor) InsertNewline()    { m.newline++ }

type mockActivator struct {
	count int
}

func (m *mockActivator) Activate() { m.count++ }

func TestNormalHandlerNavigation(t *testing.T) {
	editor := &mockEditor{}
	handler := NewNormalHandler(NewManager(), nil, nil, editor)

	cases := []struct {
		key   rune
		field *int
	}{
		{'h', &editor.left},
		{'j', &editor.down},
		{'k', &editor.up},
		{'l', &editor.right},
		{'w', &editor.wordForward},
		{'b', &editor.wordBackward},
	}

	for _, testCase := range cases {
		event := tcell.NewEventKey(tcell.KeyRune, testCase.key, tcell.ModNone)
		if !handler.HandleKey(event) {
			t.Fatalf("expected key %q to be handled", testCase.key)
		}
		if *testCase.field == 0 {
			t.Fatalf("expected handler for key %q to increment counter", testCase.key)
		}
	}
}

func TestNormalHandlerSequences(t *testing.T) {
	editor := &mockEditor{}
	handler := NewNormalHandler(NewManager(), nil, nil, editor)
	handler.SetSequenceTimeout(time.Hour)

	if handler.HandleKey(tcell.NewEventKey(tcell.KeyRune, 'g', tcell.ModNone)) != true {
		t.Fatal("expected first g to be handled")
	}
	if handler.HandleKey(tcell.NewEventKey(tcell.KeyRune, 'g', tcell.ModNone)) != true {
		t.Fatal("expected second g to be handled")
	}
	if editor.jumpStart != 1 {
		t.Fatalf("expected JumpToStart once, got %d", editor.jumpStart)
	}

	if !handler.HandleKey(tcell.NewEventKey(tcell.KeyRune, 'G', tcell.ModNone)) {
		t.Fatal("expected G to be handled")
	}
	if editor.jumpEnd != 1 {
		t.Fatalf("expected JumpToEnd once, got %d", editor.jumpEnd)
	}
}

func TestNormalHandlerInsertTransitions(t *testing.T) {
	editor := &mockEditor{}
	manager := NewManager()
	handler := NewNormalHandler(manager, nil, nil, editor)

	handler.HandleKey(tcell.NewEventKey(tcell.KeyRune, 'i', tcell.ModNone))
	if manager.Mode() != ModeInsert {
		t.Fatal("expected mode to switch to insert on i")
	}

	manager.EnterNormal()
	handler.HandleKey(tcell.NewEventKey(tcell.KeyRune, 'I', tcell.ModNone))
	if editor.lineStart == 0 {
		t.Fatal("expected I to move to line start")
	}

	manager.EnterNormal()
	handler.HandleKey(tcell.NewEventKey(tcell.KeyRune, 'a', tcell.ModNone))
	if editor.right == 0 {
		t.Fatal("expected a to move right before insert")
	}

	manager.EnterNormal()
	handler.HandleKey(tcell.NewEventKey(tcell.KeyRune, 'A', tcell.ModNone))
	if editor.lineEnd == 0 {
		t.Fatal("expected A to move to line end")
	}

	manager.EnterNormal()
	handler.HandleKey(tcell.NewEventKey(tcell.KeyRune, 'o', tcell.ModNone))
	if editor.newline == 0 {
		t.Fatal("expected o to insert newline")
	}

	manager.EnterNormal()
	handler.HandleKey(tcell.NewEventKey(tcell.KeyRune, 'O', tcell.ModNone))
	if editor.newline < 2 {
		t.Fatal("expected O to insert newline")
	}
}

func TestNormalHandlerLeaderAndCommandLine(t *testing.T) {
	editor := &mockEditor{}
	leader := &mockActivator{}
	commandLine := &mockActivator{}
	handler := NewNormalHandler(NewManager(), leader, commandLine, editor)

	handler.HandleKey(tcell.NewEventKey(tcell.KeyRune, '\\', tcell.ModNone))
	if leader.count != 1 {
		t.Fatal("expected leader activation")
	}

	handler.HandleKey(tcell.NewEventKey(tcell.KeyRune, ':', tcell.ModNone))
	if commandLine.count != 1 {
		t.Fatal("expected command line activation")
	}
}

func TestNormalHandlerIgnoresUnknownKey(t *testing.T) {
	handler := NewNormalHandler(NewManager(), nil, nil, &mockEditor{})
	if handler.HandleKey(tcell.NewEventKey(tcell.KeyRune, 'x', tcell.ModNone)) {
		t.Fatal("expected unbound key to be unhandled")
	}
}
