package leader_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/lancekrogers/lazysql/internal/vim/buffer"
	"github.com/lancekrogers/lazysql/internal/vim/leader"
	"github.com/lancekrogers/lazysql/internal/vim/namespace"
	"github.com/lancekrogers/lazysql/internal/vim/whichkey"
)

type stubBufferListHandler struct{}

type stubShellController struct{}

type stubRunController struct{}

type stubDescribeController struct{}

type stubFindController struct{}

type stubTreeController struct{}

type stubWorkspaceController struct{}

type stubExplainController struct{}

type stubCancelController struct{}

func (s *stubBufferListHandler) Show([]*buffer.Buffer) error { return nil }

func (s *stubShellController) ToggleShell() error      { return nil }
func (s *stubShellController) CloseShell() error       { return nil }
func (s *stubShellController) FocusShellInput() error  { return nil }
func (s *stubShellController) ClearShellOutput() error { return nil }
func (s *stubShellController) ResizeShell() error      { return nil }
func (s *stubShellController) ShowShellHistory() error { return nil }

func (s *stubRunController) RunCurrentQuery() error  { return nil }
func (s *stubRunController) RunAllQueries() error    { return nil }
func (s *stubRunController) RunSelectedLines() error { return nil }
func (s *stubRunController) RunWithParams() error    { return nil }
func (s *stubRunController) ShowQueryHistory() error { return nil }

func (s *stubDescribeController) DescribeTable() error       { return nil }
func (s *stubDescribeController) DescribeView() error        { return nil }
func (s *stubDescribeController) DescribeIndexes() error     { return nil }
func (s *stubDescribeController) DescribeSequences() error   { return nil }
func (s *stubDescribeController) DescribeFunctions() error   { return nil }
func (s *stubDescribeController) DescribeDatabase() error    { return nil }
func (s *stubDescribeController) DescribeUnderCursor() error { return nil }

func (s *stubFindController) FindTable() error     { return nil }
func (s *stubFindController) FindColumn() error    { return nil }
func (s *stubFindController) FindFunction() error  { return nil }
func (s *stubFindController) FindView() error      { return nil }
func (s *stubFindController) FindSchema() error    { return nil }
func (s *stubFindController) FindInQueries() error { return nil }
func (s *stubFindController) FindRecent() error    { return nil }

func (s *stubTreeController) ToggleTree() error     { return nil }
func (s *stubTreeController) ExpandNode() error     { return nil }
func (s *stubTreeController) CollapseNode() error   { return nil }
func (s *stubTreeController) ExpandAll() error      { return nil }
func (s *stubTreeController) CollapseAll() error    { return nil }
func (s *stubTreeController) FocusTree() error      { return nil }
func (s *stubTreeController) RefreshTree() error    { return nil }
func (s *stubTreeController) SyncWithBuffer() error { return nil }

func (s *stubWorkspaceController) Connect() error          { return nil }
func (s *stubWorkspaceController) SwitchConnection() error { return nil }
func (s *stubWorkspaceController) ListConnections() error  { return nil }
func (s *stubWorkspaceController) Disconnect() error       { return nil }

func (s *stubExplainController) ExplainCurrentQuery() error        { return nil }
func (s *stubExplainController) ExplainAnalyzeCurrentQuery() error { return nil }

func (s *stubCancelController) CancelQuery() error  { return nil }
func (s *stubCancelController) ClearResults() error { return nil }

func TestLeaderNamespaceKeybindingCoverage(t *testing.T) {
	registry := buildLeaderRegistry(t)
	if registry == nil {
		t.Fatal("expected registry")
	}
	tree := registry.Tree()
	if tree == nil || tree.Root == nil {
		t.Fatal("expected leader key tree")
	}

	expectedGroups := map[rune]string{
		'b': "buffer",
		's': "shell",
		'r': "run",
		'd': "describe",
		'f': "find",
		't': "tree",
		'w': "workspace",
		'x': "explain",
		'c': "cancel",
	}

	for prefix, name := range expectedGroups {
		node := tree.Root.Children[prefix]
		if node == nil {
			t.Fatalf("expected group for %q", prefix)
		}
		if node.Description != name {
			t.Fatalf("expected group %q description %q, got %q", prefix, name, node.Description)
		}
	}

	expected := expectedNamespaceBindings()
	registered := map[string]*whichkey.KeyNode{}
	collectLeaderCommands(tree.Root, nil, registered)

	if len(registered) != len(expected) {
		t.Fatalf("expected %d bindings, got %d", len(expected), len(registered))
	}

	for sequence, expectedDesc := range expected {
		cmd, ok := registry.Lookup([]rune(sequence))
		if !ok {
			t.Fatalf("missing leader binding for %q", formatSequence(sequence))
		}
		if cmd.Handler == nil {
			t.Fatalf("missing handler for %q", formatSequence(sequence))
		}
		if cmd.Description != expectedDesc {
			t.Fatalf("expected description %q for %q, got %q", expectedDesc, formatSequence(sequence), cmd.Description)
		}

		node := nodeForSequence(tree.Root, sequence)
		if node == nil {
			t.Fatalf("missing which-key node for %q", formatSequence(sequence))
		}
		if node.Action == nil {
			t.Fatalf("missing which-key action for %q", formatSequence(sequence))
		}
		if node.Description != expectedDesc {
			t.Fatalf("expected which-key description %q for %q, got %q", expectedDesc, formatSequence(sequence), node.Description)
		}
	}

	for sequence := range registered {
		if _, ok := expected[sequence]; !ok {
			t.Fatalf("unexpected leader binding registered: %q", formatSequence(sequence))
		}
	}
}

