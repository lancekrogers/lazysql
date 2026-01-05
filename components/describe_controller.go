package components

import (
	"errors"
	"fmt"
	"sort"

	"github.com/rivo/tview"

	"github.com/lancekrogers/lazysql/app"
)

func (home *Home) DescribeTable() error {
	database, table := home.currentTableContext()
	if table == "" {
		return home.describeTablePicker()
	}
	return home.describeTableByName(database, table)
}

func (home *Home) DescribeView() error {
	selection, err := home.treeSelection()
	if err == nil && selection.Type == NodeTypeView {
		return home.describeViewByName(selection.Database, selection.qualifiedName())
	}
	return home.describeViewPicker()
}

func (home *Home) DescribeIndexes() error {
	database, table := home.currentTableContext()
	if table == "" {
		return home.describeIndexesPicker()
	}
	return home.describeIndexesByName(database, table)
}

func (home *Home) DescribeSequences() error {
	return errors.New("sequences are not supported by this driver")
}

func (home *Home) DescribeFunctions() error {
	selection, err := home.treeSelection()
	if err == nil && selection.Type == NodeTypeFunction {
		return home.describeFunctionByName(selection.Database, selection.qualifiedName())
	}
	return home.describeFunctionPicker()
}

func (home *Home) DescribeDatabase() error {
	if home.DBDriver == nil {
		return errors.New("database driver not configured")
	}
	databases, err := home.DBDriver.GetDatabases()
	if err != nil {
		return err
	}
	if len(databases) == 0 {
		return errors.New("no databases found")
	}
	sort.Strings(databases)
	return home.showDescribePicker("Databases", databases, func(selected string) {
		home.showStatusInfo(fmt.Sprintf("Database: %s", selected))
	})
}

func (home *Home) DescribeUnderCursor() error {
	selection, err := home.treeSelection()
	if err != nil {
		return err
	}
	switch selection.Type {
	case NodeTypeDatabase:
		return home.DescribeDatabase()
	case NodeTypeTable:
		return home.describeTableByName(selection.Database, selection.qualifiedName())
	case NodeTypeView:
		return home.describeViewByName(selection.Database, selection.qualifiedName())
	case NodeTypeFunction:
		return home.describeFunctionByName(selection.Database, selection.qualifiedName())
	case NodeTypeProcedure:
		return home.describeProcedureByName(selection.Database, selection.qualifiedName())
	default:
		return errors.New("no describable object under cursor")
	}
}

func (home *Home) currentTableContext() (string, string) {
	if home != nil && home.TabbedPane != nil {
		if tab := home.TabbedPane.GetCurrentTab(); tab != nil {
			if table, ok := tab.Content.(*ResultsTable); ok {
				if tableName := table.GetTableName(); tableName != "" {
					return table.GetDatabaseName(), tableName
				}
			}
		}
	}
	if home != nil && home.Tree != nil {
		return home.Tree.GetSelectedDatabase(), home.Tree.GetSelectedTable()
	}
	return "", ""
}

func (home *Home) describeTableByName(database string, table string) error {
	if home == nil {
		return errors.New("home not configured")
	}
	if database == "" || table == "" {
		return errors.New("table selection required")
	}
	home.showTable(database, table)
	tab := home.TabbedPane.GetCurrentTab()
	if tab == nil {
		return errors.New("table view not available")
	}
	results, ok := tab.Content.(*ResultsTable)
	if !ok || results == nil {
		return errors.New("table view not available")
	}
	results.ShowColumns()
	return nil
}

func (home *Home) describeIndexesByName(database string, table string) error {
	if home == nil {
		return errors.New("home not configured")
	}
	if database == "" || table == "" {
		return errors.New("table selection required")
	}
	home.showTable(database, table)
	tab := home.TabbedPane.GetCurrentTab()
	if tab == nil {
		return errors.New("table view not available")
	}
	results, ok := tab.Content.(*ResultsTable)
	if !ok || results == nil {
		return errors.New("table view not available")
	}
	results.ShowIndexes()
	return nil
}

func (home *Home) describeViewByName(database string, view string) error {
	if home.DBDriver == nil {
		return errors.New("database driver not configured")
	}
	if view == "" {
		return errors.New("view name required")
	}
	definition, err := home.DBDriver.GetViewDefinition(database, view)
	if err != nil {
		return err
	}
	return home.showDefinition(definition)
}

func (home *Home) describeFunctionByName(database string, name string) error {
	if home.DBDriver == nil {
		return errors.New("database driver not configured")
	}
	if name == "" {
		return errors.New("function name required")
	}
	definition, err := home.DBDriver.GetFunctionDefinition(database, name)
	if err != nil {
		return err
	}
	return home.showDefinition(definition)
}

func (home *Home) describeProcedureByName(database string, name string) error {
	if home.DBDriver == nil {
		return errors.New("database driver not configured")
	}
	if name == "" {
		return errors.New("procedure name required")
	}
	definition, err := home.DBDriver.GetProcedureDefinition(database, name)
	if err != nil {
		return err
	}
	return home.showDefinition(definition)
}

func (home *Home) showDefinition(definition string) error {
	home.createOrFocusEditorTab()
	tab := home.TabbedPane.GetCurrentTab()
	if tab == nil {
		return errors.New("editor tab not available")
	}
	results, ok := tab.Content.(*ResultsTable)
	if !ok || results.Editor == nil {
		return errors.New("editor not available")
	}
	results.Editor.SetText(definition, false)
	App.ForceDraw()
	return nil
}

