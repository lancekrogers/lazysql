package cmdline

import "testing"

func TestCommandHistoryNavigation(t *testing.T) {
	history := NewCommandHistory(3)
	history.Add("one")
	history.Add("two")
	history.Add("two")
	history.Add("three")

	if len(history.commands) != 3 {
		t.Fatalf("expected 3 commands, got %d", len(history.commands))
	}

	if cmd, ok := history.Previous(); !ok || cmd != "three" {
		t.Fatalf("expected previous to be three, got %q (ok=%v)", cmd, ok)
	}
	if cmd, ok := history.Previous(); !ok || cmd != "two" {
		t.Fatalf("expected previous to be two, got %q (ok=%v)", cmd, ok)
	}
	if cmd, ok := history.Next(); !ok || cmd != "three" {
		t.Fatalf("expected next to be three, got %q (ok=%v)", cmd, ok)
	}
}

func TestCommandHistoryCapacity(t *testing.T) {
	history := NewCommandHistory(2)
	history.Add("one")
	history.Add("two")
	history.Add("three")

	if len(history.commands) != 2 {
		t.Fatalf("expected 2 commands, got %d", len(history.commands))
	}
	if history.commands[0] != "two" || history.commands[1] != "three" {
		t.Fatalf("unexpected history contents: %v", history.commands)
	}
}
