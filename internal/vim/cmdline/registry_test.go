package cmdline

import (
	"context"
	"testing"
)

type fakeCommand struct {
	name    string
	called  bool
	lastInv Invocation
	lastCtx context.Context
	lastEnv *CommandContext
	execErr error
}

func (f *fakeCommand) Name() string {
	return f.name
}

func (f *fakeCommand) Execute(ctx context.Context, inv Invocation, env *CommandContext) error {
	f.called = true
	f.lastInv = inv
	f.lastCtx = ctx
	f.lastEnv = env
	return f.execErr
}

func TestRegistryLookup(t *testing.T) {
	registry := NewRegistry()
	if registry.Lookup("") != nil {
		t.Fatal("expected empty lookup to return nil")
	}
	registry.Register(nil)
	registry.RegisterAlias("", "x")
	registry.RegisterAlias("alias", "")

	cmd := &fakeCommand{name: "write"}
	registry.Register(cmd)
	registry.RegisterAlias("w", "write")

	if registry.Lookup("write") != cmd {
		t.Fatal("expected lookup by name to return command")
	}
	if registry.Lookup("w") != cmd {
		t.Fatal("expected lookup by alias to return command")
	}
	if registry.Lookup("missing") != nil {
		t.Fatal("expected missing lookup to return nil")
	}
}

func TestCommandExecutor(t *testing.T) {
	registry := NewRegistry()
	cmd := &fakeCommand{name: "echo"}
	registry.Register(cmd)
	registry.RegisterAlias("e", "echo")

	exec := NewCommandExecutor(registry, &CommandContext{})
	if err := exec.Execute(context.Background(), ""); err != nil {
		t.Fatalf("expected empty command to be ignored, got %v", err)
	}

	if err := exec.Execute(nil, "e hello"); err != nil {
		t.Fatalf("execute command: %v", err)
	}
	if !cmd.called {
		t.Fatal("expected command to execute")
	}
	if cmd.lastCtx == nil {
		t.Fatal("expected context to be set")
	}
	if cmd.lastInv.Name != "e" {
		t.Fatalf("expected invocation name to be 'e', got %q", cmd.lastInv.Name)
	}
	if len(cmd.lastInv.Args) != 1 || cmd.lastInv.Args[0] != "hello" {
		t.Fatalf("unexpected invocation args: %v", cmd.lastInv.Args)
	}
}

func TestCommandExecutorErrors(t *testing.T) {
	exec := NewCommandExecutor(nil, nil)
	if err := exec.Execute(context.Background(), "w"); err == nil {
		t.Fatal("expected error for missing registry")
	}

	registry := NewRegistry()
	exec = NewCommandExecutor(registry, &CommandContext{})
	if err := exec.Execute(context.Background(), "unknown"); err == nil {
		t.Fatal("expected error for unknown command")
	}

	if err := exec.Execute(context.Background(), `e "unterminated`); err == nil {
		t.Fatal("expected error for unterminated quote")
	}
}
