package leader_test

import (
	"errors"
	"testing"
	"time"

	"github.com/gdamore/tcell/v2"

	"github.com/lancekrogers/lazysql/internal/vim/leader"
	"github.com/lancekrogers/lazysql/internal/vim/namespace"
)

type fakeFindController struct {
	tableCalls int
	tableErr   error
}

func (f *fakeFindController) FindTable() error {
	f.tableCalls++
	return f.tableErr
}

func (f *fakeFindController) FindColumn() error   { return nil }
func (f *fakeFindController) FindFunction() error { return nil }
func (f *fakeFindController) FindView() error     { return nil }
func (f *fakeFindController) FindSchema() error   { return nil }
func (f *fakeFindController) FindInQueries() error {
	return nil
}
func (f *fakeFindController) FindRecent() error { return nil }

type fakeTreeController struct {
	toggleCalls int
	toggleErr   error
}

func (f *fakeTreeController) ToggleTree() error {
	f.toggleCalls++
	return f.toggleErr
}

func (f *fakeTreeController) ExpandNode() error   { return nil }
func (f *fakeTreeController) CollapseNode() error { return nil }
func (f *fakeTreeController) ExpandAll() error    { return nil }
func (f *fakeTreeController) CollapseAll() error  { return nil }
func (f *fakeTreeController) FocusTree() error    { return nil }
func (f *fakeTreeController) RefreshTree() error  { return nil }
func (f *fakeTreeController) SyncWithBuffer() error {
	return nil
}

type fakeWorkspaceController struct {
	connectCalls int
	connectErr   error
}

func (f *fakeWorkspaceController) Connect() error {
	f.connectCalls++
	return f.connectErr
}

func (f *fakeWorkspaceController) SwitchConnection() error { return nil }
func (f *fakeWorkspaceController) ListConnections() error  { return nil }
func (f *fakeWorkspaceController) Disconnect() error       { return nil }

type fakeExplainController struct {
	explainCalls int
	analyzeCalls int
	explainErr   error
	analyzeErr   error
}

func (f *fakeExplainController) ExplainCurrentQuery() error {
	f.explainCalls++
	return f.explainErr
}

func (f *fakeExplainController) ExplainAnalyzeCurrentQuery() error {
	f.analyzeCalls++
	return f.analyzeErr
}

