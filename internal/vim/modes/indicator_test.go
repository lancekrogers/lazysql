package modes

import "testing"

func TestModeIndicatorUpdates(t *testing.T) {
	manager := NewManager()
	indicator := NewModeIndicator(manager)

	if got := indicator.GetText(true); got != "-- NORMAL --" {
		t.Fatalf("expected NORMAL indicator, got %q", got)
	}

	manager.EnterInsert()
	if got := indicator.GetText(true); got != "-- INSERT --" {
		t.Fatalf("expected INSERT indicator, got %q", got)
	}
}

func TestModeIndicatorUnknown(t *testing.T) {
	manager := NewManager()
	indicator := NewModeIndicator(manager)
	indicator.update(VimMode(99))

	if got := indicator.GetText(true); got != "-- UNKNOWN --" {
		t.Fatalf("expected UNKNOWN indicator, got %q", got)
	}
}
