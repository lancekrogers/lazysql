package namespace

import (
	"errors"
	"testing"
)

type fakeRunController struct {
	currentCalls int
	allCalls     int
	linesCalls   int
	paramsCalls  int
	historyCalls int
	currentErr   error
	allErr       error
	linesErr     error
	paramsErr    error
	historyErr   error
}

func (f *fakeRunController) RunCurrentQuery() error {
	f.currentCalls++
	return f.currentErr
}

func (f *fakeRunController) RunAllQueries() error {
	f.allCalls++
	return f.allErr
}

func (f *fakeRunController) RunSelectedLines() error {
	f.linesCalls++
	return f.linesErr
}

func (f *fakeRunController) RunWithParams() error {
	f.paramsCalls++
	return f.paramsErr
}

func (f *fakeRunController) ShowQueryHistory() error {
	f.historyCalls++
	return f.historyErr
}

func TestRunNamespaceMetadata(t *testing.T) {
	ns := NewRunNamespace(&fakeRunController{})
	if ns.Prefix() != 'r' {
		t.Fatalf("expected prefix r, got %q", ns.Prefix())
	}
	if ns.Name() != "run" {
		t.Fatalf("expected name run, got %q", ns.Name())
	}
}

func TestRunNamespaceCommands(t *testing.T) {
	ns := NewRunNamespace(&fakeRunController{})
	cmds := ns.Commands()
	if len(cmds) == 0 {
		t.Fatal("expected commands")
	}

	found := map[string]bool{}
	for _, cmd := range cmds {
		found[string(cmd.Sequence)] = true
	}

	expect := []string{"r", "a", "l", "p", "h"}
	for _, key := range expect {
		if !found[key] {
			t.Fatalf("expected command %q", key)
		}
	}
}

func TestRunNamespaceHandlers(t *testing.T) {
	controller := &fakeRunController{}
	ns := NewRunNamespace(controller)

	handlers := map[string]func() error{}
	for _, cmd := range ns.Commands() {
		handlers[string(cmd.Sequence)] = cmd.Handler
	}

	tests := []struct {
		key    string
		verify func()
	}{
		{"r", func() {
			if controller.currentCalls != 1 {
				t.Fatalf("expected run current call, got %d", controller.currentCalls)
			}
		}},
		{"a", func() {
			if controller.allCalls != 1 {
				t.Fatalf("expected run all call, got %d", controller.allCalls)
			}
		}},
		{"l", func() {
			if controller.linesCalls != 1 {
				t.Fatalf("expected run lines call, got %d", controller.linesCalls)
			}
		}},
		{"p", func() {
			if controller.paramsCalls != 1 {
				t.Fatalf("expected run params call, got %d", controller.paramsCalls)
			}
		}},
		{"h", func() {
			if controller.historyCalls != 1 {
				t.Fatalf("expected history call, got %d", controller.historyCalls)
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

func TestRunNamespaceErrors(t *testing.T) {
	ns := NewRunNamespace(nil)
	if err := ns.runCurrentQuery(); err == nil {
		t.Fatal("expected current query error")
	}
	if err := ns.runAllQueries(); err == nil {
		t.Fatal("expected all queries error")
	}
	if err := ns.runSelectedLines(); err == nil {
		t.Fatal("expected selected lines error")
	}
	if err := ns.runWithParams(); err == nil {
		t.Fatal("expected params error")
	}
	if err := ns.showHistory(); err == nil {
		t.Fatal("expected history error")
	}

	controller := &fakeRunController{currentErr: errors.New("boom")}
	ns = NewRunNamespace(controller)
	if err := ns.runCurrentQuery(); err == nil {
		t.Fatal("expected controller error")
	}
}
