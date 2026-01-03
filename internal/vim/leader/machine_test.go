package leader

import (
	"testing"

	"github.com/lancekrogers/lazysql/internal/vim/whichkey"
)

func TestLeaderStateMachineTransitions(t *testing.T) {
	tree := whichkey.NewKeyTree()
	called := false
	tree.AddGroup([]rune{'f'}, "find")
	tree.AddCommand([]rune{'f', 't'}, "table", func() {
		called = true
	})

	machine := NewStateMachine(tree)
	if machine.State() != StateIdle {
		t.Fatalf("expected idle state, got %v", machine.State())
	}

	if !machine.Activate() {
		t.Fatal("expected Activate to succeed")
	}
	if machine.State() != StateActive {
		t.Fatalf("expected active state, got %v", machine.State())
	}

	result := machine.AddKey('f')
	if result.Error != nil || result.NewState != StateSequence {
		t.Fatalf("expected sequence state, got %v (err=%v)", result.NewState, result.Error)
	}

	result = machine.AddKey('t')
	if result.Error != nil || result.NewState != StateExecuting {
		t.Fatalf("expected executing state, got %v (err=%v)", result.NewState, result.Error)
	}
	if result.Action == nil {
		t.Fatal("expected action for leaf node")
	}
	result.Action()
	if !called {
		t.Fatal("expected action to be executed")
	}

	machine.Complete()
	if machine.State() != StateIdle {
		t.Fatalf("expected idle after completion, got %v", machine.State())
	}
	if len(machine.Sequence()) != 0 {
		t.Fatal("expected sequence to be cleared")
	}
}

func TestLeaderStateMachineInvalidKey(t *testing.T) {
	tree := whichkey.NewKeyTree()
	tree.AddGroup([]rune{'f'}, "find")

	machine := NewStateMachine(tree)
	machine.Activate()
	result := machine.AddKey('x')
	if result.Error != ErrInvalidKey {
		t.Fatalf("expected invalid key error, got %v", result.Error)
	}
	if machine.State() != StateIdle {
		t.Fatalf("expected idle after invalid key, got %v", machine.State())
	}
}

func TestLeaderStateMachineInactive(t *testing.T) {
	tree := whichkey.NewKeyTree()
	machine := NewStateMachine(tree)
	result := machine.AddKey('f')
	if result.Error != ErrInactive {
		t.Fatalf("expected inactive error, got %v", result.Error)
	}
}

func TestLeaderStateMachineNoTree(t *testing.T) {
	machine := NewStateMachine(nil)
	machine.Activate()
	result := machine.AddKey('f')
	if result.Error != ErrNoTree {
		t.Fatalf("expected no-tree error, got %v", result.Error)
	}
}

func TestLeaderStateMachineExecutingError(t *testing.T) {
	tree := whichkey.NewKeyTree()
	tree.AddCommand([]rune{'x'}, "execute", func() {})

	machine := NewStateMachine(tree)
	machine.Activate()
	result := machine.AddKey('x')
	if result.Error != nil || result.NewState != StateExecuting {
		t.Fatalf("expected executing state, got %v (err=%v)", result.NewState, result.Error)
	}

	result = machine.AddKey('x')
	if result.Error != ErrExecuting {
		t.Fatalf("expected executing error, got %v", result.Error)
	}
}

func TestLeaderStateMachineCancel(t *testing.T) {
	tree := whichkey.NewKeyTree()
	tree.AddGroup([]rune{'f'}, "find")

	machine := NewStateMachine(tree)
	machine.Activate()
	machine.AddKey('f')
	machine.Cancel()
	if machine.State() != StateIdle {
		t.Fatalf("expected idle after cancel, got %v", machine.State())
	}
}

func TestLeaderStateMachineActivateWhenActive(t *testing.T) {
	tree := whichkey.NewKeyTree()
	machine := NewStateMachine(tree)
	if !machine.Activate() {
		t.Fatal("expected first activate to succeed")
	}
	if machine.Activate() {
		t.Fatal("expected second activate to be ignored")
	}
}

func TestLeaderStateString(t *testing.T) {
	if StateIdle.String() != "Idle" {
		t.Fatalf("unexpected StateIdle string: %s", StateIdle.String())
	}
	if StateActive.String() != "Active" {
		t.Fatalf("unexpected StateActive string: %s", StateActive.String())
	}
	if StateSequence.String() != "Sequence" {
		t.Fatalf("unexpected StateSequence string: %s", StateSequence.String())
	}
	if StateExecuting.String() != "Executing" {
		t.Fatalf("unexpected StateExecuting string: %s", StateExecuting.String())
	}
	if LeaderState(99).String() != "Unknown" {
		t.Fatalf("unexpected unknown state string: %s", LeaderState(99).String())
	}
}
