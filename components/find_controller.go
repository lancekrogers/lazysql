package components

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/lancekrogers/lazysql/internal/vim/buffer"
)

type findKind string

const (
	findKindTable    findKind = "table"
	findKindColumn   findKind = "column"
	findKindFunction findKind = "function"
	findKindView     findKind = "view"
	findKindSchema   findKind = "schema"
)

type findHistoryItem struct {
	kind      findKind
	database  string
	schema    string
	name      string
	column    string
	timestamp time.Time
}

type findHistory struct {
	mu    sync.Mutex
	items []findHistoryItem
	limit int
}

func newFindHistory(limit int) *findHistory {
	if limit <= 0 {
		limit = 50
	}
	return &findHistory{limit: limit}
}

func (h *findHistory) add(item findHistoryItem) {
	if h == nil {
		return
	}
	h.mu.Lock()
	defer h.mu.Unlock()

	item.timestamp = time.Now().UTC()
	key := findHistoryKey(item)
	filtered := make([]findHistoryItem, 0, len(h.items))
	for _, existing := range h.items {
		if findHistoryKey(existing) == key {
			continue
		}
		filtered = append(filtered, existing)
	}
	h.items = append([]findHistoryItem{item}, filtered...)
	if len(h.items) > h.limit {
		h.items = h.items[:h.limit]
	}
}

func (h *findHistory) list() []findHistoryItem {
	if h == nil {
		return nil
	}
	h.mu.Lock()
	defer h.mu.Unlock()
	items := make([]findHistoryItem, len(h.items))
	copy(items, h.items)
	return items
}

func findHistoryKey(item findHistoryItem) string {
	return string(item.kind) + "|" + item.database + "|" + item.schema + "|" + item.name + "|" + item.column
}

type findObject struct {
	database string
	schema   string
	name     string
}

func (o findObject) qualifiedName() string {
	if o.schema == "" {
		return o.name
	}
	return fmt.Sprintf("%s.%s", o.schema, o.name)
}

func (home *Home) FindTable() error {
	items, err := home.findTableItems()
	if err != nil {
		return err
	}
	return home.showFindPicker("Tables", "Find table...", items)
}

func (home *Home) FindColumn() error {
	items, err := home.findColumnItems()
	if err != nil {
		return err
	}
	return home.showFindPicker("Columns", "Find column...", items)
}

func (home *Home) FindFunction() error {
	items, err := home.findFunctionItems()
	if err != nil {
		return err
	}
	return home.showFindPicker("Functions", "Find function...", items)
}

func (home *Home) FindView() error {
	items, err := home.findViewItems()
	if err != nil {
		return err
	}
	return home.showFindPicker("Views", "Find view...", items)
}

func (home *Home) FindSchema() error {
	items, err := home.findSchemaItems()
	if err != nil {
		return err
	}
	return home.showFindPicker("Schemas", "Find schema...", items)
}

func (home *Home) FindInQueries() error {
	items, err := home.findQueryItems()
	if err != nil {
		return err
	}
	return home.showFindPicker("Queries", "Find in queries...", items)
}

func (home *Home) FindRecent() error {
	items, err := home.findRecentItems()
	if err != nil {
		return err
	}
	return home.showFindPicker("Recent", "Find recent...", items)
}

func (home *Home) findTableItems() ([]FindPickerItem, error) {
	if home.DBDriver == nil {
		return nil, errors.New("database driver not configured")
	}
	objects, includeDatabase, err := home.collectObjects(home.DBDriver.GetTables)
	if err != nil {
		return nil, err
	}
	if len(objects) == 0 {
		return nil, errors.New("no tables found")
	}

	items := make([]FindPickerItem, 0, len(objects))
	for _, obj := range objects {
		obj := obj
		label, secondary := formatFindLabel(obj, includeDatabase)
		items = append(items, FindPickerItem{
			Label:     label,
			Secondary: secondary,
			OnSelect: func() {
				home.showTable(obj.database, obj.qualifiedName())
				home.addFindHistory(findHistoryItem{
					kind:     findKindTable,
					database: obj.database,
					schema:   obj.schema,
					name:     obj.name,
				})
			},
		})
	}
	return items, nil
}

