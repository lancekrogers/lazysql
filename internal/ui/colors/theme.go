package colors

import (
	"fmt"

	"github.com/gdamore/tcell/v2"

	"github.com/lancekrogers/lazysql/internal/ui/icons"
)

const (
	ThemeDefault      = "default"
	ThemeDark         = "dark"
	ThemeLight        = "light"
	ThemeHighContrast = "high-contrast"
)

type ThemeManager struct {
	current             *Scheme
	themes              map[string]*Scheme
	highContrastEnabled bool
	lastNonContrast     string
}

func NewThemeManager() *ThemeManager {
	manager := &ThemeManager{
		themes: make(map[string]*Scheme),
	}
	manager.RegisterTheme(ThemeDefault, &DefaultScheme)
	manager.RegisterTheme(ThemeDark, &DarkScheme)
	manager.RegisterTheme(ThemeLight, &LightScheme)
	manager.RegisterTheme(ThemeHighContrast, &HighContrastScheme)
	_ = manager.SetTheme(ThemeDefault)
	return manager
}

func (tm *ThemeManager) RegisterTheme(name string, scheme *Scheme) {
	if tm == nil || name == "" || scheme == nil {
		return
	}
	tm.themes[name] = scheme
}

func (tm *ThemeManager) SetTheme(name string) error {
	if tm == nil {
		return fmt.Errorf("theme manager not configured")
	}
	scheme, ok := tm.themes[name]
	if !ok {
		return fmt.Errorf("unknown theme: %s", name)
	}
	tm.current = scheme
	if name == ThemeHighContrast {
		tm.highContrastEnabled = true
	} else {
		tm.highContrastEnabled = false
		tm.lastNonContrast = name
	}
	return nil
}

func (tm *ThemeManager) SetHighContrast(enabled bool) error {
	if enabled {
		return tm.SetTheme(ThemeHighContrast)
	}
	if tm.lastNonContrast != "" {
		return tm.SetTheme(tm.lastNonContrast)
	}
	return tm.SetTheme(ThemeDefault)
}

func (tm *ThemeManager) Theme() *Scheme {
	if tm == nil {
		return nil
	}
	return tm.current
}

func (tm *ThemeManager) GetColor(objType icons.ObjectType) tcell.Color {
	if tm == nil || tm.current == nil {
		return tcell.ColorDefault
	}
	switch objType {
	case icons.TypeDatabase:
		return tm.current.Database
	case icons.TypeSchema:
		return tm.current.Schema
	case icons.TypeTable:
		return tm.current.Table
	case icons.TypeView:
		return tm.current.View
	case icons.TypeColumn:
		return tm.current.Column
	case icons.TypePrimaryKey:
		return tm.current.PrimaryKey
	case icons.TypeForeignKey:
		return tm.current.ForeignKey
	case icons.TypeIndex:
		return tm.current.Index
	case icons.TypeFunction:
		return tm.current.Function
	case icons.TypeProcedure:
		return tm.current.Procedure
	case icons.TypeTrigger:
		return tm.current.Trigger
	case icons.TypeSequence:
		return tm.current.Sequence
	case icons.TypeType:
		return tm.current.Type
	case icons.TypeExtension:
		return tm.current.Extension
	default:
		return tcell.ColorDefault
	}
}

var DefaultThemeManager = NewThemeManager()

func ColorFor(objType icons.ObjectType) tcell.Color {
	return DefaultThemeManager.GetColor(objType)
}

func SetTheme(name string) error {
	return DefaultThemeManager.SetTheme(name)
}

func SetHighContrast(enabled bool) error {
	return DefaultThemeManager.SetHighContrast(enabled)
}
