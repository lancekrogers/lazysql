package tree

import (
	"strings"
	"testing"

	"github.com/gdamore/tcell/v2"

	"github.com/lancekrogers/lazysql/internal/ui/colors"
	"github.com/lancekrogers/lazysql/internal/ui/icons"
)

func TestRendererRenderDefaults(t *testing.T) {
	renderer := NewRenderer(colors.NewThemeManager())
	result := renderer.Render(icons.TypeTable, "users")
	if !strings.Contains(result.Text, "users") {
		t.Fatalf("expected label to include name, got %q", result.Text)
	}
	if !strings.Contains(result.Text, icons.Icon(icons.TypeTable)) {
		t.Fatalf("expected label to include icon, got %q", result.Text)
	}

	foreground, _, _ := result.TextStyle.Decompose()
	if foreground != colors.DefaultScheme.Table {
		t.Fatalf("expected text style to use table color, got %v", foreground)
	}
}

func TestRendererIconOverride(t *testing.T) {
	renderer := NewRenderer(colors.NewThemeManager())
	renderer.SetIconColorOverride(icons.TypeTable, tcell.ColorRed)

	result := renderer.Render(icons.TypeTable, "users")
	expectedTag := "[" + colorToTag(tcell.ColorRed) + "]"
	if !strings.Contains(result.Text, expectedTag) {
		t.Fatalf("expected override tag %q in %q", expectedTag, result.Text)
	}
}

func TestRendererHandlesUnknownTypes(t *testing.T) {
	renderer := NewRenderer(colors.NewThemeManager())
	result := renderer.Render(icons.ObjectType(999), "unknown")
	if !strings.Contains(result.Text, "unknown") {
		t.Fatalf("expected label to include name, got %q", result.Text)
	}
}

func TestRendererNilTheme(t *testing.T) {
	renderer := NewRenderer(nil)
	result := renderer.Render(icons.TypeTable, "users")
	if result.Text == "" {
		t.Fatal("expected label text")
	}
}

func TestRendererNilThemeManager(t *testing.T) {
	renderer := &Renderer{}
	result := renderer.Render(icons.TypeTable, "users")
	foreground, _, _ := result.TextStyle.Decompose()
	if foreground != tcell.ColorDefault {
		t.Fatalf("expected default color for nil theme, got %v", foreground)
	}
}

func TestRendererClearIconOverride(t *testing.T) {
	renderer := NewRenderer(colors.NewThemeManager())
	renderer.SetIconColorOverride(icons.TypeTable, tcell.ColorRed)
	renderer.ClearIconColorOverride(icons.TypeTable)

	result := renderer.Render(icons.TypeTable, "users")
	if strings.Contains(result.Text, "["+colorToTag(tcell.ColorRed)+"]") {
		t.Fatal("expected override tag to be cleared")
	}
}

func TestColorToTagDefault(t *testing.T) {
	if colorToTag(tcell.ColorDefault) != "default" {
		t.Fatal("expected default color tag")
	}
}
