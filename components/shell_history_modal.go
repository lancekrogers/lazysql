package components

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"github.com/lancekrogers/lazysql/app"
	"github.com/lancekrogers/lazysql/commands"
)

type ShellHistoryModal struct {
	tview.Primitive
	historyComponent     *QueryHistoryComponent
	onQuerySelected      func(query string)
	connectionIdentifier string
	grid                 *tview.Grid
}

func NewShellHistoryModal(connectionIdentifier string, onSelect func(query string)) *ShellHistoryModal {
	modal := &ShellHistoryModal{
		onQuerySelected:      onSelect,
		connectionIdentifier: connectionIdentifier,
	}

	modal.historyComponent = NewQueryHistoryComponent(connectionIdentifier, func(query string) {
		mainPages.RemovePage(pageNameShellHistory)
		if onSelect != nil {
			onSelect(query)
		}
	}, nil)

	frame := tview.NewFrame(modal.historyComponent)
	frame.SetBackgroundColor(app.Styles.PrimitiveBackgroundColor)
	frame.SetBorder(true)
	frame.SetBorders(0, 0, 0, 0, 0, 0)

	smallScreenGrid := tview.NewGrid().
		SetRows(1, 0, 1).
		SetColumns(1, 0, 1).
		SetMinSize(1, 1)
	smallScreenGrid.AddItem(frame, 0, 0, 3, 3, 0, 0, true)

	largeScreenGrid := tview.NewGrid().
		SetRows(0, 20, 0).
		SetColumns(0, 150, 0).
		SetMinSize(1, 1)
	largeScreenGrid.AddItem(frame, 1, 1, 1, 1, 0, 0, true)

	mainGrid := tview.NewGrid().
		SetRows(0).
		SetColumns(0)
	mainGrid.AddItem(smallScreenGrid, 0, 0, 1, 1, 0, 0, true)
	mainGrid.AddItem(largeScreenGrid, 0, 0, 1, 1, 0, 100, true)

	modal.grid = mainGrid
	modal.Primitive = modal.grid

	modal.grid.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEsc && !modal.historyComponent.GetIsFiltering() {
			mainPages.RemovePage(pageNameShellHistory)
			return nil
		}

		command := app.Keymaps.Group(app.QueryHistoryGroup).Resolve(event)
		switch command {
		case commands.ToggleQueryHistory, commands.Quit:
			if !modal.historyComponent.GetIsFiltering() {
				mainPages.RemovePage(pageNameShellHistory)
				return nil
			}
		}

		return event
	})

	return modal
}

func (modal *ShellHistoryModal) LoadHistory(connectionIdentifier string) {
	if modal.historyComponent == nil {
		return
	}
	modal.historyComponent.LoadHistory(connectionIdentifier)
}

func (modal *ShellHistoryModal) GetPrimitive() tview.Primitive {
	return modal.Primitive
}
