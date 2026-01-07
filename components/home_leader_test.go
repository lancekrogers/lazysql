package components

import (
	"testing"

	"github.com/lancekrogers/lazysql/models"
)

func TestHomeRegistersLeaderToggleTree(t *testing.T) {
	driver := &fakeFindDriver{}
	connection := models.Connection{
		Name:   "test",
		URL:    "postgres://localhost/test",
		DBName: "test",
	}

	home := NewHomePage(connection, driver)
	if home == nil || home.LeaderRegistry == nil {
		t.Fatal("expected home leader registry")
	}

	cmd, ok := home.LeaderRegistry.Lookup([]rune{'e'})
	if !ok {
		t.Fatal("expected leader binding for \\e")
	}
	if cmd.Description != "Toggle tree" {
		t.Fatalf("expected description %q, got %q", "Toggle tree", cmd.Description)
	}
	if cmd.Handler == nil {
		t.Fatal("expected leader handler")
	}
	if err := cmd.Handler(); err != nil {
		t.Fatalf("expected handler to succeed: %v", err)
	}
}
