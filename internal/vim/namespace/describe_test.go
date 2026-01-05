package namespace

import (
	"errors"
	"testing"
)

type fakeDescribeController struct {
	tableCalls    int
	viewCalls     int
	indexCalls    int
	sequenceCalls int
	funcCalls     int
	dbCalls       int
	cursorCalls   int
	tableErr      error
}

func (f *fakeDescribeController) DescribeTable() error {
	f.tableCalls++
	return f.tableErr
}

func (f *fakeDescribeController) DescribeView() error {
	f.viewCalls++
	return nil
}

func (f *fakeDescribeController) DescribeIndexes() error {
	f.indexCalls++
	return nil
}

func (f *fakeDescribeController) DescribeSequences() error {
	f.sequenceCalls++
	return nil
}

func (f *fakeDescribeController) DescribeFunctions() error {
	f.funcCalls++
	return nil
}

func (f *fakeDescribeController) DescribeDatabase() error {
	f.dbCalls++
	return nil
}

func (f *fakeDescribeController) DescribeUnderCursor() error {
	f.cursorCalls++
	return nil
}

func TestDescribeNamespaceMetadata(t *testing.T) {
	ns := NewDescribeNamespace(&fakeDescribeController{})
	if ns.Prefix() != 'd' {
		t.Fatalf("expected prefix d, got %q", ns.Prefix())
	}
	if ns.Name() != "describe" {
		t.Fatalf("expected name describe, got %q", ns.Name())
	}
}

func TestDescribeNamespaceCommands(t *testing.T) {
	ns := NewDescribeNamespace(&fakeDescribeController{})
	cmds := ns.Commands()
	if len(cmds) == 0 {
		t.Fatal("expected commands")
	}

	found := map[string]bool{}
	for _, cmd := range cmds {
		found[string(cmd.Sequence)] = true
	}

	expect := []string{"t", "v", "i", "s", "f", "d", "\n"}
	for _, key := range expect {
		if !found[key] {
			t.Fatalf("expected command %q", key)
		}
	}
}

func TestDescribeNamespaceHandlers(t *testing.T) {
	controller := &fakeDescribeController{}
	ns := NewDescribeNamespace(controller)

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
				t.Fatalf("expected table call, got %d", controller.tableCalls)
			}
		}},
		{"v", func() {
			if controller.viewCalls != 1 {
				t.Fatalf("expected view call, got %d", controller.viewCalls)
			}
		}},
		{"i", func() {
			if controller.indexCalls != 1 {
				t.Fatalf("expected index call, got %d", controller.indexCalls)
			}
		}},
		{"s", func() {
			if controller.sequenceCalls != 1 {
				t.Fatalf("expected sequence call, got %d", controller.sequenceCalls)
			}
		}},
		{"f", func() {
			if controller.funcCalls != 1 {
				t.Fatalf("expected function call, got %d", controller.funcCalls)
			}
		}},
		{"d", func() {
			if controller.dbCalls != 1 {
				t.Fatalf("expected database call, got %d", controller.dbCalls)
			}
		}},
		{"\n", func() {
			if controller.cursorCalls != 1 {
				t.Fatalf("expected cursor call, got %d", controller.cursorCalls)
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

func TestDescribeNamespaceErrors(t *testing.T) {
	ns := NewDescribeNamespace(nil)
	if err := ns.describeTable(); err == nil {
		t.Fatal("expected table error")
	}
	if err := ns.describeView(); err == nil {
		t.Fatal("expected view error")
	}
	if err := ns.describeIndexes(); err == nil {
		t.Fatal("expected index error")
	}
	if err := ns.describeSequences(); err == nil {
		t.Fatal("expected sequence error")
	}
	if err := ns.describeFunctions(); err == nil {
		t.Fatal("expected function error")
	}
	if err := ns.describeDatabase(); err == nil {
		t.Fatal("expected database error")
	}
	if err := ns.describeUnderCursor(); err == nil {
		t.Fatal("expected cursor error")
	}

	controller := &fakeDescribeController{tableErr: errors.New("boom")}
	ns = NewDescribeNamespace(controller)
	if err := ns.describeTable(); err == nil {
		t.Fatal("expected controller error")
	}
}