func (home *Home) findColumnItems() ([]FindPickerItem, error) {
	if home.DBDriver == nil {
		return nil, errors.New("database driver not configured")
	}
	objects, includeDatabase, err := home.collectObjects(home.DBDriver.GetTables)
	if err != nil {
		return nil, err
	}
	if len(objects) == 0 {
		return nil, errors.New("no tables found")
	}

	items := make([]FindPickerItem, 0, len(objects))
	for _, obj := range objects {
		obj := obj
		columns, err := home.columnNamesForTable(obj)
		if err != nil {
			return nil, err
		}
		for _, column := range columns {
			column := column
			label, secondary := formatFindColumnLabel(obj, column, includeDatabase)
			items = append(items, FindPickerItem{
				Label:     label,
				Secondary: secondary,
				OnSelect: func() {
					if err := home.describeTableByName(obj.database, obj.qualifiedName()); err != nil {
						home.showStatusError(err.Error())
						return
					}
					home.addFindHistory(findHistoryItem{
						kind:     findKindColumn,
						database: obj.database,
						schema:   obj.schema,
						name:     obj.name,
						column:   column,
					})
				},
			})
		}
	}
	if len(items) == 0 {
		return nil, errors.New("no columns found")
	}
	return items, nil
}

func (home *Home) findFunctionItems() ([]FindPickerItem, error) {
	if home.DBDriver == nil {
		return nil, errors.New("database driver not configured")
	}
	if !home.DBDriver.SupportsProgramming() {
		return nil, errors.New("functions are not supported by this driver")
	}
	objects, includeDatabase, err := home.collectObjects(home.DBDriver.GetFunctions)
	if err != nil {
		return nil, err
	}
	if len(objects) == 0 {
		return nil, errors.New("no functions found")
	}
	return home.buildObjectItems(objects, includeDatabase, findKindFunction, func(obj findObject) error {
		return home.describeFunctionByName(obj.database, obj.qualifiedName())
	})
}

func (home *Home) findViewItems() ([]FindPickerItem, error) {
	if home.DBDriver == nil {
		return nil, errors.New("database driver not configured")
	}
	if !home.DBDriver.SupportsProgramming() {
		return nil, errors.New("views are not supported by this driver")
	}
	objects, includeDatabase, err := home.collectObjects(home.DBDriver.GetViews)
	if err != nil {
		return nil, err
	}
	if len(objects) == 0 {
		return nil, errors.New("no views found")
	}
	return home.buildObjectItems(objects, includeDatabase, findKindView, func(obj findObject) error {
		return home.describeViewByName(obj.database, obj.qualifiedName())
	})
}

func (home *Home) findSchemaItems() ([]FindPickerItem, error) {
	schemas, includeDatabase, err := home.collectSchemas()
	if err != nil {
		return nil, err
	}
	if len(schemas) == 0 {
		return nil, errors.New("no schemas found")
	}
	items := make([]FindPickerItem, 0, len(schemas))
	for _, schema := range schemas {
		schema := schema
		label, secondary := formatSchemaLabel(schema, includeDatabase)
		items = append(items, FindPickerItem{
			Label:     label,
			Secondary: secondary,
			OnSelect: func() {
				home.showStatusInfo(fmt.Sprintf("Schema: %s", schema.name))
				home.addFindHistory(findHistoryItem{
					kind:     findKindSchema,
					database: schema.database,
					name:     schema.name,
				})
			},
		})
	}
	return items, nil
}

func (home *Home) findQueryItems() ([]FindPickerItem, error) {
	if home == nil || home.BufferManager == nil {
		return nil, errors.New("buffer manager not configured")
	}

	buffers := home.BufferManager.ListBuffers()
	items := make([]FindPickerItem, 0, len(buffers))
	for _, buf := range buffers {
		if buf == nil || strings.TrimSpace(buf.Content) == "" {
			continue
		}
		items = append(items, home.buildQueryItems(buf)...)
	}
	if len(items) == 0 {
		return nil, errors.New("no query content found")
	}
	return items, nil
}

func (home *Home) findRecentItems() ([]FindPickerItem, error) {
	history := home.ensureFindHistory()
	if history == nil {
		return nil, errors.New("find history not configured")
	}
	items := history.list()
	if len(items) == 0 {
		return nil, errors.New("no recent items")
	}

	includeDatabase := hasMultipleDatabases(items)
	findItems := make([]FindPickerItem, 0, len(items))
	for _, entry := range items {
		entry := entry
		label, secondary := formatHistoryLabel(entry, includeDatabase)
		findItems = append(findItems, FindPickerItem{
			Label:     label,
			Secondary: secondary,
			OnSelect: func() {
				if err := home.runHistoryEntry(entry); err != nil {
					home.showStatusError(err.Error())
				}
			},
		})
	}
	return findItems, nil
}

