package buffer

import (
	"testing"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func TestTabBarUpdates(t *testing.T) {
	manager := NewManager()
	first := manager.Create("one")
	second := manager.Create("two")
	if err := manager.UpdateContent(second.ID, "select 2"); err != nil {
		t.Fatalf("update content: %v", err)
	}

	tabBar := NewTabBar(manager)
	tabBar.updateTabs()

	if len(tabBar.tabs) != 2 {
		t.Fatalf("expected 2 tabs, got %d", len(tabBar.tabs))
	}
	if tabBar.tabs[0].bufferID != first.ID {
		t.Fatalf("expected first tab to be buffer %d, got %d", first.ID, tabBar.tabs[0].bufferID)
	}
	if !tabBar.tabs[1].dirty {
		t.Fatal("expected second tab to be dirty")
	}
}

func TestTabBarMouseClick(t *testing.T) {
	manager := NewManager()
	first := manager.Create("one")
	second := manager.Create("two")

	tabBar := NewTabBar(manager)
	tabBar.SetRect(0, 0, 40, 1)
	tabBar.updateTabs()

	clickX := tabBar.tabs[1].x + 1
	event := tcell.NewEventMouse(clickX, 0, tcell.Button1, 0)
	handler := tabBar.MouseHandler()
	consumed, _ := handler(tview.MouseLeftClick, event, nil)
	if !consumed {
		t.Fatal("expected click to be consumed")
	}
	if manager.ActiveID() != second.ID {
		t.Fatalf("expected active buffer to be %d, got %d", second.ID, manager.ActiveID())
	}

	_ = manager.SetActive(first.ID)
}

func TestTabBarDrawAndTruncate(t *testing.T) {
	manager := NewManager()
	buffer := manager.Create("very-long-buffer-name")
	if err := manager.UpdateContent(buffer.ID, "select 1"); err != nil {
		t.Fatalf("update content: %v", err)
	}

	tabBar := NewTabBar(manager)
	tabBar.SetMaxTabWidth(5)
	tabBar.SetStyles(DefaultTabStyles)
	tabBar.SetRect(0, 0, 20, 1)

	screen := tcell.NewSimulationScreen("UTF-8")
	if err := screen.Init(); err != nil {
		t.Fatalf("init screen: %v", err)
	}
	screen.SetSize(20, 1)
	tabBar.Draw(screen)
}

func TestTabBarNilManager(t *testing.T) {
	tabBar := NewTabBar(nil)
	tabBar.updateTabs()
	if tabBar.tabs != nil {
		t.Fatal("expected tabs to be nil when manager is nil")
	}
	tabBar.SetOnChange(nil)
	tabBar.SetMaxTabWidth(0)
}

func TestTabBarOnChange(t *testing.T) {
	manager := NewManager()
	tabBar := NewTabBar(manager)
	calls := 0
	tabBar.SetOnChange(func() {
		calls++
	})

	manager.Create("one")
	if calls == 0 {
		t.Fatal("expected onChange to be called on buffer event")
	}
}
