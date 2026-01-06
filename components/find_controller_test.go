package components

import (
	"errors"
	"strings"
	"testing"

	"github.com/lancekrogers/lazysql/internal/vim/buffer"
	"github.com/lancekrogers/lazysql/models"
)

type fakeFindDriver struct {
	databases           []string
	tables              map[string]map[string][]string
	columns             map[string][][]string
	functions           map[string]map[string][]string
	views               map[string]map[string][]string
	useSchemas          bool
	supportsProgramming bool
	provider            string
}

func (d *fakeFindDriver) Connect(string) error        { return nil }
func (d *fakeFindDriver) TestConnection(string) error { return nil }
func (d *fakeFindDriver) GetDatabases() ([]string, error) {
	return append([]string(nil), d.databases...), nil
}
func (d *fakeFindDriver) GetProvider() string                              { return d.provider }
func (d *fakeFindDriver) SetProvider(provider string)                      { d.provider = provider }
func (d *fakeFindDriver) SupportsProgramming() bool                        { return d.supportsProgramming }
func (d *fakeFindDriver) UseSchemas() bool                                 { return d.useSchemas }
func (d *fakeFindDriver) ExecuteDMLStatement(string) (string, error)       { return "", nil }
func (d *fakeFindDriver) ExecuteQuery(string) ([][]string, int, error)     { return nil, 0, nil }
func (d *fakeFindDriver) ExecutePendingChanges([]models.DBDMLChange) error { return nil }
func (d *fakeFindDriver) UpdateRecord(string, string, string, string, string, string) error {
	return nil
}
func (d *fakeFindDriver) DeleteRecord(string, string, string, string) error { return nil }
func (d *fakeFindDriver) GetConstraints(string, string) ([][]string, error) {
	return nil, nil
}
func (d *fakeFindDriver) GetForeignKeys(string, string) ([][]string, error) {
	return nil, nil
}
func (d *fakeFindDriver) GetIndexes(string, string) ([][]string, error) {
	return nil, nil
}
func (d *fakeFindDriver) GetRecords(string, string, string, string, int, int) ([][]string, int, string, error) {
	return [][]string{{"id"}, {"1"}}, 1, "", nil
}
func (d *fakeFindDriver) GetPrimaryKeyColumnNames(string, string) ([]string, error) {
	return []string{"id"}, nil
}
func (d *fakeFindDriver) GetFunctions(database string) (map[string][]string, error) {
	if !d.supportsProgramming {
		return nil, errors.New("not supported")
	}
	return copyNamespaceMap(d.functions[database], d.useSchemas, database), nil
}
func (d *fakeFindDriver) GetProcedures(string) (map[string][]string, error) { return nil, nil }
func (d *fakeFindDriver) GetViews(database string) (map[string][]string, error) {
	if !d.supportsProgramming {
		return nil, errors.New("not supported")
	}
	return copyNamespaceMap(d.views[database], d.useSchemas, database), nil
}
func (d *fakeFindDriver) GetFunctionDefinition(string, string) (string, error)  { return "", nil }
func (d *fakeFindDriver) GetProcedureDefinition(string, string) (string, error) { return "", nil }
func (d *fakeFindDriver) GetViewDefinition(string, string) (string, error)      { return "", nil }
func (d *fakeFindDriver) FormatArg(arg any, _ models.CellValueType) any         { return arg }
func (d *fakeFindDriver) FormatArgForQueryString(_ any) string                  { return "" }
func (d *fakeFindDriver) FormatReference(reference string) string               { return reference }
func (d *fakeFindDriver) FormatPlaceholder(_ int) string                        { return "?" }
func (d *fakeFindDriver) DMLChangeToQueryString(models.DBDMLChange) (string, error) {
	return "", nil
}

func (d *fakeFindDriver) GetTables(database string) (map[string][]string, error) {
	if d.tables == nil {
		return nil, nil
	}
	entries := d.tables[database]
	return copyNamespaceMap(entries, d.useSchemas, database), nil
}

func (d *fakeFindDriver) GetTableColumns(database, table string) ([][]string, error) {
	key := database + "|" + table
	if columns, ok := d.columns[key]; ok {
		return columns, nil
	}
	return [][]string{{"column_name"}, {"id"}}, nil
}

func copyNamespaceMap(entries map[string][]string, useSchemas bool, database string) map[string][]string {
	result := make(map[string][]string)
	if entries == nil {
		return result
	}
	if useSchemas {
		for key, values := range entries {
			result[key] = append([]string(nil), values...)
		}
		return result
	}

	var flattened []string
	for _, values := range entries {
		flattened = append(flattened, values...)
	}
	result[database] = flattened
	return result
}