func (home *Home) buildObjectItems(objects []findObject, includeDatabase bool, kind findKind, action func(findObject) error) ([]FindPickerItem, error) {
	items := make([]FindPickerItem, 0, len(objects))
	for _, obj := range objects {
		obj := obj
		label, secondary := formatFindLabel(obj, includeDatabase)
		items = append(items, FindPickerItem{
			Label:     label,
			Secondary: secondary,
			OnSelect: func() {
				if err := action(obj); err != nil {
					home.showStatusError(err.Error())
					return
				}
				home.addFindHistory(findHistoryItem{
					kind:     kind,
					database: obj.database,
					schema:   obj.schema,
					name:     obj.name,
				})
			},
		})
	}
	return items, nil
}

func (home *Home) buildQueryItems(buf *buffer.Buffer) []FindPickerItem {
	lines := strings.Split(buf.Content, "\n")
	items := make([]FindPickerItem, 0, len(lines))
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		lineNum := i + 1
		label := fmt.Sprintf("%s:%d %s", buf.Name, lineNum, trimmed)
		secondary := buf.FilePath
		bufferID := buf.ID
		bufferName := buf.Name
		items = append(items, FindPickerItem{
			Label:     label,
			Secondary: secondary,
			OnSelect: func() {
				if err := home.openQueryBuffer(bufferID, bufferName, lineNum); err != nil {
					home.showStatusError(err.Error())
				}
			},
		})
	}
	return items
}

func (home *Home) openQueryBuffer(bufferID buffer.BufferID, bufferName string, line int) error {
	if home == nil || home.BufferManager == nil {
		return errors.New("buffer manager not configured")
	}
	home.createOrFocusEditorTab()
	if err := home.BufferManager.SetActive(bufferID); err != nil {
		return err
	}
	home.focusEditor()
	home.showStatusInfo(fmt.Sprintf("Jumped to %s:%d", bufferName, line))
	return nil
}

func (home *Home) focusEditor() {
	if home == nil || home.TabbedPane == nil {
		return
	}
	tab := home.TabbedPane.GetCurrentTab()
	if tab == nil {
		return
	}
	results, ok := tab.Content.(*ResultsTable)
	if !ok || results.Editor == nil {
		return
	}
	App.SetFocus(results.Editor)
}

func (home *Home) showFindPicker(title string, placeholder string, items []FindPickerItem) error {
	if home == nil || home.FindPicker == nil {
		return errors.New("find picker not configured")
	}
	return home.FindPicker.Show(FindPickerConfig{
		Title:       title,
		Placeholder: placeholder,
		Items:       items,
	})
}

func (home *Home) collectObjects(fetch func(string) (map[string][]string, error)) ([]findObject, bool, error) {
	databases, includeDatabase, err := home.listDatabases()
	if err != nil {
		return nil, includeDatabase, err
	}

	useSchemas := home.DBDriver.UseSchemas()
	objects := make([]findObject, 0)
	for _, database := range databases {
		entries, err := fetch(database)
		if err != nil {
			return nil, includeDatabase, err
		}
		for schema, names := range entries {
			sort.Strings(names)
			for _, name := range names {
				obj := findObject{database: database, name: name}
				if useSchemas {
					obj.schema = schema
				}
				objects = append(objects, obj)
			}
		}
	}

	sort.SliceStable(objects, func(i, j int) bool {
		left := objects[i].qualifiedName()
		right := objects[j].qualifiedName()
		if includeDatabase && objects[i].database != objects[j].database {
			return objects[i].database < objects[j].database
		}
		return left < right
	})

	return objects, includeDatabase, nil
}

type findSchema struct {
	database string
	name     string
}

func (home *Home) collectSchemas() ([]findSchema, bool, error) {
	if home.DBDriver == nil {
		return nil, false, errors.New("database driver not configured")
	}
	if !home.DBDriver.UseSchemas() {
		return nil, false, errors.New("schemas are not supported by this driver")
	}
	databases, includeDatabase, err := home.listDatabases()
	if err != nil {
		return nil, includeDatabase, err
	}

	seen := make(map[string]struct{})
	schemas := make([]findSchema, 0)
	for _, database := range databases {
		entries, err := home.DBDriver.GetTables(database)
		if err != nil {
			return nil, includeDatabase, err
		}
		for schema := range entries {
			key := database + "|" + schema
			if _, ok := seen[key]; ok {
				continue
			}
			seen[key] = struct{}{}
			schemas = append(schemas, findSchema{database: database, name: schema})
		}
	}

	sort.SliceStable(schemas, func(i, j int) bool {
		if includeDatabase && schemas[i].database != schemas[j].database {
			return schemas[i].database < schemas[j].database
		}
		return schemas[i].name < schemas[j].name
	})

	return schemas, includeDatabase, nil
}