func (home *Home) treeSelection() (*TreeNodeData, error) {
	if home == nil || home.Tree == nil {
		return nil, errors.New("tree not configured")
	}
	node := home.Tree.GetCurrentNode()
	if node == nil {
		return nil, errors.New("no tree selection")
	}
	data := home.Tree.GetTreeNodeData(node)
	if data == nil || data.Type == NodeTypeSection {
		return nil, errors.New("no describable object selected")
	}
	return data, nil
}

func (home *Home) describeTablePicker() error {
	if home.DBDriver == nil {
		return errors.New("database driver not configured")
	}
	database := ""
	if home.Tree != nil {
		database = home.Tree.GetSelectedDatabase()
	}
	if database == "" {
		return errors.New("no database selected")
	}
	tables, err := home.DBDriver.GetTables(database)
	if err != nil {
		return err
	}
	items := flattenQualifiedNames(tables, database)
	if len(items) == 0 {
		return errors.New("no tables found")
	}
	return home.showDescribePicker("Tables", items, func(selected string) {
		if err := home.describeTableByName(database, selected); err != nil {
			home.showStatusError(err.Error())
		}
	})
}

func (home *Home) describeIndexesPicker() error {
	if home.DBDriver == nil {
		return errors.New("database driver not configured")
	}
	database := ""
	if home.Tree != nil {
		database = home.Tree.GetSelectedDatabase()
	}
	if database == "" {
		return errors.New("no database selected")
	}
	tables, err := home.DBDriver.GetTables(database)
	if err != nil {
		return err
	}
	items := flattenQualifiedNames(tables, database)
	if len(items) == 0 {
		return errors.New("no tables found")
	}
	return home.showDescribePicker("Tables", items, func(selected string) {
		if err := home.describeIndexesByName(database, selected); err != nil {
			home.showStatusError(err.Error())
		}
	})
}

func (home *Home) describeViewPicker() error {
	if home.DBDriver == nil {
		return errors.New("database driver not configured")
	}
	database := ""
	if home.Tree != nil {
		database = home.Tree.GetSelectedDatabase()
	}
	if database == "" {
		return errors.New("no database selected")
	}
	views, err := home.DBDriver.GetViews(database)
	if err != nil {
		return err
	}
	items := flattenQualifiedNames(views, database)
	if len(items) == 0 {
		return errors.New("no views found")
	}
	return home.showDescribePicker("Views", items, func(selected string) {
		if err := home.describeViewByName(database, selected); err != nil {
			home.showStatusError(err.Error())
		}
	})
}

func (home *Home) describeFunctionPicker() error {
	if home.DBDriver == nil {
		return errors.New("database driver not configured")
	}
	database := ""
	if home.Tree != nil {
		database = home.Tree.GetSelectedDatabase()
	}
	if database == "" {
		return errors.New("no database selected")
	}
	functions, err := home.DBDriver.GetFunctions(database)
	if err != nil {
		return err
	}
	items := flattenQualifiedNames(functions, database)
	if len(items) == 0 {
		return errors.New("no functions found")
	}
	return home.showDescribePicker("Functions", items, func(selected string) {
		if err := home.describeFunctionByName(database, selected); err != nil {
			home.showStatusError(err.Error())
		}
	})
}

func (home *Home) showDescribePicker(title string, items []string, onSelect func(string)) error {
	if mainPages == nil {
		return errors.New("main pages not configured")
	}
	list := tview.NewList()
	list.SetBorder(true)
	list.SetTitle(fmt.Sprintf(" %s ", title))
	list.SetMainTextColor(app.Styles.PrimaryTextColor)
	list.SetSecondaryTextColor(app.Styles.InverseTextColor)
	list.ShowSecondaryText(false)

	previousFocus := app.App.GetFocus()
	for _, item := range items {
		entry := item
		list.AddItem(entry, "", 0, func() {
			mainPages.RemovePage(pageNameDescribePicker)
			if onSelect != nil {
				onSelect(entry)
			}
			if previousFocus != nil {
				app.App.SetFocus(previousFocus)
			}
		})
	}

	list.SetDoneFunc(func() {
		mainPages.RemovePage(pageNameDescribePicker)
		if previousFocus != nil {
			app.App.SetFocus(previousFocus)
		}
	})

	if mainPages.HasPage(pageNameDescribePicker) {
		mainPages.RemovePage(pageNameDescribePicker)
	}
	mainPages.AddPage(pageNameDescribePicker, list, true, true)
	app.App.SetFocus(list)
	return nil
}

func flattenQualifiedNames(grouped map[string][]string, skipPrefix string) []string {
	items := make([]string, 0)
	for group, names := range grouped {
		for _, name := range names {
			item := name
			if group != "" && group != skipPrefix {
				item = fmt.Sprintf("%s.%s", group, name)
			}
			items = append(items, item)
		}
	}
	sort.Strings(items)
	return items
}

func (data *TreeNodeData) qualifiedName() string {
	if data == nil {
		return ""
	}
	if data.Schema == "" {
		return data.Name
	}
	return fmt.Sprintf("%s.%s", data.Schema, data.Name)
}
