package components

import (
	"testing"

	"github.com/rivo/tview"
)

func TestCellTextUnescapesBrackets(t *testing.T) {
	cell := tview.NewTableCell(tview.Escape(`["x"]`))
	if got := cellText(cell); got != `["x"]` {
		t.Fatalf("cellText = %q", got)
	}
	if cellText(nil) != "" {
		t.Fatal("nil cell should be empty")
	}
}

func TestGetPrimaryKeyValue(t *testing.T) {
	table := &ResultsTable{state: &ResultsTableState{
		records: [][]string{
			{"id", "name"},
			{"1", "ada"},
		},
		columns: [][]string{
			{"Field", "Type"},
			{"id", "int"},
			{"name", "text"},
		},
		primaryKeyColumnNames: []string{"id"},
	}}

	got := table.GetPrimaryKeyValue(1)
	if len(got) != 1 || got[0].Name != "id" || got[0].Value != "1" {
		t.Fatalf("primary key = %#v", got)
	}
	if out := table.GetPrimaryKeyValue(3); len(out) != 0 {
		t.Fatalf("out of range = %#v", out)
	}

	table.state.primaryKeyColumnNames = nil
	got = table.GetPrimaryKeyValue(1)
	if len(got) != 2 || got[0].Value != "1" || got[1].Value != "ada" {
		t.Fatalf("no primary key = %#v", got)
	}

	table.state.records = [][]string{{"id", "name"}, {"only-id"}}
	got = table.GetPrimaryKeyValue(1)
	if len(got) != 1 || got[0].Value != "only-id" {
		t.Fatalf("short row = %#v", got)
	}
}