func TestFindTableItems(t *testing.T) {
	driver := &fakeFindDriver{
		databases:  []string{"db1", "db2"},
		useSchemas: false,
		tables: map[string]map[string][]string{
			"db1": {"db1": {"users"}},
			"db2": {"db2": {"orders"}},
		},
	}
	home := &Home{DBDriver: driver, findHistory: newFindHistory(10)}

	items, err := home.findTableItems()
	if err != nil {
		t.Fatalf("findTableItems error: %v", err)
	}
	if len(items) != 2 {
		t.Fatalf("expected 2 items, got %d", len(items))
	}

	want := map[string]string{"users": "db1", "orders": "db2"}
	for _, item := range items {
		if secondary, ok := want[item.Label]; ok {
			if item.Secondary != secondary {
				t.Fatalf("expected secondary %q for %q, got %q", secondary, item.Label, item.Secondary)
			}
			delete(want, item.Label)
		}
	}
	if len(want) != 0 {
		t.Fatalf("missing items: %+v", want)
	}
}

func TestFindColumnItems(t *testing.T) {
	driver := &fakeFindDriver{
		databases:  []string{"db1"},
		useSchemas: true,
		tables: map[string]map[string][]string{
			"db1": {"public": {"users"}},
		},
		columns: map[string][][]string{
			"db1|public.users": {{"column_name"}, {"id"}, {"email"}},
		},
	}
	home := &Home{DBDriver: driver, findHistory: newFindHistory(10)}

	items, err := home.findColumnItems()
	if err != nil {
		t.Fatalf("findColumnItems error: %v", err)
	}
	if len(items) != 2 {
		t.Fatalf("expected 2 items, got %d", len(items))
	}

	labels := []string{items[0].Label, items[1].Label}
	if !containsLabel(labels, "public.users.email") || !containsLabel(labels, "public.users.id") {
		t.Fatalf("unexpected labels: %+v", labels)
	}
}

func TestFindFunctionAndViewItems(t *testing.T) {
	driver := &fakeFindDriver{
		databases:           []string{"db1"},
		useSchemas:          false,
		supportsProgramming: true,
		functions: map[string]map[string][]string{
			"db1": {"db1": {"fn_one"}},
		},
		views: map[string]map[string][]string{
			"db1": {"db1": {"view_one"}},
		},
	}
	home := &Home{DBDriver: driver, findHistory: newFindHistory(10)}

	functions, err := home.findFunctionItems()
	if err != nil {
		t.Fatalf("findFunctionItems error: %v", err)
	}
	if len(functions) != 1 || functions[0].Label != "fn_one" {
		t.Fatalf("unexpected functions: %+v", functions)
	}

	views, err := home.findViewItems()
	if err != nil {
		t.Fatalf("findViewItems error: %v", err)
	}
	if len(views) != 1 || views[0].Label != "view_one" {
		t.Fatalf("unexpected views: %+v", views)
	}
}

func TestFindSchemaItems(t *testing.T) {
	driver := &fakeFindDriver{
		databases:  []string{"db1"},
		useSchemas: true,
		tables: map[string]map[string][]string{
			"db1": {"public": {"users"}, "audit": {"events"}},
		},
	}
	home := &Home{DBDriver: driver, findHistory: newFindHistory(10)}

	items, err := home.findSchemaItems()
	if err != nil {
		t.Fatalf("findSchemaItems error: %v", err)
	}
	if len(items) != 2 {
		t.Fatalf("expected 2 items, got %d", len(items))
	}
	labels := []string{items[0].Label, items[1].Label}
	if !containsLabel(labels, "audit") || !containsLabel(labels, "public") {
		t.Fatalf("unexpected schema labels: %+v", labels)
	}
}

func TestFindQueryItems(t *testing.T) {
	manager := buffer.NewManager()
	buf := manager.Create("Query")
	_ = manager.UpdateContent(buf.ID, "select * from users;\n\nselect * from orders;")

	home := &Home{BufferManager: manager}
	items, err := home.findQueryItems()
	if err != nil {
		t.Fatalf("findQueryItems error: %v", err)
	}
	if len(items) != 2 {
		t.Fatalf("expected 2 items, got %d", len(items))
	}
	if !strings.Contains(items[0].Label, "Query:1") {
		t.Fatalf("expected label to include line 1, got %q", items[0].Label)
	}
	if !strings.Contains(items[1].Label, "Query:3") {
		t.Fatalf("expected label to include line 3, got %q", items[1].Label)
	}
}

func TestFindRecentItems(t *testing.T) {
	history := newFindHistory(10)
	home := &Home{findHistory: history}

	history.add(findHistoryItem{kind: findKindTable, database: "db1", name: "users"})
	history.add(findHistoryItem{kind: findKindTable, database: "db2", name: "orders"})

	items, err := home.findRecentItems()
	if err != nil {
		t.Fatalf("findRecentItems error: %v", err)
	}
	if len(items) != 2 {
		t.Fatalf("expected 2 items, got %d", len(items))
	}
	if items[0].Secondary == "" || items[1].Secondary == "" {
		t.Fatalf("expected secondary database labels, got %+v", items)
	}
}

func containsLabel(labels []string, target string) bool {
	for _, label := range labels {
		if label == target {
			return true
		}
	}
	return false
}
