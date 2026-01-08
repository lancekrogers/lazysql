package tree

import (
	"fmt"

	"github.com/gdamore/tcell/v2"

	"github.com/lancekrogers/lazysql/internal/ui/colors"
	"github.com/lancekrogers/lazysql/internal/ui/icons"
)

type RenderedNode struct {
	Text              string
	TextStyle         tcell.Style
	SelectedTextStyle tcell.Style
}

type Renderer struct {
	theme         *colors.ThemeManager
	iconOverrides map[icons.ObjectType]tcell.Color
}

func NewRenderer(theme *colors.ThemeManager) *Renderer {
	if theme == nil {
		theme = colors.DefaultThemeManager
	}
	return &Renderer{
		theme:         theme,
		iconOverrides: make(map[icons.ObjectType]tcell.Color),
	}
}

func (r *Renderer) SetIconColorOverride(objType icons.ObjectType, color tcell.Color) {
	if r == nil {
		return
	}
	if r.iconOverrides == nil {
		r.iconOverrides = make(map[icons.ObjectType]tcell.Color)
	}
	r.iconOverrides[objType] = color
}

func (r *Renderer) ClearIconColorOverride(objType icons.ObjectType) {
	if r == nil || r.iconOverrides == nil {
		return
	}
	delete(r.iconOverrides, objType)
}

func (r *Renderer) Render(objType icons.ObjectType, name string) RenderedNode {
	icon := icons.Icon(objType)
	textColor := r.colorFor(objType)
	selectedBg := r.selectedBackground()

	textStyle := tcell.StyleDefault.Foreground(textColor)
	selectedStyle := tcell.StyleDefault.Foreground(textColor).Background(selectedBg)

	label := name
	if icon != "" {
		label = fmt.Sprintf("%s %s", icon, name)
	}

	if override, ok := r.iconOverrides[objType]; ok && icon != "" {
		tag := colorToTag(override)
		if tag != "" {
			label = fmt.Sprintf("[%s]%s[-] %s", tag, icon, name)
		}
	}

	return RenderedNode{
		Text:              label,
		TextStyle:         textStyle,
		SelectedTextStyle: selectedStyle,
	}
}

func (r *Renderer) colorFor(objType icons.ObjectType) tcell.Color {
	if r == nil || r.theme == nil {
		return tcell.ColorDefault
	}
	return r.theme.GetColor(objType)
}

func (r *Renderer) selectedBackground() tcell.Color {
	if r == nil || r.theme == nil {
		return tcell.ColorDefault
	}
	theme := r.theme.Theme()
	if theme == nil {
		return tcell.ColorDefault
	}
	return theme.Selected
}

func colorToTag(color tcell.Color) string {
	if color == tcell.ColorDefault {
		return "default"
	}
	if name := color.Name(true); name != "" {
		return name
	}
	return "default"
}
