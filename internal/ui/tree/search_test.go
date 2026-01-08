package tree

import (
	"testing"

	"github.com/rivo/tview"
)

func buildSearchTree() (*tview.TreeNode, *tview.TreeNode, *tview.TreeNode) {
	root := tview.NewTreeNode("root")
	db1 := tview.NewTreeNode("db1")
	db2 := tview.NewTreeNode("db2")

	users := tview.NewTreeNode("users")
	ordersDB1 := tview.NewTreeNode("orders")
	ordersDB2 := tview.NewTreeNode("orders")

	db1.AddChild(users)
	db1.AddChild(ordersDB1)
	db2.AddChild(ordersDB2)

	root.AddChild(db1)
	root.AddChild(db2)

	return root, db1, ordersDB1
}

func TestSearchFindsMatches(t *testing.T) {
	root, db1, _ := buildSearchTree()
	matches := Search(root, "orders")
	if len(matches) != 2 {
		t.Fatalf("expected 2 matches, got %d", len(matches))
	}
	if !db1.IsExpanded() {
		t.Fatal("expected parent to be expanded when matches found")
	}
}

func TestSearchWithDatabaseFilter(t *testing.T) {
	root, _, ordersDB1 := buildSearchTree()
	matches := Search(root, "db1 orders")
	if len(matches) != 1 {
		t.Fatalf("expected 1 match, got %d", len(matches))
	}
	if matches[0] != ordersDB1 {
		t.Fatalf("expected db1 orders match, got %q", matches[0].GetText())
	}
}

func TestPrioritizeResultScoring(t *testing.T) {
	if score := prioritizeResult("foo", "foo", 0); score != 0 {
		t.Fatalf("expected exact match score 0, got %d", score)
	}
	if score := prioritizeResult("foo", "foobar", 1); score != 4 {
		t.Fatalf("expected prefix score 4, got %d", score)
	}
	if score := prioritizeResult("foo", "barfoo", 1); score != 106 {
		t.Fatalf("expected contains score 106, got %d", score)
	}
	if score := prioritizeResult("foo", "bar", 7); score != 10007 {
		t.Fatalf("expected fallback score 10007, got %d", score)
	}
}
