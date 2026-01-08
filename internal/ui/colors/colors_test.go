package colors

import (
	"testing"

	"github.com/gdamore/tcell/v2"

	"github.com/lancekrogers/lazysql/internal/ui/icons"
)

func TestThemeManagerGetColorDefaults(t *testing.T) {
	manager := NewThemeManager()
	if manager.Theme() == nil {
		t.Fatal("expected default theme")
	}
	if got := manager.GetColor(icons.TypeDatabase); got != DefaultScheme.Database {
		t.Fatalf("expected default database color, got %v", got)
	}
	if got := manager.GetColor(icons.ObjectType(999)); got != tcell.ColorDefault {
		t.Fatalf("expected default color for unknown type, got %v", got)
	}
}

func TestThemeManagerGetColorAllTypes(t *testing.T) {
	manager := NewThemeManager()
	expectations := map[icons.ObjectType]tcell.Color{
		icons.TypeDatabase:   DefaultScheme.Database,
		icons.TypeSchema:     DefaultScheme.Schema,
		icons.TypeTable:      DefaultScheme.Table,
		icons.TypeView:       DefaultScheme.View,
		icons.TypeColumn:     DefaultScheme.Column,
		icons.TypePrimaryKey: DefaultScheme.PrimaryKey,
		icons.TypeForeignKey: DefaultScheme.ForeignKey,
		icons.TypeIndex:      DefaultScheme.Index,
		icons.TypeFunction:   DefaultScheme.Function,
		icons.TypeProcedure:  DefaultScheme.Procedure,
		icons.TypeTrigger:    DefaultScheme.Trigger,
		icons.TypeSequence:   DefaultScheme.Sequence,
		icons.TypeType:       DefaultScheme.Type,
		icons.TypeExtension:  DefaultScheme.Extension,
	}

	for objType, expected := range expectations {
		if got := manager.GetColor(objType); got != expected {
			t.Fatalf("unexpected color for %v: got %v, want %v", objType, got, expected)
		}
	}
}

func TestThemeManagerSetThemeUnknown(t *testing.T) {
	manager := NewThemeManager()
	if err := manager.SetTheme("missing"); err == nil {
		t.Fatal("expected error for unknown theme")
	}
}

func TestThemeManagerHighContrastToggle(t *testing.T) {
	manager := NewThemeManager()
	if err := manager.SetTheme(ThemeLight); err != nil {
		t.Fatalf("set theme: %v", err)
	}
	if err := manager.SetHighContrast(true); err != nil {
		t.Fatalf("set high contrast: %v", err)
	}
	if got := manager.GetColor(icons.TypeDatabase); got != HighContrastScheme.Database {
		t.Fatalf("expected high contrast color, got %v", got)
	}
	if err := manager.SetHighContrast(false); err != nil {
		t.Fatalf("unset high contrast: %v", err)
	}
	if got := manager.GetColor(icons.TypeDatabase); got != LightScheme.Database {
		t.Fatalf("expected light theme color after reset, got %v", got)
	}
}

func TestThemeManagerNilColor(t *testing.T) {
	var manager *ThemeManager
	if got := manager.GetColor(icons.TypeDatabase); got != tcell.ColorDefault {
		t.Fatalf("expected default color for nil manager, got %v", got)
	}
	if manager.Theme() != nil {
		t.Fatal("expected nil theme for nil manager")
	}
}

func TestThemeManagerRegisterThemeIgnoresInvalid(t *testing.T) {
	manager := NewThemeManager()
	manager.RegisterTheme("", &Scheme{})
	manager.RegisterTheme("ignored", nil)
	if err := manager.SetTheme("ignored"); err == nil {
		t.Fatal("expected invalid theme not to be registered")
	}
}

func TestGlobalThemeHelpers(t *testing.T) {
	t.Cleanup(func() {
		_ = DefaultThemeManager.SetTheme(ThemeDefault)
	})

	if err := SetTheme(ThemeDark); err != nil {
		t.Fatalf("set theme: %v", err)
	}
	if got := ColorFor(icons.TypeTable); got != DarkScheme.Table {
		t.Fatalf("expected dark theme color, got %v", got)
	}
	if err := SetHighContrast(true); err != nil {
		t.Fatalf("set high contrast: %v", err)
	}
	if got := ColorFor(icons.TypeTable); got != HighContrastScheme.Table {
		t.Fatalf("expected high contrast color, got %v", got)
	}
}