func (home *Home) listDatabases() ([]string, bool, error) {
	if home == nil || home.DBDriver == nil {
		return nil, false, errors.New("database driver not configured")
	}
	databases, err := home.DBDriver.GetDatabases()
	if err != nil {
		return nil, false, err
	}
	if len(databases) == 0 {
		return nil, false, errors.New("no databases found")
	}
	sort.Strings(databases)
	return databases, len(databases) > 1, nil
}

func (home *Home) columnNamesForTable(obj findObject) ([]string, error) {
	if home.DBDriver == nil {
		return nil, errors.New("database driver not configured")
	}
	columns, err := home.DBDriver.GetTableColumns(obj.database, obj.qualifiedName())
	if err != nil {
		return nil, err
	}
	if len(columns) <= 1 {
		return nil, nil
	}

	names := make([]string, 0, len(columns)-1)
	for i := 1; i < len(columns); i++ {
		row := columns[i]
		if len(row) == 0 {
			continue
		}
		name := row[0]
		if name != "" {
			names = append(names, name)
		}
	}
	sort.Strings(names)
	return names, nil
}

func formatFindLabel(obj findObject, includeDatabase bool) (string, string) {
	label := obj.qualifiedName()
	if includeDatabase {
		return label, obj.database
	}
	return label, ""
}

func formatFindColumnLabel(obj findObject, column string, includeDatabase bool) (string, string) {
	label := fmt.Sprintf("%s.%s", obj.qualifiedName(), column)
	if includeDatabase {
		return label, obj.database
	}
	return label, ""
}

func formatSchemaLabel(schema findSchema, includeDatabase bool) (string, string) {
	if includeDatabase {
		return schema.name, schema.database
	}
	return schema.name, ""
}

func (home *Home) ensureFindHistory() *findHistory {
	if home == nil {
		return nil
	}
	if home.findHistory == nil {
		home.findHistory = newFindHistory(50)
	}
	return home.findHistory
}

func (home *Home) addFindHistory(item findHistoryItem) {
	history := home.ensureFindHistory()
	if history == nil {
		return
	}
	history.add(item)
}

func hasMultipleDatabases(items []findHistoryItem) bool {
	seen := make(map[string]struct{})
	for _, item := range items {
		if item.database == "" {
			continue
		}
		seen[item.database] = struct{}{}
		if len(seen) > 1 {
			return true
		}
	}
	return false
}

func formatHistoryLabel(item findHistoryItem, includeDatabase bool) (string, string) {
	obj := findObject{
		database: item.database,
		schema:   item.schema,
		name:     item.name,
	}
	label := ""
	switch item.kind {
	case findKindColumn:
		label, _ = formatFindColumnLabel(obj, item.column, false)
	default:
		label, _ = formatFindLabel(obj, false)
	}
	return fmt.Sprintf("%s: %s", formatHistoryKind(item.kind), label), formatHistorySecondary(item, includeDatabase)
}

func formatHistoryKind(kind findKind) string {
	switch kind {
	case findKindTable:
		return "table"
	case findKindColumn:
		return "column"
	case findKindFunction:
		return "function"
	case findKindView:
		return "view"
	case findKindSchema:
		return "schema"
	default:
		return "item"
	}
}

func formatHistorySecondary(item findHistoryItem, includeDatabase bool) string {
	if includeDatabase {
		return item.database
	}
	return ""
}

func (home *Home) runHistoryEntry(entry findHistoryItem) error {
	switch entry.kind {
	case findKindTable:
		home.showTable(entry.database, joinSchema(entry.schema, entry.name))
		return nil
	case findKindColumn:
		return home.describeTableByName(entry.database, joinSchema(entry.schema, entry.name))
	case findKindFunction:
		return home.describeFunctionByName(entry.database, joinSchema(entry.schema, entry.name))
	case findKindView:
		return home.describeViewByName(entry.database, joinSchema(entry.schema, entry.name))
	case findKindSchema:
		home.showStatusInfo(fmt.Sprintf("Schema: %s", entry.name))
		return nil
	default:
		return errors.New("unsupported history entry")
	}
}

func joinSchema(schema string, name string) string {
	if schema == "" {
		return name
	}
	return fmt.Sprintf("%s.%s", schema, name)
}
