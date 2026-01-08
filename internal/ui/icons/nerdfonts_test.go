package icons

import (
	"os"
	"testing"
)

func TestIconMappingsCoverAllTypes(t *testing.T) {
	types := []ObjectType{
		TypeDatabase,
		TypeSchema,
		TypeTable,
		TypeView,
		TypeColumn,
		TypeIndex,
		TypePrimaryKey,
		TypeForeignKey,
		TypeFunction,
		TypeProcedure,
		TypeTrigger,
		TypeSequence,
		TypeType,
		TypeExtension,
	}

	for _, objType := range types {
		if icon := nerdFontIcons[objType]; icon == "" {
			t.Fatalf("missing nerd font icon for %v", objType)
		}
		if icon := asciiIcons[objType]; icon == "" {
			t.Fatalf("missing ascii icon for %v", objType)
		}
	}
}

func TestIconFallbacks(t *testing.T) {
	previous := useNerdFonts
	defer SetNerdFonts(previous)

	SetNerdFonts(true)
	if icon := Icon(TypeDatabase); icon != IconDatabase {
		t.Fatalf("expected nerd icon for database, got %q", icon)
	}

	SetNerdFonts(false)
	if icon := Icon(TypeDatabase); icon != ASCIIDatabase {
		t.Fatalf("expected ascii icon for database, got %q", icon)
	}

	if icon := Icon(ObjectType(999)); icon != asciiUnknown {
		t.Fatalf("expected fallback icon %q, got %q", asciiUnknown, icon)
	}
}

func TestApplyEnvConfig(t *testing.T) {
	previous := useNerdFonts
	t.Cleanup(func() {
		useNerdFonts = previous
		_ = os.Unsetenv("LAZYSQL_NERD_FONTS")
	})

	if err := os.Setenv("LAZYSQL_NERD_FONTS", "false"); err != nil {
		t.Fatalf("set env: %v", err)
	}
	useNerdFonts = true
	applyEnvConfig()
	if useNerdFonts {
		t.Fatal("expected env to disable nerd fonts")
	}

	if err := os.Setenv("LAZYSQL_NERD_FONTS", "notabool"); err != nil {
		t.Fatalf("set env: %v", err)
	}
	useNerdFonts = true
	applyEnvConfig()
	if !useNerdFonts {
		t.Fatal("expected invalid env value to leave setting unchanged")
	}
}
