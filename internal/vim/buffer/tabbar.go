package buffer

import (
	"unicode/utf8"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type tabEntry struct {
	bufferID BufferID
	label    string
	dirty    bool
	x        int
	width    int
}

type TabBar struct {
	*tview.Box
	manager     *Manager
	tabs        []tabEntry
	styles      TabStyles
	maxTabWidth int
	onChange    func()
}

func NewTabBar(manager *Manager) *TabBar {
	tabBar := &TabBar{
		Box:         tview.NewBox(),
		manager:     manager,
		styles:      DefaultTabStyles,
		maxTabWidth: 24,
		onChange:    func() {},
	}
	if manager != nil {
		manager.AddListener(func(_ BufferEvent) {
			tabBar.updateTabs()
			tabBar.onChange()
		})
	}
	return tabBar
}

func (t *TabBar) SetStyles(styles TabStyles) {
	t.styles = styles
}

func (t *TabBar) SetOnChange(onChange func()) {
	if onChange == nil {
		t.onChange = func() {}
		return
	}
	t.onChange = onChange
}

func (t *TabBar) SetMaxTabWidth(width int) {
	if width > 0 {
		t.maxTabWidth = width
	}
}

func (t *TabBar) Draw(screen tcell.Screen) {
	t.DrawForSubclass(screen, t)
	x, y, width, _ := t.GetInnerRect()

	t.updateTabs()

	pos := x
	for _, tab := range t.tabs {
		if pos >= x+width {
			break
		}

		style := t.styles.Normal
		if tab.bufferID == t.manager.ActiveID() {
			style = t.styles.Active
		} else if tab.dirty {
			style = t.styles.Dirty
		}

		label := tab.label
		if tab.dirty {
			label = label + " +"
		}

		if utf8.RuneCountInString(label) > t.maxTabWidth {
			label = truncate(label, t.maxTabWidth)
		}

		pos = drawLabel(screen, label, pos, y, x+width, style)
		if pos < x+width {
			screen.SetContent(pos, y, '|', nil, t.styles.Separator)
			pos++
		}
	}
}

func (t *TabBar) MouseHandler() func(action tview.MouseAction, event *tcell.EventMouse, setFocus func(p tview.Primitive)) (bool, tview.Primitive) {
	return t.WrapMouseHandler(func(action tview.MouseAction, event *tcell.EventMouse, _ func(p tview.Primitive)) (bool, tview.Primitive) {
		if action != tview.MouseLeftClick {
			return false, nil
		}
		if t.manager == nil {
			return false, nil
		}

		x, _ := event.Position()
		innerX, _, _, _ := t.GetInnerRect()
		clickX := x - innerX

		for _, tab := range t.tabs {
			if clickX >= tab.x && clickX < tab.x+tab.width {
				_ = t.manager.SetActive(tab.bufferID)
				return true, nil
			}
		}
		return false, nil
	})
}

func (t *TabBar) updateTabs() {
	if t.manager == nil {
		t.tabs = nil
		return
	}
	buffers := t.manager.ListBuffers()
	tabs := make([]tabEntry, 0, len(buffers))
	pos := 0
	for _, buf := range buffers {
		label := buf.Name
		if label == "" {
			label = "Untitled"
		}
		displayLabel := label
		if buf.Dirty {
			displayLabel = displayLabel + " +"
		}
		if utf8.RuneCountInString(displayLabel) > t.maxTabWidth {
			displayLabel = truncate(displayLabel, t.maxTabWidth)
		}
		entry := tabEntry{
			bufferID: buf.ID,
			label:    label,
			dirty:    buf.Dirty,
			x:        pos,
			width:    utf8.RuneCountInString(displayLabel) + 1,
		}
		tabs = append(tabs, entry)
		pos += entry.width
	}
	t.tabs = tabs
}

func truncate(label string, maxWidth int) string {
	runes := []rune(label)
	if len(runes) <= maxWidth {
		return label
	}
	if maxWidth <= 3 {
		return string(runes[:maxWidth])
	}
	return string(runes[:maxWidth-3]) + "..."
}

func drawLabel(screen tcell.Screen, label string, x, y, maxX int, style tcell.Style) int {
	pos := x
	for _, r := range label {
		if pos >= maxX {
			break
		}
		screen.SetContent(pos, y, r, nil, style)
		pos++
	}
	return pos
}
