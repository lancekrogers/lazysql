package namespace

import (
	"errors"
	"testing"
)

type fakeExplainController struct {
	explainCalls      int
	analyzeCalls      int
	explainErr        error
	explainAnalyzeErr error
}

func (f *fakeExplainController) ExplainCurrentQuery() error {
	f.explainCalls++
	return f.explainErr
}

func (f *fakeExplainController) ExplainAnalyzeCurrentQuery() error {
	f.analyzeCalls++
	return f.explainAnalyzeErr
}

func TestExplainNamespaceMetadata(t *testing.T) {
	ns := NewExplainNamespace(&fakeExplainController{})
	if ns.Prefix() != 'x' {
		t.Fatalf("expected prefix x, got %q", ns.Prefix())
	}
	if ns.Name() != "explain" {
		t.Fatalf("expected name explain, got %q", ns.Name())
	}
}

func TestExplainNamespaceCommands(t *testing.T) {
	ns := NewExplainNamespace(&fakeExplainController{})
	cmds := ns.Commands()
	if len(cmds) == 0 {
		t.Fatal("expected commands")
	}

	found := map[string]bool{}
	for _, cmd := range cmds {
		found[string(cmd.Sequence)] = true
	}

	expect := []string{"x", "X"}
	for _, key := range expect {
		if !found[key] {
			t.Fatalf("expected command %q", key)
		}
	}
}

func TestExplainNamespaceHandlers(t *testing.T) {
	controller := &fakeExplainController{}
	ns := NewExplainNamespace(controller)

	handlers := map[string]func() error{}
	for _, cmd := range ns.Commands() {
		handlers[string(cmd.Sequence)] = cmd.Handler
	}

	tests := []struct {
		key    string
		verify func()
	}{
		{"x", func() {
			if controller.explainCalls != 1 {
				t.Fatalf("expected explain call, got %d", controller.explainCalls)
			}
		}},
		{"X", func() {
			if controller.analyzeCalls != 1 {
				t.Fatalf("expected analyze call, got %d", controller.analyzeCalls)
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

func TestExplainNamespaceErrors(t *testing.T) {
	ns := NewExplainNamespace(nil)
	if err := ns.explain(); err == nil {
		t.Fatal("expected explain error")
	}
	if err := ns.explainAnalyze(); err == nil {
		t.Fatal("expected analyze error")
	}

	controller := &fakeExplainController{explainErr: errors.New("boom")}
	ns = NewExplainNamespace(controller)
	if err := ns.explain(); err == nil {
		t.Fatal("expected controller error")
	}
}
