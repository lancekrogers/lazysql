package namespace

import (
	"errors"
	"testing"
)

type fakeFindController struct {
	tableCalls    int
	columnCalls   int
	functionCalls int
	viewCalls     int
	schemaCalls   int
	queryCalls    int
	recentCalls   int
	tableErr      error
	columnErr     error
	functionErr   error
	viewErr       error
	schemaErr     error
	queryErr      error
	recentErr     error
}

func (f *fakeFindController) FindTable() error {
	f.tableCalls++
	return f.tableErr
}

func (f *fakeFindController) FindColumn() error {
	f.columnCalls++
	return f.columnErr
}

func (f *fakeFindController) FindFunction() error {
	f.functionCalls++
	return f.functionErr
}

func (f *fakeFindController) FindView() error {
	f.viewCalls++
	return f.viewErr
}

func (f *fakeFindController) FindSchema() error {
	f.schemaCalls++
	return f.schemaErr
}

func (f *fakeFindController) FindInQueries() error {
	f.queryCalls++
	return f.queryErr
}

func (f *fakeFindController) FindRecent() error {
	f.recentCalls++
	return f.recentErr
}

func TestFindNamespaceMetadata(t *testing.T) {
	ns := NewFindNamespace(&fakeFindController{})
	if ns.Prefix() != 'f' {
		t.Fatalf("expected prefix f, got %q", ns.Prefix())
	}
	if ns.Name() != "find" {
		t.Fatalf("expected name find, got %q", ns.Name())
	}
}

func TestFindNamespaceCommands(t *testing.T) {
	ns := NewFindNamespace(&fakeFindController{})
	cmds := ns.Commands()
	if len(cmds) == 0 {
		t.Fatal("expected commands")
	}

	found := map[string]bool{}
	for _, cmd := range cmds {
		found[string(cmd.Sequence)] = true
	}

	expect := []string{"t", "c", "f", "v", "s", "q", "r"}
	for _, key := range expect {
		if !found[key] {
			t.Fatalf("expected command %q", key)
		}
	}
}

func TestFindNamespaceHandlers(t *testing.T) {
	controller := &fakeFindController{}
	ns := NewFindNamespace(controller)

	handlers := map[string]func() error{}
	for _, cmd := range ns.Commands() {
		handlers[string(cmd.Sequence)] = cmd.Handler
	}

	tests := []struct {
		key    string
		verify func()
	}{
		{"t", func() {
			if controller.tableCalls != 1 {
				t.Fatalf("expected find table call, got %d", controller.tableCalls)
			}
		}},
		{"c", func() {
			if controller.columnCalls != 1 {
				t.Fatalf("expected find column call, got %d", controller.columnCalls)
			}
		}},
		{"f", func() {
			if controller.functionCalls != 1 {
				t.Fatalf("expected find function call, got %d", controller.functionCalls)
			}
		}},
		{"v", func() {
			if controller.viewCalls != 1 {
				t.Fatalf("expected find view call, got %d", controller.viewCalls)
			}
		}},
		{"s", func() {
			if controller.schemaCalls != 1 {
				t.Fatalf("expected find schema call, got %d", controller.schemaCalls)
			}
		}},
		{"q", func() {
			if controller.queryCalls != 1 {
				t.Fatalf("expected find in queries call, got %d", controller.queryCalls)
			}
		}},
		{"r", func() {
			if controller.recentCalls != 1 {
				t.Fatalf("expected find recent call, got %d", controller.recentCalls)
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

func TestFindNamespaceErrors(t *testing.T) {
	ns := NewFindNamespace(nil)
	if err := ns.findTable(); err == nil {
		t.Fatal("expected table error")
	}
	if err := ns.findColumn(); err == nil {
		t.Fatal("expected column error")
	}
	if err := ns.findFunction(); err == nil {
		t.Fatal("expected function error")
	}
	if err := ns.findView(); err == nil {
		t.Fatal("expected view error")
	}
	if err := ns.findSchema(); err == nil {
		t.Fatal("expected schema error")
	}
	if err := ns.findInQueries(); err == nil {
		t.Fatal("expected query error")
	}
	if err := ns.findRecent(); err == nil {
		t.Fatal("expected recent error")
	}

	controller := &fakeFindController{tableErr: errors.New("boom")}
	ns = NewFindNamespace(controller)
	if err := ns.findTable(); err == nil {
		t.Fatal("expected controller error")
	}
}