func buildLeaderRegistry(t *testing.T) *leader.Registry {
	t.Helper()
	registry := leader.NewRegistry()

	bufferManager := buffer.NewManager()
	bufferNavigator := buffer.NewNavigator(bufferManager)
	bufferNamespace := namespace.NewBufferNamespace(bufferManager, bufferNavigator, &stubBufferListHandler{})
	shellNamespace := namespace.NewShellNamespace(&stubShellController{})
	runNamespace := namespace.NewRunNamespace(&stubRunController{})
	describeNamespace := namespace.NewDescribeNamespace(&stubDescribeController{})
	findNamespace := namespace.NewFindNamespace(&stubFindController{})
	treeNamespace := namespace.NewTreeNamespace(&stubTreeController{})
	workspaceNamespace := namespace.NewWorkspaceNamespace(&stubWorkspaceController{})
	explainNamespace := namespace.NewExplainNamespace(&stubExplainController{})
	cancelNamespace := namespace.NewCancelNamespace(&stubCancelController{})

	if _, err := namespace.Initialize(
		registry,
		bufferNamespace,
		shellNamespace,
		runNamespace,
		describeNamespace,
		findNamespace,
		treeNamespace,
		workspaceNamespace,
		explainNamespace,
		cancelNamespace,
	); err != nil {
		t.Fatalf("initialize namespaces: %v", err)
	}

	return registry
}

func expectedNamespaceBindings() map[string]string {
	expected := map[string]string{
		"bn":  "Next buffer",
		"bp":  "Previous buffer",
		"bb":  "Alternate buffer",
		"bd":  "Delete buffer",
		"bD":  "Force delete buffer",
		"bl":  "List buffers",
		"bN":  "New buffer",
		"ss":  "Toggle shell pane",
		"sq":  "Close shell pane",
		"sf":  "Focus shell input",
		"sc":  "Clear shell output",
		"sr":  "Resize shell pane",
		"sh":  "Shell history",
		"rr":  "Run current query",
		"ra":  "Run all queries",
		"rl":  "Run selected lines",
		"rp":  "Run with parameters",
		"rh":  "Query history",
		"dt":  "Describe table",
		"dv":  "Describe view",
		"di":  "Describe indexes",
		"ds":  "Describe sequences",
		"df":  "Describe functions",
		"dd":  "Describe database",
		"d\n": "Describe under cursor",
		"ft":  "Find table",
		"fc":  "Find column",
		"ff":  "Find function",
		"fv":  "Find view",
		"fs":  "Find schema",
		"fq":  "Find in queries",
		"fr":  "Find recent",
		"tt":  "Toggle tree",
		"te":  "Expand node",
		"tc":  "Collapse node",
		"tE":  "Expand all",
		"tC":  "Collapse all",
		"tf":  "Focus tree",
		"tr":  "Refresh tree",
		"ts":  "Sync with buffer",
		"wc":  "Connect",
		"ws":  "Switch connection",
		"wl":  "List connections",
		"wd":  "Disconnect",
		"xx":  "Explain query",
		"xX":  "Explain analyze",
		"cc":  "Cancel running query",
		"cl":  "Clear results",
	}

	for i := 1; i <= 9; i++ {
		key := fmt.Sprintf("b%d", i)
		expected[key] = fmt.Sprintf("Jump to buffer %d", i)
	}

	return expected
}

func collectLeaderCommands(node *whichkey.KeyNode, prefix []rune, out map[string]*whichkey.KeyNode) {
	if node == nil {
		return
	}

	if node.Action != nil {
		out[string(prefix)] = node
	}

	for key, child := range node.Children {
		next := append(append([]rune(nil), prefix...), key)
		collectLeaderCommands(child, next, out)
	}
}

func nodeForSequence(root *whichkey.KeyNode, sequence string) *whichkey.KeyNode {
	node := root
	for _, key := range sequence {
		if node == nil {
			return nil
		}
		node = node.Children[key]
	}
	return node
}

func formatSequence(sequence string) string {
	if strings.Contains(sequence, "\n") {
		sequence = strings.ReplaceAll(sequence, "\n", "<Enter>")
	}
	return sequence
}
