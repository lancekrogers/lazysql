package buffer

import (
	"errors"
	"testing"
)

func TestBufferCreateAndActivate(t *testing.T) {
	manager := NewManager()
	buffer := manager.Create("")
	if buffer.Name != "Untitled-1" {
		t.Fatalf("expected default name Untitled-1, got %q", buffer.Name)
	}
	if manager.ActiveID() != buffer.ID {
		t.Fatal("expected created buffer to be active")
	}

	other := manager.Create("query.sql")
	if other.Name != "query.sql" {
		t.Fatalf("expected name query.sql, got %q", other.Name)
	}
	if manager.ActiveID() != other.ID {
		t.Fatal("expected latest buffer to be active")
	}
}

func TestBufferUpdateContent(t *testing.T) {
	manager := NewManager()
	buffer := manager.Create("test")
	if buffer.Dirty {
		t.Fatal("expected new buffer to be clean")
	}
	if err := manager.UpdateContent(buffer.ID, "select 1"); err != nil {
		t.Fatalf("update content: %v", err)
	}
	if !buffer.Dirty {
		t.Fatal("expected buffer to be dirty after update")
	}
	if buffer.Content != "select 1" {
		t.Fatalf("unexpected content: %q", buffer.Content)
	}
}

func TestBufferMarkSaved(t *testing.T) {
	manager := NewManager()
	buffer := manager.Create("test")
	if err := manager.UpdateContent(buffer.ID, "select 1"); err != nil {
		t.Fatalf("update content: %v", err)
	}

	if err := manager.MarkSaved(buffer.ID); err != nil {
		t.Fatalf("mark saved: %v", err)
	}
	if buffer.Dirty {
		t.Fatal("expected buffer to be clean after save")
	}
}

func TestBufferClose(t *testing.T) {
	manager := NewManager()
	first := manager.Create("one")
	second := manager.Create("two")
	third := manager.Create("three")

	if err := manager.SetActive(second.ID); err != nil {
		t.Fatalf("set active: %v", err)
	}
	if err := manager.Close(second.ID); err != nil {
		t.Fatalf("close buffer: %v", err)
	}
	if manager.ActiveID() != third.ID {
		t.Fatalf("expected next buffer to be active, got %d", manager.ActiveID())
	}
	if err := manager.Close(first.ID); err != nil {
		t.Fatalf("close buffer: %v", err)
	}
	if err := manager.Close(third.ID); err != nil {
		t.Fatalf("close buffer: %v", err)
	}
	if manager.ActiveID() != 0 {
		t.Fatalf("expected no active buffer, got %d", manager.ActiveID())
	}
}

func TestBufferCloseSelectsNeighbor(t *testing.T) {
	manager := NewManager()
	first := manager.Create("one")
	second := manager.Create("two")
	third := manager.Create("three")

	if err := manager.SetActive(first.ID); err != nil {
		t.Fatalf("set active: %v", err)
	}
	if err := manager.Close(first.ID); err != nil {
		t.Fatalf("close buffer: %v", err)
	}
	if manager.ActiveID() != second.ID {
		t.Fatalf("expected next buffer to be active, got %d", manager.ActiveID())
	}
	if err := manager.Close(second.ID); err != nil {
		t.Fatalf("close buffer: %v", err)
	}
	if manager.ActiveID() != third.ID {
		t.Fatalf("expected remaining buffer to be active, got %d", manager.ActiveID())
	}
}

func TestBufferErrors(t *testing.T) {
	manager := NewManager()
	if err := manager.SetActive(99); err == nil {
		t.Fatal("expected error for missing buffer")
	}
	if err := manager.UpdateContent(99, "nope"); err == nil {
		t.Fatal("expected error for missing buffer")
	}
	if err := manager.MarkSaved(99); err == nil {
		t.Fatal("expected error for missing buffer")
	}
	if err := manager.Close(99); err == nil {
		t.Fatal("expected error for missing buffer")
	}
}

func TestBufferCountAndFindByPath(t *testing.T) {
	manager := NewManager()
	if manager.Count() != 0 {
		t.Fatalf("expected count 0, got %d", manager.Count())
	}
	first := manager.Create("one")
	second := manager.Create("two")

	if err := manager.UpdateMetadata(first.ID, "", "/tmp/one.sql"); err != nil {
		t.Fatalf("update metadata: %v", err)
	}
	if err := manager.UpdateMetadata(second.ID, "", "/tmp/two.sql"); err != nil {
		t.Fatalf("update metadata: %v", err)
	}

	if manager.Count() != 2 {
		t.Fatalf("expected count 2, got %d", manager.Count())
	}
	found := manager.FindByPath("/tmp/two.sql")
	if found == nil || found.ID != second.ID {
		t.Fatalf("expected to find second buffer, got %+v", found)
	}
}

func TestBufferUpdateMetadataErrors(t *testing.T) {
	manager := NewManager()
	if err := manager.UpdateMetadata(42, "missing", "/tmp/missing.sql"); err == nil {
		t.Fatal("expected error for missing buffer")
	}
}

func TestBufferCloseDirtyRequiresForce(t *testing.T) {
	manager := NewManager()
	buffer := manager.Create("one")

	if err := manager.UpdateContent(buffer.ID, "select 1"); err != nil {
		t.Fatalf("update content: %v", err)
	}
	if err := manager.Close(buffer.ID); !errors.Is(err, ErrBufferDirty) {
		t.Fatalf("expected ErrBufferDirty, got %v", err)
	}
	if err := manager.CloseForce(buffer.ID); err != nil {
		t.Fatalf("close force: %v", err)
	}
}

func TestBufferDirtyList(t *testing.T) {
	manager := NewManager()
	first := manager.Create("one")
	second := manager.Create("two")

	if err := manager.UpdateContent(second.ID, "select 2"); err != nil {
		t.Fatalf("update content: %v", err)
	}
	dirty := manager.GetDirtyBuffers()
	if len(dirty) != 1 {
		t.Fatalf("expected 1 dirty buffer, got %d", len(dirty))
	}
	if dirty[0].ID != second.ID {
		t.Fatalf("expected second buffer dirty, got %d", dirty[0].ID)
	}
	if first.Dirty {
		t.Fatal("expected first buffer to remain clean")
	}
}

func TestBufferListeners(t *testing.T) {
	manager := NewManager()
	events := make([]EventType, 0)
	manager.AddListener(func(event BufferEvent) {
		events = append(events, event.Type)
	})

	buffer := manager.Create("one")
	if err := manager.UpdateContent(buffer.ID, "select 1"); err != nil {
		t.Fatalf("update content: %v", err)
	}
	if err := manager.MarkSaved(buffer.ID); err != nil {
		t.Fatalf("mark saved: %v", err)
	}
	if err := manager.Close(buffer.ID); err != nil {
		t.Fatalf("close buffer: %v", err)
	}

	if len(events) < 4 {
		t.Fatalf("expected at least 4 events, got %d", len(events))
	}
	if events[0] != EventCreated || events[1] != EventActivated {
		t.Fatalf("unexpected event order: %v", events[:2])
	}
}
