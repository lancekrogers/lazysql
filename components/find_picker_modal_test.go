package components

import "testing"

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
