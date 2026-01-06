package namespace

import (
	"errors"

	"github.com/lancekrogers/lazysql/internal/vim/leader"
)

type TreeController interface {
	ToggleTree() error
	ExpandNode() error
	CollapseNode() error
	ExpandAll() error
	CollapseAll() error
	FocusTree() error
	RefreshTree() error
	SyncWithBuffer() error
}

type TreeNamespace struct {
	controller TreeController
}

func NewTreeNamespace(controller TreeController) *TreeNamespace {
	return &TreeNamespace{controller: controller}
}

func (t *TreeNamespace) Prefix() rune {
	return 't'
}

func (t *TreeNamespace) Name() string {
	return "tree"
}

func (t *TreeNamespace) Commands() []leader.Command {
	return []leader.Command{
		{Sequence: []rune{'t'}, Description: "Toggle tree", Handler: t.toggleTree},
		{Sequence: []rune{'e'}, Description: "Expand node", Handler: t.expandNode},
		{Sequence: []rune{'c'}, Description: "Collapse node", Handler: t.collapseNode},
		{Sequence: []rune{'E'}, Description: "Expand all", Handler: t.expandAll},
		{Sequence: []rune{'C'}, Description: "Collapse all", Handler: t.collapseAll},
		{Sequence: []rune{'f'}, Description: "Focus tree", Handler: t.focusTree},
		{Sequence: []rune{'r'}, Description: "Refresh tree", Handler: t.refreshTree},
		{Sequence: []rune{'s'}, Description: "Sync with buffer", Handler: t.syncWithBuffer},
	}
}

func (t *TreeNamespace) toggleTree() error {
	if t.controller == nil {
		return errors.New("tree controller not configured")
	}
	return t.controller.ToggleTree()
}

func (t *TreeNamespace) expandNode() error {
	if t.controller == nil {
		return errors.New("tree controller not configured")
	}
	return t.controller.ExpandNode()
}

func (t *TreeNamespace) collapseNode() error {
	if t.controller == nil {
		return errors.New("tree controller not configured")
	}
	return t.controller.CollapseNode()
}

func (t *TreeNamespace) expandAll() error {
	if t.controller == nil {
		return errors.New("tree controller not configured")
	}
	return t.controller.ExpandAll()
}

func (t *TreeNamespace) collapseAll() error {
	if t.controller == nil {
		return errors.New("tree controller not configured")
	}
	return t.controller.CollapseAll()
}

func (t *TreeNamespace) focusTree() error {
	if t.controller == nil {
		return errors.New("tree controller not configured")
	}
	return t.controller.FocusTree()
}

func (t *TreeNamespace) refreshTree() error {
	if t.controller == nil {
		return errors.New("tree controller not configured")
	}
	return t.controller.RefreshTree()
}

func (t *TreeNamespace) syncWithBuffer() error {
	if t.controller == nil {
		return errors.New("tree controller not configured")
	}
	return t.controller.SyncWithBuffer()
}
