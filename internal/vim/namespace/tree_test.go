package namespace

import (
	"errors"
	"testing"
)

type fakeTreeController struct {
	toggleCalls    int
	expandCalls    int
	collapseCalls  int
	expandAll      int
	collapseAll    int
	focusCalls     int
	refreshCalls   int
	syncCalls      int
	toggleErr      error
	expandErr      error
	collapseErr    error
	expandAllErr   error
	collapseAllErr error
	focusErr       error
	refreshErr     error
	syncErr        error
}

func (f *fakeTreeController) ToggleTree() error {
	f.toggleCalls++
	return f.toggleErr
}

func (f *fakeTreeController) ExpandNode() error {
	f.expandCalls++
	return f.expandErr
}

func (f *fakeTreeController) CollapseNode() error {
	f.collapseCalls++
	return f.collapseErr
}

func (f *fakeTreeController) ExpandAll() error {
	f.expandAll++
	return f.expandAllErr
}

func (f *fakeTreeController) CollapseAll() error {
	f.collapseAll++
	return f.collapseAllErr
}

func (f *fakeTreeController) FocusTree() error {
	f.focusCalls++
	return f.focusErr
}

func (f *fakeTreeController) RefreshTree() error {
	f.refreshCalls++
	return f.refreshErr
}

func (f *fakeTreeController) SyncWithBuffer() error {
	f.syncCalls++
	return f.syncErr
}

func TestTreeNamespaceMetadata(t *testing.T) {
	ns := NewTreeNamespace(&fakeTreeController{})
	if ns.Prefix() != 't' {
		t.Fatalf("expected prefix t, got %q", ns.Prefix())
	}
	if ns.Name() != "tree" {
		t.Fatalf("expected name tree, got %q", ns.Name())
	}
}

func TestTreeNamespaceCommands(t *testing.T) {
	ns := NewTreeNamespace(&fakeTreeController{})
	cmds := ns.Commands()
	if len(cmds) == 0 {
		t.Fatal("expected commands")
	}

	found := map[string]bool{}
	for _, cmd := range cmds {
		found[string(cmd.Sequence)] = true
	}

	expect := []string{"t", "e", "c", "E", "C", "f", "r", "s"}
	for _, key := range expect {
		if !found[key] {
			t.Fatalf("expected command %q", key)
		}
	}
}

func TestTreeNamespaceHandlers(t *testing.T) {
	controller := &fakeTreeController{}
	ns := NewTreeNamespace(controller)

	handlers := map[string]func() error{}
	for _, cmd := range ns.Commands() {
		handlers[string(cmd.Sequence)] = cmd.Handler
	}

	tests := []struct {
		key    string
		verify func()
	}{
		{"t", func() {
			if controller.toggleCalls != 1 {
				t.Fatalf("expected toggle call, got %d", controller.toggleCalls)
			}
		}},
		{"e", func() {
			if controller.expandCalls != 1 {
				t.Fatalf("expected expand call, got %d", controller.expandCalls)
			}
		}},
		{"c", func() {
			if controller.collapseCalls != 1 {
				t.Fatalf("expected collapse call, got %d", controller.collapseCalls)
			}
		}},
		{"E", func() {
			if controller.expandAll != 1 {
				t.Fatalf("expected expand all call, got %d", controller.expandAll)
			}
		}},
		{"C", func() {
			if controller.collapseAll != 1 {
				t.Fatalf("expected collapse all call, got %d", controller.collapseAll)
			}
		}},
		{"f", func() {
			if controller.focusCalls != 1 {
				t.Fatalf("expected focus call, got %d", controller.focusCalls)
			}
		}},
		{"r", func() {
			if controller.refreshCalls != 1 {
				t.Fatalf("expected refresh call, got %d", controller.refreshCalls)
			}
		}},
		{"s", func() {
			if controller.syncCalls != 1 {
				t.Fatalf("expected sync call, got %d", controller.syncCalls)
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

func TestTreeNamespaceErrors(t *testing.T) {
	ns := NewTreeNamespace(nil)
	if err := ns.toggleTree(); err == nil {
		t.Fatal("expected toggle error")
	}
	if err := ns.expandNode(); err == nil {
		t.Fatal("expected expand error")
	}
	if err := ns.collapseNode(); err == nil {
		t.Fatal("expected collapse error")
	}
	if err := ns.expandAll(); err == nil {
		t.Fatal("expected expand all error")
	}
	if err := ns.collapseAll(); err == nil {
		t.Fatal("expected collapse all error")
	}
	if err := ns.focusTree(); err == nil {
		t.Fatal("expected focus error")
	}
	if err := ns.refreshTree(); err == nil {
		t.Fatal("expected refresh error")
	}
	if err := ns.syncWithBuffer(); err == nil {
		t.Fatal("expected sync error")
	}

	controller := &fakeTreeController{toggleErr: errors.New("boom")}
	ns = NewTreeNamespace(controller)
	if err := ns.toggleTree(); err == nil {
		t.Fatal("expected controller error")
	}
}
