package commands

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/lancekrogers/lazysql/internal/vim/buffer"
	"github.com/lancekrogers/lazysql/internal/vim/cmdline"
)

func TestWriteCommandWritesFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "query.sql")

	manager := buffer.NewManager()
	buf := manager.Create("query")
	if err := manager.UpdateContent(buf.ID, "select 1;"); err != nil {
		t.Fatalf("update content: %v", err)
	}

	status := &fakeStatus{}
	command := &WriteCommand{}
	if command.Name() != "w" {
		t.Fatalf("expected command name w, got %q", command.Name())
	}
	err := command.Execute(context.Background(), cmdline.Invocation{Name: "w", Args: []string{path}}, &cmdline.CommandContext{
		Buffers: manager,
		Status:  status,
	})
	if err != nil {
		t.Fatalf("write command: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read file: %v", err)
	}
	if string(data) != "select 1;" {
		t.Fatalf("unexpected file content: %q", string(data))
	}
	if manager.Active().FilePath != path {
		t.Fatalf("expected filepath to be %q, got %q", path, manager.Active().FilePath)
	}
	if manager.Active().Dirty {
		t.Fatal("expected buffer to be clean after write")
	}
	if len(status.infos) == 0 {
		t.Fatal("expected status message to be set")
	}
}

func TestWriteCommandBlocksOverwrite(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "existing.sql")
	if err := os.WriteFile(path, []byte("existing"), 0o644); err != nil {
		t.Fatalf("write existing file: %v", err)
	}

	manager := buffer.NewManager()
	buf := manager.Create("query")
	if err := manager.UpdateContent(buf.ID, "select 2;"); err != nil {
		t.Fatalf("update content: %v", err)
	}

	command := &WriteCommand{}
	err := command.Execute(context.Background(), cmdline.Invocation{Name: "w", Args: []string{path}}, &cmdline.CommandContext{
		Buffers: manager,
	})
	if err == nil {
		t.Fatal("expected error when overwriting without force")
	}
}

func TestWriteCommandForceOverwrite(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "existing.sql")
	if err := os.WriteFile(path, []byte("existing"), 0o644); err != nil {
		t.Fatalf("write existing file: %v", err)
	}

	manager := buffer.NewManager()
	buf := manager.Create("query")
	if err := manager.UpdateContent(buf.ID, "select 3;"); err != nil {
		t.Fatalf("update content: %v", err)
	}

	command := &WriteCommand{}
	err := command.Execute(context.Background(), cmdline.Invocation{Name: "w", Args: []string{path}, Force: true}, &cmdline.CommandContext{
		Buffers: manager,
	})
	if err != nil {
		t.Fatalf("write command: %v", err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read file: %v", err)
	}
	if string(data) != "select 3;" {
		t.Fatalf("expected overwritten content, got %q", string(data))
	}
}

func TestWriteCommandErrors(t *testing.T) {
	command := &WriteCommand{}
	if err := command.Execute(context.Background(), cmdline.Invocation{Name: "w"}, nil); err == nil {
		t.Fatal("expected error for missing buffer manager")
	}

	manager := buffer.NewManager()
	if err := command.Execute(context.Background(), cmdline.Invocation{Name: "w"}, &cmdline.CommandContext{
		Buffers: manager,
	}); err == nil {
		t.Fatal("expected error for missing active buffer")
	}

	buf := manager.Create("query")
	if err := manager.UpdateContent(buf.ID, "select 1"); err != nil {
		t.Fatalf("update content: %v", err)
	}
	if err := command.Execute(context.Background(), cmdline.Invocation{Name: "w"}, &cmdline.CommandContext{
		Buffers: manager,
	}); err == nil {
		t.Fatal("expected error for missing file name")
	}
}
