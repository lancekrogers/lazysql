package tree

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Navigator struct {
	pendingGotoTop bool
}

func NewNavigator() *Navigator {
	return &Navigator{}
}

func (n *Navigator) HandleVimJump(event *tcell.EventKey, view *tview.TreeView) bool {
	if event == nil || view == nil {
		if n != nil {
			n.pendingGotoTop = false
		}
		return false
	}
	if event.Key() != tcell.KeyRune {
		n.pendingGotoTop = false
		return false
	}

	switch event.Rune() {
	case 'g':
		if n.pendingGotoTop {
			n.pendingGotoTop = false
			GotoTop(view)
		} else {
			n.pendingGotoTop = true
		}
		return true
	case 'G':
		n.pendingGotoTop = false
		GotoBottom(view)
		return true
	default:
		n.pendingGotoTop = false
		return false
	}
}

func GotoTop(view *tview.TreeView) {
	if view == nil {
		return
	}
	root := view.GetRoot()
	if root == nil {
		return
	}
	view.SetCurrentNode(root)
}

func GotoBottom(view *tview.TreeView) {
	if view == nil {
		return
	}
	root := view.GetRoot()
	if root == nil {
		return
	}
	children := root.GetChildren()
	if len(children) == 0 {
		return
	}
	lastNode := children[len(children)-1]
	if lastNode.IsExpanded() {
		childNodes := lastNode.GetChildren()
		if len(childNodes) > 0 {
			view.SetCurrentNode(childNodes[len(childNodes)-1])
			return
		}
	}
	view.SetCurrentNode(lastNode)
}
