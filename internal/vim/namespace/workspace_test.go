package namespace

import (
	"errors"
	"testing"
)

type fakeWorkspaceController struct {
	connectCalls    int
	switchCalls     int
	listCalls       int
	disconnectCalls int
	connectErr      error
	switchErr       error
	listErr         error
	disconnectErr   error
}

func (f *fakeWorkspaceController) Connect() error {
	f.connectCalls++
	return f.connectErr
}

func (f *fakeWorkspaceController) SwitchConnection() error {
	f.switchCalls++
	return f.switchErr
}

func (f *fakeWorkspaceController) ListConnections() error {
	f.listCalls++
	return f.listErr
}

func (f *fakeWorkspaceController) Disconnect() error {
	f.disconnectCalls++
	return f.disconnectErr
}

func TestWorkspaceNamespaceMetadata(t *testing.T) {
	ns := NewWorkspaceNamespace(&fakeWorkspaceController{})
	if ns.Prefix() != 'w' {
		t.Fatalf("expected prefix w, got %q", ns.Prefix())
	}
	if ns.Name() != "workspace" {
		t.Fatalf("expected name workspace, got %q", ns.Name())
	}
}

func TestWorkspaceNamespaceCommands(t *testing.T) {
	ns := NewWorkspaceNamespace(&fakeWorkspaceController{})
	cmds := ns.Commands()
	if len(cmds) == 0 {
		t.Fatal("expected commands")
	}

	found := map[string]bool{}
	for _, cmd := range cmds {
		found[string(cmd.Sequence)] = true
	}

	expect := []string{"c", "s", "l", "d"}
	for _, key := range expect {
		if !found[key] {
			t.Fatalf("expected command %q", key)
		}
	}
}

func TestWorkspaceNamespaceHandlers(t *testing.T) {
	controller := &fakeWorkspaceController{}
	ns := NewWorkspaceNamespace(controller)

	handlers := map[string]func() error{}
	for _, cmd := range ns.Commands() {
		handlers[string(cmd.Sequence)] = cmd.Handler
	}

	tests := []struct {
		key    string
		verify func()
	}{
		{"c", func() {
			if controller.connectCalls != 1 {
				t.Fatalf("expected connect call, got %d", controller.connectCalls)
			}
		}},
		{"s", func() {
			if controller.switchCalls != 1 {
				t.Fatalf("expected switch call, got %d", controller.switchCalls)
			}
		}},
		{"l", func() {
			if controller.listCalls != 1 {
				t.Fatalf("expected list call, got %d", controller.listCalls)
			}
		}},
		{"d", func() {
			if controller.disconnectCalls != 1 {
				t.Fatalf("expected disconnect call, got %d", controller.disconnectCalls)
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

func TestWorkspaceNamespaceErrors(t *testing.T) {
	ns := NewWorkspaceNamespace(nil)
	if err := ns.connect(); err == nil {
		t.Fatal("expected connect error")
	}
	if err := ns.switchConnection(); err == nil {
		t.Fatal("expected switch error")
	}
	if err := ns.listConnections(); err == nil {
		t.Fatal("expected list error")
	}
	if err := ns.disconnect(); err == nil {
		t.Fatal("expected disconnect error")
	}

	controller := &fakeWorkspaceController{connectErr: errors.New("boom")}
	ns = NewWorkspaceNamespace(controller)
	if err := ns.connect(); err == nil {
		t.Fatal("expected controller error")
	}
}
