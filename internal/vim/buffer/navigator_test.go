package buffer

import "testing"

func TestNavigationNextPrevious(t *testing.T) {
	manager := NewManager()
	first := manager.Create("one")
	second := manager.Create("two")
	third := manager.Create("three")

	navigator := NewNavigator(manager)
	if err := navigator.GoTo(second.ID); err != nil {
		t.Fatalf("goto: %v", err)
	}

	if err := navigator.Next(); err != nil {
		t.Fatalf("next: %v", err)
	}
	if manager.ActiveID() != third.ID {
		t.Fatalf("expected third buffer, got %d", manager.ActiveID())
	}

	if err := navigator.Next(); err != nil {
		t.Fatalf("next wrap: %v", err)
	}
	if manager.ActiveID() != first.ID {
		t.Fatalf("expected first buffer after wrap, got %d", manager.ActiveID())
	}

	if err := navigator.Previous(); err != nil {
		t.Fatalf("previous wrap: %v", err)
	}
	if manager.ActiveID() != third.ID {
		t.Fatalf("expected third buffer after previous wrap, got %d", manager.ActiveID())
	}
}

func TestNavigationAlternate(t *testing.T) {
	manager := NewManager()
	first := manager.Create("one")
	second := manager.Create("two")
	navigator := NewNavigator(manager)

	if err := navigator.GoTo(first.ID); err != nil {
		t.Fatalf("goto: %v", err)
	}
	if err := navigator.GoTo(second.ID); err != nil {
		t.Fatalf("goto: %v", err)
	}

	alternate, ok := navigator.Alternate()
	if !ok {
		t.Fatal("expected alternate buffer")
	}
	if alternate != first.ID {
		t.Fatalf("expected alternate to be first, got %d", alternate)
	}
}

func TestNavigationIndex(t *testing.T) {
	manager := NewManager()
	first := manager.Create("one")
	second := manager.Create("two")
	navigator := NewNavigator(manager)

	if err := navigator.GoToIndex(2); err != nil {
		t.Fatalf("goto index: %v", err)
	}
	if manager.ActiveID() != second.ID {
		t.Fatalf("expected second buffer, got %d", manager.ActiveID())
	}

	if err := navigator.GoToIndex(1); err != nil {
		t.Fatalf("goto index: %v", err)
	}
	if manager.ActiveID() != first.ID {
		t.Fatalf("expected first buffer, got %d", manager.ActiveID())
	}
}

func TestHistoryCapacity(t *testing.T) {
	history := NewNavigationHistory(2)
	history.Push(1)
	history.Push(2)
	history.Push(3)

	alternate, ok := history.Alternate()
	if !ok {
		t.Fatal("expected alternate")
	}
	if alternate != 2 {
		t.Fatalf("expected alternate to be 2, got %d", alternate)
	}
}
