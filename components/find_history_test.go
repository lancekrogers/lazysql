package components

import "testing"

func TestFindHistoryAddDedupes(t *testing.T) {
	history := newFindHistory(5)

	history.add(findHistoryItem{kind: findKindTable, database: "db", name: "users"})
	history.add(findHistoryItem{kind: findKindTable, database: "db", name: "users"})
	history.add(findHistoryItem{kind: findKindView, database: "db", name: "active_users"})

	items := history.list()
	if len(items) != 2 {
		t.Fatalf("expected 2 items, got %d", len(items))
	}
	if items[0].kind != findKindView {
		t.Fatalf("expected most recent item to be view, got %q", items[0].kind)
	}
}

func TestFindHistoryLimit(t *testing.T) {
	history := newFindHistory(1)
	history.add(findHistoryItem{kind: findKindTable, database: "db", name: "users"})
	history.add(findHistoryItem{kind: findKindView, database: "db", name: "active_users"})

	items := history.list()
	if len(items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(items))
	}
	if items[0].kind != findKindView {
		t.Fatalf("expected most recent item to be view, got %q", items[0].kind)
	}
}
