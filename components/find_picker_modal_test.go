package components

import (
	"testing"

	"github.com/rivo/tview"
)

func TestFilterFindItemsEmptyQuery(t *testing.T) {
	items := []FindPickerItem{
		{Label: "users", Secondary: "db1"},
		{Label: "orders", Secondary: "db2"},
		{Label: "user_profiles", Secondary: "db1"},
	}

	filtered := filterFindItems(items, "  ")
	if len(filtered) != len(items) {
		t.Fatalf("expected %d items, got %d", len(items), len(filtered))
	}
	for i, item := range items {
		if filtered[i].Label != item.Label {
			t.Fatalf("expected item %q at index %d, got %q", item.Label, i, filtered[i].Label)
		}
	}
}

func TestFilterFindItemsMatchesPrimary(t *testing.T) {
	items := []FindPickerItem{
		{Label: "users", Secondary: "db1"},
		{Label: "orders", Secondary: "db2"},
		{Label: "user_profiles", Secondary: "db1"},
	}

	filtered := filterFindItems(items, "user")
	if len(filtered) != 2 {
		t.Fatalf("expected 2 items, got %d", len(filtered))
	}
	found := map[string]bool{}
	for _, item := range filtered {
		found[item.Label] = true
	}
	if !found["users"] || !found["user_profiles"] {
		t.Fatalf("expected users and user_profiles, got %+v", found)
	}
}

func TestFilterFindItemsMatchesSecondary(t *testing.T) {
	items := []FindPickerItem{
		{Label: "users", Secondary: "db1"},
		{Label: "orders", Secondary: "db2"},
		{Label: "user_profiles", Secondary: "db1"},
	}

	filtered := filterFindItems(items, "db2")
	if len(filtered) != 1 {
		t.Fatalf("expected 1 item, got %d", len(filtered))
	}
	if filtered[0].Label != "orders" {
		t.Fatalf("expected orders, got %q", filtered[0].Label)
	}
}

func TestFindPickerShowAllowsEmptyItems(t *testing.T) {
	mainPages = tview.NewPages()
	picker := NewFindPicker()

	if err := picker.Show(FindPickerConfig{Title: "Find", Placeholder: "Find...", Items: nil}); err != nil {
		t.Fatalf("expected show to succeed, got %v", err)
	}
	if count := picker.list.GetItemCount(); count != 1 {
		t.Fatalf("expected 1 list item, got %d", count)
	}
	main, _ := picker.list.GetItemText(0)
	if main != "No matches" {
		t.Fatalf("expected No matches, got %q", main)
	}
}

func TestFindPickerShowLoadingThenSetItems(t *testing.T) {
	mainPages = tview.NewPages()
	picker := NewFindPicker()

	if err := picker.ShowLoading(FindPickerConfig{Title: "Find", Placeholder: "Find..."}); err != nil {
		t.Fatalf("expected show loading to succeed, got %v", err)
	}
	main, _ := picker.list.GetItemText(0)
	if main != "Loading..." {
		t.Fatalf("expected Loading..., got %q", main)
	}

	picker.SetItems([]FindPickerItem{{Label: "users"}})
	if count := picker.list.GetItemCount(); count != 1 {
		t.Fatalf("expected 1 list item, got %d", count)
	}
	main, _ = picker.list.GetItemText(0)
	if main != "users" {
		t.Fatalf("expected users, got %q", main)
	}
}

func TestFindPickerSelectCurrentTriggersOnSelect(t *testing.T) {
	mainPages = tview.NewPages()
	picker := NewFindPicker()

	called := 0
	items := []FindPickerItem{
		{Label: "users", OnSelect: func() { called++ }},
	}

	if err := picker.Show(FindPickerConfig{Title: "Find", Placeholder: "Find...", Items: items}); err != nil {
		t.Fatalf("expected show to succeed, got %v", err)
	}

	picker.selectCurrent()
	if called != 1 {
		t.Fatalf("expected on select to be called once, got %d", called)
	}
}

func TestFindPickerDismissRemovesPageAndCallsCancel(t *testing.T) {
	mainPages = tview.NewPages()
	picker := NewFindPicker()
	called := 0

	if err := picker.Show(FindPickerConfig{
		Title:    "Find",
		Items:    []FindPickerItem{{Label: "users"}},
		OnCancel: func() { called++ },
	}); err != nil {
		t.Fatalf("expected show to succeed, got %v", err)
	}

	if !mainPages.HasPage(pageNameFindPicker) {
		t.Fatal("expected find picker page to exist")
	}

	picker.dismiss(true)
	if mainPages.HasPage(pageNameFindPicker) {
		t.Fatal("expected find picker page to be removed")
	}
	if called != 1 {
		t.Fatalf("expected cancel to be called once, got %d", called)
	}
}

func TestFindPickerSelectCurrentWithNoItemsDoesNotPanic(t *testing.T) {
	mainPages = tview.NewPages()
	picker := NewFindPicker()

	if err := picker.Show(FindPickerConfig{Title: "Find", Items: nil}); err != nil {
		t.Fatalf("expected show to succeed, got %v", err)
	}

	picker.selectCurrent()
}
