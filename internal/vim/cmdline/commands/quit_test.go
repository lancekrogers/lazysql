package commands

import (
	"context"
	"testing"

	"github.com/lancekrogers/lazysql/internal/vim/buffer"
	"github.com/lancekrogers/lazysql/internal/vim/cmdline"
)

func TestQuitCommandDirtyBuffer(t *testing.T) {
	manager := buffer.NewManager()
	buf := manager.Create("one")
	if err := manager.UpdateContent(buf.ID, "select 1"); err != nil {
		t.Fatalf("update content: %v", err)
	}

	app := &fakeApp{}
	command := &QuitCommand{}
	err := command.Execute(context.Background(), cmdline.Invocation{Name: "q"}, &cmdline.CommandContext{
		Buffers: manager,
		App:     app,
	})
	if err == nil {
		t.Fatal("expected error for dirty buffer")
	}
	if app.stopped {
		t.Fatal("did not expect app to stop on dirty buffer")
	}
}

func TestQuitCommandForce(t *testing.T) {
	manager := buffer.NewManager()
	buf := manager.Create("one")
	if err := manager.UpdateContent(buf.ID, "select 1"); err != nil {
		t.Fatalf("update content: %v", err)
	}

	app := &fakeApp{}
	command := &QuitCommand{}
	if command.Name() != "q" {
		t.Fatalf("expected command name q, got %q", command.Name())
	}
	err := command.Execute(context.Background(), cmdline.Invocation{Name: "q", Force: true}, &cmdline.CommandContext{
		Buffers: manager,
		App:     app,
	})
	if err != nil {
		t.Fatalf("quit command: %v", err)
	}
	if manager.Count() != 0 {
		t.Fatalf("expected buffers closed, got %d", manager.Count())
	}
	if !app.stopped {
		t.Fatal("expected app to stop when last buffer closes")
	}
}

func TestQuitCommandNoActiveBuffer(t *testing.T) {
	manager := buffer.NewManager()
	app := &fakeApp{}
	command := &QuitCommand{}
	err := command.Execute(context.Background(), cmdline.Invocation{Name: "q"}, &cmdline.CommandContext{
		Buffers: manager,
		App:     app,
	})
	if err != nil {
		t.Fatalf("quit command: %v", err)
	}
	if !app.stopped {
		t.Fatal("expected app to stop when no buffers exist")
	}
}

func TestQuitAllCommand(t *testing.T) {
	manager := buffer.NewManager()
	first := manager.Create("one")
	second := manager.Create("two")
	if err := manager.UpdateContent(first.ID, "select 1"); err != nil {
		t.Fatalf("update content: %v", err)
	}
	if err := manager.UpdateContent(second.ID, "select 2"); err != nil {
		t.Fatalf("update content: %v", err)
	}

	app := &fakeApp{}
	command := &QuitAllCommand{}
	if command.Name() != "qa" {
		t.Fatalf("expected command name qa, got %q", command.Name())
	}
	err := command.Execute(context.Background(), cmdline.Invocation{Name: "qa"}, &cmdline.CommandContext{
		Buffers: manager,
		App:     app,
	})
	if err == nil {
		t.Fatal("expected error for dirty buffers")
	}

	err = command.Execute(context.Background(), cmdline.Invocation{Name: "qa", Force: true}, &cmdline.CommandContext{
		Buffers: manager,
		App:     app,
	})
	if err != nil {
		t.Fatalf("quit all: %v", err)
	}
	if manager.Count() != 0 {
		t.Fatalf("expected buffers closed, got %d", manager.Count())
	}
	if !app.stopped {
		t.Fatal("expected app to stop after quit all")
	}
}

func TestWriteQuitCommand(t *testing.T) {
	dir := t.TempDir()
	path := dir + "/query.sql"

	manager := buffer.NewManager()
	buf := manager.Create("one")
	if err := manager.UpdateContent(buf.ID, "select 1"); err != nil {
		t.Fatalf("update content: %v", err)
	}

	app := &fakeApp{}
	command := &WriteQuitCommand{}
	if command.Name() != "wq" {
		t.Fatalf("expected command name wq, got %q", command.Name())
	}
	err := command.Execute(context.Background(), cmdline.Invocation{Name: "wq", Args: []string{path}}, &cmdline.CommandContext{
		Buffers: manager,
		App:     app,
	})
	if err != nil {
		t.Fatalf("write quit: %v", err)
	}
	if manager.Count() != 0 {
		t.Fatalf("expected buffers closed, got %d", manager.Count())
	}
	if !app.stopped {
		t.Fatal("expected app to stop after write quit")
	}
}

func TestQuitAllCommandCleanBuffers(t *testing.T) {
	manager := buffer.NewManager()
	manager.Create("one")
	manager.Create("two")

	app := &fakeApp{}
	command := &QuitAllCommand{}
	err := command.Execute(context.Background(), cmdline.Invocation{Name: "qa"}, &cmdline.CommandContext{
		Buffers: manager,
		App:     app,
	})
	if err != nil {
		t.Fatalf("quit all: %v", err)
	}
	if manager.Count() != 0 {
		t.Fatalf("expected buffers closed, got %d", manager.Count())
	}
	if !app.stopped {
		t.Fatal("expected app to stop after quit all")
	}
}

func TestQuitCommandErrors(t *testing.T) {
	command := &QuitCommand{}
	if err := command.Execute(context.Background(), cmdline.Invocation{Name: "q"}, nil); err == nil {
		t.Fatal("expected error for missing buffer manager")
	}

	manager := buffer.NewManager()
	if err := command.Execute(context.Background(), cmdline.Invocation{Name: "q"}, &cmdline.CommandContext{
		Buffers: manager,
	}); err == nil {
		t.Fatal("expected error for missing app controller")
	}
}

func TestQuitAllCommandErrors(t *testing.T) {
	command := &QuitAllCommand{}
	if err := command.Execute(context.Background(), cmdline.Invocation{Name: "qa"}, nil); err == nil {
		t.Fatal("expected error for missing buffer manager")
	}

	manager := buffer.NewManager()
	if err := command.Execute(context.Background(), cmdline.Invocation{Name: "qa"}, &cmdline.CommandContext{
		Buffers: manager,
	}); err == nil {
		t.Fatal("expected error for missing app controller")
	}
}
