package namespace

import (
	"errors"
	"testing"
)

type fakeShellController struct {
	toggled int
	closed  int
	focused int
	cleared int
	resized int
	history int

	toggleErr  error
	closeErr   error
	focusErr   error
	clearErr   error
	resizeErr  error
	historyErr error
}

func (f *fakeShellController) ToggleShell() error {
	f.toggled++
	return f.toggleErr
}

func (f *fakeShellController) CloseShell() error {
	f.closed++
	return f.closeErr
}

func (f *fakeShellController) FocusShellInput() error {
	f.focused++
	return f.focusErr
}

func (f *fakeShellController) ClearShellOutput() error {
	f.cleared++
	return f.clearErr
}

func (f *fakeShellController) ResizeShell() error {
	f.resized++
	return f.resizeErr
}

func (f *fakeShellController) ShowShellHistory() error {
	f.history++
	return f.historyErr
}

func TestShellNamespaceMetadata(t *testing.T) {
	ns := NewShellNamespace(&fakeShellController{})
	if ns.Prefix() != 's' {
		t.Fatalf("expected prefix s, got %q", ns.Prefix())
	}
	if ns.Name() != "shell" {
		t.Fatalf("expected name shell, got %q", ns.Name())
	}
}

func TestShellNamespaceCommands(t *testing.T) {
	ns := NewShellNamespace(&fakeShellController{})
	cmds := ns.Commands()
	if len(cmds) == 0 {
		t.Fatal("expected commands")
	}

	found := map[string]bool{}
	for _, cmd := range cmds {
		found[string(cmd.Sequence)] = true
	}

	expect := []string{"s", "q", "f", "c", "r", "h"}
	for _, key := range expect {
		if !found[key] {
			t.Fatalf("expected command %q", key)
		}
	}
}

func TestShellNamespaceHandlers(t *testing.T) {
	controller := &fakeShellController{}
	ns := NewShellNamespace(controller)

	handlers := map[string]func() error{}
	for _, cmd := range ns.Commands() {
		handlers[string(cmd.Sequence)] = cmd.Handler
	}

	tests := []struct {
		key    string
		verify func()
	}{
		{"s", func() {
			if controller.toggled != 1 {
				t.Fatalf("expected toggle to be called once, got %d", controller.toggled)
			}
		}},
		{"q", func() {
			if controller.closed != 1 {
				t.Fatalf("expected close to be called once, got %d", controller.closed)
			}
		}},
		{"f", func() {
			if controller.focused != 1 {
				t.Fatalf("expected focus to be called once, got %d", controller.focused)
			}
		}},
		{"c", func() {
			if controller.cleared != 1 {
				t.Fatalf("expected clear to be called once, got %d", controller.cleared)
			}
		}},
		{"r", func() {
			if controller.resized != 1 {
				t.Fatalf("expected resize to be called once, got %d", controller.resized)
			}
		}},
		{"h", func() {
			if controller.history != 1 {
				t.Fatalf("expected history to be called once, got %d", controller.history)
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

func TestShellNamespaceErrors(t *testing.T) {
	ns := NewShellNamespace(nil)
	if err := ns.toggleShell(); err == nil {
		t.Fatal("expected toggle error")
	}
	if err := ns.closeShell(); err == nil {
		t.Fatal("expected close error")
	}
	if err := ns.focusShell(); err == nil {
		t.Fatal("expected focus error")
	}
	if err := ns.clearShell(); err == nil {
		t.Fatal("expected clear error")
	}
	if err := ns.resizeShell(); err == nil {
		t.Fatal("expected resize error")
	}
	if err := ns.showHistory(); err == nil {
		t.Fatal("expected history error")
	}

	controller := &fakeShellController{toggleErr: errors.New("boom")}
	ns = NewShellNamespace(controller)
	if err := ns.toggleShell(); err == nil {
		t.Fatal("expected controller error")
	}
}
