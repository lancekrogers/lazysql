package namespace

import (
	"errors"
	"testing"
)

type fakeCancelController struct {
	cancelCalls int
	clearCalls  int
	cancelErr   error
	clearErr    error
}

func (f *fakeCancelController) CancelQuery() error {
	f.cancelCalls++
	return f.cancelErr
}

func (f *fakeCancelController) ClearResults() error {
	f.clearCalls++
	return f.clearErr
}

func TestCancelNamespaceMetadata(t *testing.T) {
	ns := NewCancelNamespace(&fakeCancelController{})
	if ns.Prefix() != 'c' {
		t.Fatalf("expected prefix c, got %q", ns.Prefix())
	}
	if ns.Name() != "cancel" {
		t.Fatalf("expected name cancel, got %q", ns.Name())
	}
}

func TestCancelNamespaceCommands(t *testing.T) {
	ns := NewCancelNamespace(&fakeCancelController{})
	cmds := ns.Commands()
	if len(cmds) == 0 {
		t.Fatal("expected commands")
	}

	found := map[string]bool{}
	for _, cmd := range cmds {
		found[string(cmd.Sequence)] = true
	}

	expect := []string{"c", "l"}
	for _, key := range expect {
		if !found[key] {
			t.Fatalf("expected command %q", key)
		}
	}
}

func TestCancelNamespaceHandlers(t *testing.T) {
	controller := &fakeCancelController{}
	ns := NewCancelNamespace(controller)

	handlers := map[string]func() error{}
	for _, cmd := range ns.Commands() {
		handlers[string(cmd.Sequence)] = cmd.Handler
	}

	tests := []struct {
		key    string
		verify func()
	}{
		{"c", func() {
			if controller.cancelCalls != 1 {
				t.Fatalf("expected cancel call, got %d", controller.cancelCalls)
			}
		}},
		{"l", func() {
			if controller.clearCalls != 1 {
				t.Fatalf("expected clear call, got %d", controller.clearCalls)
			}
		}},
	}

	for _, test := range tests {
		handler := handlers[test.key]
		if handler == nil {
			t.Fatalf("missing handler for %q", test.key)
		}
		if err := handler(); err != nil {
			t.Fatalf("handler %q returned error: %v", test.key, err)
		}
		test.verify()
	}
}

func TestCancelNamespaceErrors(t *testing.T) {
	ns := NewCancelNamespace(nil)
	if err := ns.cancelQuery(); err == nil {
		t.Fatal("expected cancel error")
	}
	if err := ns.clearResults(); err == nil {
		t.Fatal("expected clear error")
	}

	controller := &fakeCancelController{cancelErr: errors.New("boom")}
	ns = NewCancelNamespace(controller)
	if err := ns.cancelQuery(); err == nil {
		t.Fatal("expected controller error")
	}
}
