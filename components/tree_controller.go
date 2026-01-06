package components

import (
	"errors"
	"strings"

	"github.com/rivo/tview"
)

func (home *Home) ToggleTree() error {
	if home == nil || home.Tree == nil {
		return errors.New("tree not configured")
	}
	home.toggleLeftWrapper()
	home.treePinned = home.leftWrapperVisible
	return nil
}

func (home *Home) ExpandNode() error {
	tree, err := home.tree()
	if err != nil {
		return err
	}
	node := tree.GetCurrentNode()
	if node == nil {
		return nil
	}
	node.Expand()
	return nil
}

func (home *Home) CollapseNode() error {
	tree, err := home.tree()
	if err != nil {
		return err
	}
	node := tree.GetCurrentNode()
	if node == nil {
		return nil
	}
	node.Collapse()
	return nil
}

func (home *Home) ExpandAll() error {
	tree, err := home.tree()
	if err != nil {
		return err
	}
	tree.ExpandAll()
	return nil
}

func (home *Home) CollapseAll() error {
	tree, err := home.tree()
	if err != nil {
		return err
	}
	tree.CollapseAll()
	return nil
}

func (home *Home) FocusTree() error {
	if home == nil || home.Tree == nil {
		return errors.New("tree not configured")
	}
	if !home.leftWrapperVisible {
		home.toggleLeftWrapper()
		return nil
	}
	home.focusLeftWrapper()
	return nil
}

func (home *Home) RefreshTree() error {
	tree, err := home.tree()
	if err != nil {
		return err
	}
	if tree.DBDriver == nil {
		return errors.New("database driver not configured")
	}
	tree.Refresh(home.ConnectionDBName)
	return nil
}

func (home *Home) SyncWithBuffer() error {
	tree, err := home.tree()
	if err != nil {
		return err
	}
	if home.TabbedPane == nil {
		return errors.New("tabbed pane not configured")
	}
	tab := home.TabbedPane.GetCurrentTab()
	if tab == nil {
		return nil
	}
	table, ok := tab.Content.(*ResultsTable)
	if !ok {
		return nil
	}
	database := table.GetDatabaseName()
	tableName := table.GetTableName()
	if database == "" || tableName == "" {
		return nil
	}

	root := tree.GetRoot()
	if root == nil {
		return errors.New("tree not configured")
	}
	if len(root.GetChildren()) == 0 {
		tree.InitializeNodes(home.ConnectionDBName)
	}

	schema := ""
	name := tableName
	if tree.DBDriver != nil && tree.DBDriver.UseSchemas() {
		parts := strings.SplitN(tableName, ".", 2)
		if len(parts) == 2 {
			schema = parts[0]
			name = parts[1]
		}
	}

	targetDatabase := sanitizeDBName(database)
	var target *tview.TreeNode
	parents := map[*tview.TreeNode]*tview.TreeNode{}
	root.Walk(func(node, parent *tview.TreeNode) bool {
		if parent != nil {
			parents[node] = parent
		}
		if node.GetReference() == nil {
			return true
		}
		data := tree.GetTreeNodeData(node)
		if data.Type != NodeTypeTable {
			return true
		}
		if data.Database != targetDatabase {
			return true
		}
		if schema != "" && data.Schema != schema {
			return true
		}
		if data.Name != name {
			return true
		}
		target = node
		return false
	})

	if target == nil {
		return nil
	}

	for node := parents[target]; node != nil; node = parents[node] {
		node.Expand()
	}

	tree.SetCurrentNode(target)
	tree.state.selectedDatabase = targetDatabase
	if schema != "" {
		tree.state.selectedTable = schema + "." + name
	} else {
		tree.state.selectedTable = name
	}
	return nil
}

func (home *Home) tree() (*Tree, error) {
	if home == nil || home.Tree == nil {
		return nil, errors.New("tree not configured")
	}
	return home.Tree, nil
}
