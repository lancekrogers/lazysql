package tree

import (
	"testing"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func buildTestTree(expandLast bool) (*tview.TreeView, *tview.TreeNode, *tview.TreeNode, *tview.TreeNode) {
	root := tview.NewTreeNode("root")
	first := tview.NewTreeNode("first")
	last := tview.NewTreeNode("last")
	lastChild := tview.NewTreeNode("last-child")
	if expandLast {
		last.SetExpanded(true)
		last.AddChild(lastChild)
	}
	root.AddChild(first)
	root.AddChild(last)

	view := tview.NewTreeView()
	view.SetRoot(root)
	view.SetCurrentNode(first)

	return view, root, last, lastChild
}

func TestGotoTop(t *testing.T) {
	view, root, _, _ := buildTestTree(false)
	GotoTop(view)
	if view.GetCurrentNode() != root {
		t.Fatal("expected current node to be root after goto top")
	}
}

func TestGotoBottom(t *testing.T) {
	var last *tview.TreeNode
	view, _, _, lastChild := buildTestTree(true)
	GotoBottom(view)
	if view.GetCurrentNode() != lastChild {
		t.Fatal("expected current node to be last expanded child")
	}

	view, _, last, _ = buildTestTree(false)
	GotoBottom(view)
	if view.GetCurrentNode() != last {
		t.Fatal("expected current node to be last child when not expanded")
	}
}

func TestNavigatorHandleVimJump(t *testing.T) {
	view, root, last, _ := buildTestTree(false)
	navigator := NewNavigator()

	event := tcell.NewEventKey(tcell.KeyRune, 'g', tcell.ModNone)
	if !navigator.HandleVimJump(event, view) {
		t.Fatal("expected first g to be handled")
	}
	if view.GetCurrentNode() == root {
		t.Fatal("expected first g to set pending only")
	}

	if !navigator.HandleVimJump(event, view) {
		t.Fatal("expected second g to be handled")
	}
	if view.GetCurrentNode() != root {
		t.Fatal("expected gg to move to root")
	}

	event = tcell.NewEventKey(tcell.KeyRune, 'G', tcell.ModNone)
	if !navigator.HandleVimJump(event, view) {
		t.Fatal("expected G to be handled")
	}
	if view.GetCurrentNode() != last {
		t.Fatal("expected G to move to bottom")
	}
}

func TestNavigatorHandleVimJumpNonRune(t *testing.T) {
	view, _, _, _ := buildTestTree(false)
	navigator := NewNavigator()
	navigator.pendingGotoTop = true

	event := tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone)
	if navigator.HandleVimJump(event, view) {
		t.Fatal("expected non-rune key to be ignored")
	}
	if navigator.pendingGotoTop {
		t.Fatal("expected pending state to reset on non-rune key")
	}
}

func TestNavigatorHandleVimJumpNilView(t *testing.T) {
	navigator := NewNavigator()
	event := tcell.NewEventKey(tcell.KeyRune, 'g', tcell.ModNone)
	if navigator.HandleVimJump(event, nil) {
		t.Fatal("expected nil view to be ignored")
	}
	if navigator.pendingGotoTop {
		t.Fatal("expected pending state to reset on nil view")
	}
}

func TestGotoTopBottomNilView(t *testing.T) {
	t.Helper()
	GotoTop(nil)
	GotoBottom(nil)
}
