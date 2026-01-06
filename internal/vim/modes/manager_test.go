package modes

import "testing"

func TestInitialMode(t *testing.T) {
	manager := NewManager()
	if manager.Mode() != ModeNormal {
		t.Fatalf("expected initial mode NORMAL, got %v", manager.Mode())
	}
}

func TestEnterInsert(t *testing.T) {
	manager := NewManager()
	if !manager.EnterInsert() {
		t.Fatal("expected transition to insert mode")
	}
	if manager.Mode() != ModeInsert {
		t.Fatalf("expected mode INSERT, got %v", manager.Mode())
	}
	if manager.EnterInsert() {
		t.Fatal("expected duplicate insert transition to be blocked")
	}
}

func TestEnterNormal(t *testing.T) {
	manager := NewManager()
	if manager.EnterNormal() {
		t.Fatal("expected normal transition to be ignored when already normal")
	}
	manager.EnterInsert()
	if !manager.EnterNormal() {
		t.Fatal("expected transition back to normal")
	}
	if manager.Mode() != ModeNormal {
		t.Fatalf("expected mode NORMAL, got %v", manager.Mode())
	}
}

func TestListenerNotifications(t *testing.T) {
	manager := NewManager()
	notifications := make([]VimMode, 0, 2)
	manager.AddListener(func(mode VimMode) {
		notifications = append(notifications, mode)
	})

	manager.EnterInsert()
	manager.EnterNormal()

	if len(notifications) != 2 {
		t.Fatalf("expected 2 notifications, got %d", len(notifications))
	}
	if notifications[0] != ModeInsert || notifications[1] != ModeNormal {
		t.Fatalf("unexpected notification order: %v", notifications)
	}
}

func TestAddNilListener(_ *testing.T) {
	manager := NewManager()
	manager.AddListener(nil)
	manager.EnterInsert()
}

func TestModeString(t *testing.T) {
	if ModeNormal.String() != "NORMAL" {
		t.Fatalf("unexpected ModeNormal string: %s", ModeNormal.String())
	}
	if ModeInsert.String() != "INSERT" {
		t.Fatalf("unexpected ModeInsert string: %s", ModeInsert.String())
	}
	if VimMode(99).String() != "UNKNOWN" {
		t.Fatalf("unexpected unknown mode string: %s", VimMode(99).String())
	}
}