type fakeCancelController struct {
	cancelCalls int
	cancelErr   error
	clearCalls  int
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

type fakeStatusReporter struct {
	errors []string
}

func (f *fakeStatusReporter) Info(string) {}

func (f *fakeStatusReporter) Error(message string) {
	f.errors = append(f.errors, message)
}

func sendRunes(manager *leader.Manager, runes ...rune) {
	for _, r := range runes {
		manager.HandleEvent(tcell.NewEventKey(tcell.KeyRune, r, 0))
	}
}

func TestManagerExecutesFindCommand(t *testing.T) {
	registry := leader.NewRegistry()
	controller := &fakeFindController{}
	findNS := namespace.NewFindNamespace(controller)
	if _, err := namespace.Initialize(registry, findNS); err != nil {
		t.Fatalf("initialize namespace: %v", err)
	}

	status := &fakeStatusReporter{}
	manager := leader.NewManager(registry, nil, time.Second, status)

	sendRunes(manager, '\\', 'f', 't')

	if controller.tableCalls != 1 {
		t.Fatalf("expected table call, got %d", controller.tableCalls)
	}
	if len(status.errors) != 0 {
		t.Fatalf("unexpected status errors: %+v", status.errors)
	}
}

func TestManagerReportsFindErrors(t *testing.T) {
	registry := leader.NewRegistry()
	controller := &fakeFindController{tableErr: errors.New("boom")}
	findNS := namespace.NewFindNamespace(controller)
	if _, err := namespace.Initialize(registry, findNS); err != nil {
		t.Fatalf("initialize namespace: %v", err)
	}

	status := &fakeStatusReporter{}
	manager := leader.NewManager(registry, nil, time.Second, status)

	sendRunes(manager, '\\', 'f', 't')

	if len(status.errors) != 1 || status.errors[0] != "boom" {
		t.Fatalf("expected status error boom, got %+v", status.errors)
	}
}

func TestManagerExecutesTreeCommand(t *testing.T) {
	registry := leader.NewRegistry()
	controller := &fakeTreeController{}
	treeNS := namespace.NewTreeNamespace(controller)
	if _, err := namespace.Initialize(registry, treeNS); err != nil {
		t.Fatalf("initialize namespace: %v", err)
	}

	status := &fakeStatusReporter{}
	manager := leader.NewManager(registry, nil, time.Second, status)

	sendRunes(manager, '\\', 't', 't')

	if controller.toggleCalls != 1 {
		t.Fatalf("expected toggle call, got %d", controller.toggleCalls)
	}
	if len(status.errors) != 0 {
		t.Fatalf("unexpected status errors: %+v", status.errors)
	}
}

func TestManagerExecutesWorkspaceCommand(t *testing.T) {
	registry := leader.NewRegistry()
	controller := &fakeWorkspaceController{}
	workspaceNS := namespace.NewWorkspaceNamespace(controller)
	if _, err := namespace.Initialize(registry, workspaceNS); err != nil {
		t.Fatalf("initialize namespace: %v", err)
	}

	status := &fakeStatusReporter{}
	manager := leader.NewManager(registry, nil, time.Second, status)

	sendRunes(manager, '\\', 'w', 'c')

	if controller.connectCalls != 1 {
		t.Fatalf("expected connect call, got %d", controller.connectCalls)
	}
	if len(status.errors) != 0 {
		t.Fatalf("unexpected status errors: %+v", status.errors)
	}
}

func TestManagerReportsWorkspaceErrors(t *testing.T) {
	registry := leader.NewRegistry()
	controller := &fakeWorkspaceController{connectErr: errors.New("boom")}
	workspaceNS := namespace.NewWorkspaceNamespace(controller)
	if _, err := namespace.Initialize(registry, workspaceNS); err != nil {
		t.Fatalf("initialize namespace: %v", err)
	}

	status := &fakeStatusReporter{}
	manager := leader.NewManager(registry, nil, time.Second, status)

	sendRunes(manager, '\\', 'w', 'c')

	if len(status.errors) != 1 || status.errors[0] != "boom" {
		t.Fatalf("expected status error boom, got %+v", status.errors)
	}
}

func TestManagerExecutesExplainCommand(t *testing.T) {
	registry := leader.NewRegistry()
	controller := &fakeExplainController{}
	explainNS := namespace.NewExplainNamespace(controller)
	if _, err := namespace.Initialize(registry, explainNS); err != nil {
		t.Fatalf("initialize namespace: %v", err)
	}

	status := &fakeStatusReporter{}
	manager := leader.NewManager(registry, nil, time.Second, status)

	sendRunes(manager, '\\', 'x', 'x')

	if controller.explainCalls != 1 {
		t.Fatalf("expected explain call, got %d", controller.explainCalls)
	}
	if len(status.errors) != 0 {
		t.Fatalf("unexpected status errors: %+v", status.errors)
	}
}

func TestManagerReportsExplainErrors(t *testing.T) {
	registry := leader.NewRegistry()
	controller := &fakeExplainController{explainErr: errors.New("boom")}
	explainNS := namespace.NewExplainNamespace(controller)
	if _, err := namespace.Initialize(registry, explainNS); err != nil {
		t.Fatalf("initialize namespace: %v", err)
	}

	status := &fakeStatusReporter{}
	manager := leader.NewManager(registry, nil, time.Second, status)

	sendRunes(manager, '\\', 'x', 'x')

	if len(status.errors) != 1 || status.errors[0] != "boom" {
		t.Fatalf("expected status error boom, got %+v", status.errors)
	}
}

func TestManagerExecutesCancelCommand(t *testing.T) {
	registry := leader.NewRegistry()
	controller := &fakeCancelController{}
	cancelNS := namespace.NewCancelNamespace(controller)
	if _, err := namespace.Initialize(registry, cancelNS); err != nil {
		t.Fatalf("initialize namespace: %v", err)
	}

	status := &fakeStatusReporter{}
	manager := leader.NewManager(registry, nil, time.Second, status)

	sendRunes(manager, '\\', 'c', 'c')

	if controller.cancelCalls != 1 {
		t.Fatalf("expected cancel call, got %d", controller.cancelCalls)
	}
	if len(status.errors) != 0 {
		t.Fatalf("unexpected status errors: %+v", status.errors)
	}
}

func TestManagerExecutesClearResultsCommand(t *testing.T) {
	registry := leader.NewRegistry()
	controller := &fakeCancelController{}
	cancelNS := namespace.NewCancelNamespace(controller)
	if _, err := namespace.Initialize(registry, cancelNS); err != nil {
		t.Fatalf("initialize namespace: %v", err)
	}

	status := &fakeStatusReporter{}
	manager := leader.NewManager(registry, nil, time.Second, status)

	sendRunes(manager, '\\', 'c', 'l')

	if controller.clearCalls != 1 {
		t.Fatalf("expected clear call, got %d", controller.clearCalls)
	}
	if len(status.errors) != 0 {
		t.Fatalf("unexpected status errors: %+v", status.errors)
	}
}

func TestManagerReportsCancelErrors(t *testing.T) {
	registry := leader.NewRegistry()
	controller := &fakeCancelController{cancelErr: errors.New("boom")}
	cancelNS := namespace.NewCancelNamespace(controller)
	if _, err := namespace.Initialize(registry, cancelNS); err != nil {
		t.Fatalf("initialize namespace: %v", err)
	}

	status := &fakeStatusReporter{}
	manager := leader.NewManager(registry, nil, time.Second, status)

	sendRunes(manager, '\\', 'c', 'c')

	if len(status.errors) != 1 || status.errors[0] != "boom" {
		t.Fatalf("expected status error boom, got %+v", status.errors)
	}
}
