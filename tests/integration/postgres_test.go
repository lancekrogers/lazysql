//go:build integration
// +build integration

package integration

import (
	"strings"
	"testing"

	"github.com/lancekrogers/lazysql/drivers"
)

func TestPostgresDriverIntegration(t *testing.T) {
	container := GetSharedContainer(t)

	driver := &drivers.Postgres{}
	if err := driver.Connect(container.DSN()); err != nil {
		t.Fatalf("connect: %v", err)
	}
	t.Cleanup(func() {
		if driver.Connection != nil {
			_ = driver.Connection.Close()
		}
	})

	if _, err := driver.ExecuteDMLStatement("CREATE TABLE public.test_items (id SERIAL PRIMARY KEY, name TEXT NOT NULL)"); err != nil {
		t.Fatalf("create table: %v", err)
	}

	insertResult, err := driver.ExecuteDMLStatement("INSERT INTO public.test_items (name) VALUES ('alpha'), ('beta')")
	if err != nil {
		t.Fatalf("insert rows: %v", err)
	}
	if !strings.Contains(insertResult, "rows affected") {
		t.Fatalf("unexpected insert result: %q", insertResult)
	}

	databases, err := driver.GetDatabases()
	if err != nil {
		t.Fatalf("get databases: %v", err)
	}
	if !containsString(databases, container.Database()) {
		t.Fatalf("expected database %q in %v", container.Database(), databases)
	}

	tables, err := driver.GetTables(container.Database())
	if err != nil {
		t.Fatalf("get tables: %v", err)
	}
	if !containsString(tables["public"], "test_items") {
		t.Fatalf("expected test_items in tables: %v", tables)
	}

	columns, err := driver.GetTableColumns(container.Database(), "public.test_items")
	if err != nil {
		t.Fatalf("get columns: %v", err)
	}
	if len(columns) < 2 {
		t.Fatalf("expected column rows, got %d", len(columns))
	}
	if !containsColumn(columns[1:], "id") || !containsColumn(columns[1:], "name") {
		t.Fatalf("expected columns id/name, got %v", columns)
	}

	results, count, err := driver.ExecuteQuery("SELECT name FROM public.test_items ORDER BY id")
	if err != nil {
		t.Fatalf("execute query: %v", err)
	}
	if count != 2 {
		t.Fatalf("expected 2 rows, got %d", count)
	}
	if len(results) < 3 || results[0][0] != "name" {
		t.Fatalf("unexpected query results: %v", results)
	}
}

func containsString(items []string, target string) bool {
	for _, item := range items {
		if item == target {
			return true
		}
	}
	return false
}

func containsColumn(rows [][]string, column string) bool {
	for _, row := range rows {
		if len(row) > 0 && row[0] == column {
			return true
		}
	}
	return false
}
