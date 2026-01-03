package whichkey

import "testing"

func TestKeyTreeNavigation(t *testing.T) {
	tree := NewKeyTree()
	tree.AddGroup([]rune{'f'}, "find")
	tree.AddCommand([]rune{'f', 't'}, "table", func() {})

	if len(tree.Root.Children) != 1 {
		t.Fatalf("expected 1 root child, got %d", len(tree.Root.Children))
	}

	if !tree.Enter('f') {
		t.Fatal("expected to enter 'f' group")
	}
	if tree.Current.Key != 'f' {
		t.Fatalf("expected current key 'f', got %q", tree.Current.Key)
	}

	if !tree.Enter('t') {
		t.Fatal("expected to enter 't' command")
	}
	if tree.Current.Key != 't' || tree.Current.Action == nil {
		t.Fatal("expected leaf node with action")
	}

	if !tree.Back() || tree.Current.Key != 'f' {
		t.Fatal("expected to return to 'f' group")
	}
	if !tree.Back() || tree.Current != tree.Root {
		t.Fatal("expected to return to root")
	}
	if tree.Back() {
		t.Fatal("expected back from root to fail")
	}
}

func TestKeyTreeBreadcrumb(t *testing.T) {
	tree := NewKeyTree()
	tree.AddGroup([]rune{'f'}, "find")
	tree.AddCommand([]rune{'f', 't'}, "table", func() {})

	if !tree.Enter('f') {
		t.Fatal("expected to enter 'f'")
	}
	if !tree.Enter('t') {
		t.Fatal("expected to enter 't'")
	}

	breadcrumb := tree.Breadcrumb('\\')
	if breadcrumb != `\ > f > t > ` {
		t.Fatalf("unexpected breadcrumb: %q", breadcrumb)
	}
}

func TestKeyTreeReset(t *testing.T) {
	tree := NewKeyTree()
	tree.AddGroup([]rune{'f'}, "find")
	tree.AddCommand([]rune{'f', 't'}, "table", func() {})

	if !tree.Enter('f') {
		t.Fatal("expected to enter 'f'")
	}
	if tree.Current == tree.Root {
		t.Fatal("expected non-root before reset")
	}

	tree.Reset()
	if tree.Current != tree.Root {
		t.Fatal("expected reset to return to root")
	}

	leaf := tree.Root.Children['f'].Children['t']
	if leaf == nil || !leaf.IsLeaf() {
		t.Fatal("expected leaf node to report IsLeaf true")
	}
}
