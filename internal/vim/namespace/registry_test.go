package namespace

import (
	"testing"

	"github.com/lancekrogers/lazysql/internal/vim/leader"
)

type fakeNamespace struct {
	prefix   rune
	name     string
	commands []leader.Command
}

func (f fakeNamespace) Prefix() rune {
	return f.prefix
}

func (f fakeNamespace) Name() string {
	return f.name
}

func (f fakeNamespace) Commands() []leader.Command {
	return f.commands
}

func TestRegistryRegisterAndLookup(t *testing.T) {
	leaderRegistry := leader.NewRegistry()
	registry := NewRegistry(leaderRegistry)

	cmd := leader.Command{
		Sequence:    []rune{'n'},
		Description: "next",
		Handler:     func() error { return nil },
	}
	ns := fakeNamespace{
		prefix:   'b',
		name:     "buffer",
		commands: []leader.Command{cmd},
	}

	if err := registry.Register(ns); err != nil {
		t.Fatalf("register: %v", err)
	}

	if _, ok := registry.Get('b'); !ok {
		t.Fatal("expected namespace to be registered")
	}
	if _, ok := registry.Get('x'); ok {
		t.Fatal("did not expect unknown namespace")
	}

	registered, ok := leaderRegistry.Lookup([]rune{'b', 'n'})
	if !ok {
		t.Fatal("expected leader command to be registered")
	}
	if registered.Description != cmd.Description {
		t.Fatalf("unexpected command description: %q", registered.Description)
	}
	if node := leaderRegistry.Tree().Root.Children['b']; node == nil {
		t.Fatal("expected leader group to be registered")
	}
}

func TestRegistryRegisterErrors(t *testing.T) {
	registry := NewRegistry(nil)
	if err := registry.Register(nil); err == nil {
		t.Fatal("expected error for nil namespace")
	}

	ns := fakeNamespace{prefix: 0, name: "missing"}
	if err := registry.Register(ns); err == nil {
		t.Fatal("expected error for missing prefix")
	}

	ns = fakeNamespace{
		prefix: 'b',
		name:   "buffer",
		commands: []leader.Command{
			{Sequence: []rune{'n'}, Description: "next", Handler: func() error { return nil }},
		},
	}
	if err := registry.Register(ns); err != nil {
		t.Fatalf("register: %v", err)
	}
	if err := registry.Register(ns); err == nil {
		t.Fatal("expected duplicate prefix error")
	}
}

func TestRegistryLeaderErrors(t *testing.T) {
	leaderRegistry := leader.NewRegistry()
	registry := NewRegistry(leaderRegistry)

	ns := fakeNamespace{
		prefix: 'b',
		name:   "buffer",
		commands: []leader.Command{
			{Sequence: []rune{'n'}, Description: "next", Handler: func() error { return nil }},
			{Sequence: []rune{'n'}, Description: "duplicate", Handler: func() error { return nil }},
		},
	}
	if err := registry.Register(ns); err == nil {
		t.Fatal("expected error for duplicate command sequence")
	}
}

func TestInitializeRegistry(t *testing.T) {
	leaderRegistry := leader.NewRegistry()
	ns := fakeNamespace{
		prefix: 'b',
		name:   "buffer",
		commands: []leader.Command{
			{Sequence: []rune{'n'}, Description: "next", Handler: func() error { return nil }},
		},
	}

	registry, err := Initialize(leaderRegistry, ns)
	if err != nil {
		t.Fatalf("initialize: %v", err)
	}
	if registry == nil {
		t.Fatal("expected registry instance")
	}

	if _, ok := registry.Get('b'); !ok {
		t.Fatal("expected namespace to be registered")
	}

	if _, err := Initialize(leaderRegistry, nil); err == nil {
		t.Fatal("expected initialize error for nil namespace")
	}
}
