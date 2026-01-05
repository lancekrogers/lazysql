package commands

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/lancekrogers/lazysql/internal/vim/buffer"
	"github.com/lancekrogers/lazysql/internal/vim/cmdline"
)

func TestEditCommandExistingFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "query.sql")
	if err := os.WriteFile(path, []byte("select 1;\nselect 2;"), 0o644); err != nil {
		t.Fatalf("write file: %v", err)
	}

	manager := buffer.NewManager()
	status := &fakeStatus{}
	command := &EditCommand{}
	if command.Name() != "e" {
		t.Fatalf("expected command name e, got %q", command.Name())
	}
	err := command.Execute(context.Background(), cmdline.Invocation{Name: "e", Args: []string{path}}, &cmdline.CommandContext{
		Buffers: manager,
		Status:  status,
	})
	if err != nil {
		t.Fatalf("edit command: %v", err)
	}

	active := manager.Active()
	if active == nil {
		t.Fatal("expected active buffer")
	}
	if active.FilePath != path {
		t.Fatalf("expected filepath %q, got %q", path, active.FilePath)
	}
	if active.Content != "select 1;\nselect 2;" {
		t.Fatalf("unexpected content: %q", active.Content)
	}
	if active.Dirty {
		t.Fatal("expected buffer to be clean after load")
	}
	if len(status.infos) == 0 {
		t.Fatal("expected status message to be set")
	}
}

func TestEditCommandMissingFileCreatesBuffer(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "missing.sql")

	manager := buffer.NewManager()
	status := &fakeStatus{}
	command := &EditCommand{}
	err := command.Execute(context.Background(), cmdline.Invocation{Name: "e", Args: []string{path}}, &cmdline.CommandContext{
		Buffers: manager,
		Status:  status,
	})
	if err != nil {
		t.Fatalf("edit command: %v", err)
	}

	active := manager.Active()
	if active == nil {
		t.Fatal("expected active buffer")
	}
	if active.FilePath != path {
		t.Fatalf("expected filepath %q, got %q", path, active.FilePath)
	}
	if active.Content != "" {
		t.Fatalf("expected empty content, got %q", active.Content)
	}
	if active.Dirty {
		t.Fatal("expected buffer to be clean for new file")
	}
}

func TestEditCommandAlreadyOpen(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "existing.sql")

	manager := buffer.NewManager()
	buf := manager.Create("existing.sql")
	if err := manager.UpdateMetadata(buf.ID, "existing.sql", path); err != nil {
		t.Fatalf("update metadata: %v", err)
	}

	status := &fakeStatus{}
	command := &EditCommand{}
	err := command.Execute(context.Background(), cmdline.Invocation{Name: "e", Args: []string{path}}, &cmdline.CommandContext{
		Buffers: manager,
		Status:  status,
	})
	if err != nil {
		t.Fatalf("edit command: %v", err)
	}
	if manager.Count() != 1 {
		t.Fatalf("expected 1 buffer, got %d", manager.Count())
	}
	if manager.Active().ID != buf.ID {
		t.Fatalf("expected active buffer %d, got %d", buf.ID, manager.Active().ID)
	}
	if len(status.infos) == 0 {
		t.Fatal("expected status info when buffer already open")
	}
}

func TestEditCommandForceReloadsExisting(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "force.sql")
	if err := os.WriteFile(path, []byte("from-file"), 0o644); err != nil {
		t.Fatalf("write file: %v", err)
	}

	manager := buffer.NewManager()
	buf := manager.Create("force.sql")
	if err := manager.UpdateMetadata(buf.ID, "force.sql", path); err != nil {
		t.Fatalf("update metadata: %v", err)
	}
	if err := manager.UpdateContent(buf.ID, "dirty"); err != nil {
		t.Fatalf("update content: %v", err)
	}

	command := &EditCommand{}
	err := command.Execute(context.Background(), cmdline.Invocation{Name: "e", Args: []string{path}, Force: true}, &cmdline.CommandContext{
		Buffers: manager,
	})
	if err != nil {
		t.Fatalf("edit command: %v", err)
	}
	if manager.Active().Content != "from-file" {
		t.Fatalf("expected content from file, got %q", manager.Active().Content)
	}
	if manager.Active().Dirty {
		t.Fatal("expected buffer to be clean after force reload")
	}
}

func TestEditCommandErrors(t *testing.T) {
	command := &EditCommand{}
	if err := command.Execute(context.Background(), cmdline.Invocation{Name: "e"}, nil); err == nil {
		t.Fatal("expected error for missing buffer manager")
	}

	manager := buffer.NewManager()
	if err := command.Execute(context.Background(), cmdline.Invocation{Name: "e"}, &cmdline.CommandContext{
		Buffers: manager,
	}); err == nil {
		t.Fatal("expected error for missing file name")
	}
}

func TestEditCommandEmptyFileStatus(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "empty.sql")
	if err := os.WriteFile(path, []byte(""), 0o644); err != nil {
		t.Fatalf("write file: %v", err)
	}

	manager := buffer.NewManager()
	status := &fakeStatus{}
	command := &EditCommand{}
	err := command.Execute(context.Background(), cmdline.Invocation{Name: "e", Args: []string{path}}, &cmdline.CommandContext{
		Buffers: manager,
		Status:  status,
	})
	if err != nil {
		t.Fatalf("edit command: %v", err)
	}
	if len(status.infos) == 0 {
		t.Fatal("expected status info for empty file")
	}
}
