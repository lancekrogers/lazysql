package components

import (
	"github.com/rivo/tview"

	"github.com/lancekrogers/lazysql/app"
	"github.com/lancekrogers/lazysql/keymap"
)

type HelpStatus struct {
	*tview.TextView
}

func NewHelpStatus() HelpStatus {
	status := HelpStatus{tview.NewTextView().SetTextColor(app.Styles.TertiaryTextColor)}

	status.SetStatusOnTree()

	return status
}

func (status *HelpStatus) UpdateText(binds []keymap.Bind) {
	newtext := ""

	for i, key := range binds {

		newtext += key.Cmd.String()

		newtext += ": "

		newtext += key.Key.String()

		islast := i == len(binds)-1

		if !islast {
			newtext += " | "
		}

	}

	status.SetText(newtext)
}

func (status *HelpStatus) SetStatusOnTree() {
	status.UpdateText(app.Keymaps.Global)
}

func (status *HelpStatus) SetStatusOnEditorView() {
	status.UpdateText(app.Keymaps.Group(app.EditorGroup))
}

func (status *HelpStatus) SetStatusOnTableView() {
	status.UpdateText(app.Keymaps.Group(app.TableGroup))
}
