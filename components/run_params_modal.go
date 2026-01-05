package components

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"github.com/lancekrogers/lazysql/app"
)

type RunParamsModal struct {
	tview.Primitive
	form  *tview.Form
	grid  *tview.Grid
	onRun func(params string)
}

func NewRunParamsModal(onRun func(params string)) *RunParamsModal {
	modal := &RunParamsModal{
		onRun: onRun,
	}

	modal.form = tview.NewForm().
		AddInputField("Params", "", 60, nil, nil).
		AddButton("Run", modal.run).
		AddButton("Cancel", modal.cancel).SetFieldStyle(
		tcell.StyleDefault.
			Background(app.Styles.SecondaryTextColor).
			Foreground(app.Styles.ContrastSecondaryTextColor),
	).SetButtonActivatedStyle(tcell.StyleDefault.
		Background(app.Styles.SecondaryTextColor).
		Foreground(app.Styles.ContrastSecondaryTextColor),
	).SetButtonStyle(tcell.StyleDefault.
		Background(app.Styles.InverseTextColor).
		Foreground(app.Styles.ContrastSecondaryTextColor),
	)

	modal.form.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEsc {
			modal.cancel()
			return nil
		}
		if event.Key() == tcell.KeyEnter {
			modal.run()
			return nil
		}
		return event
	})

	modal.form.SetBorder(true).SetTitle(" Query Parameters ").SetTitleAlign(tview.AlignLeft)

	modal.grid = tview.NewGrid().
		SetRows(0, 7, 0).
		SetColumns(0, 70, 0).
		AddItem(modal.form, 1, 1, 1, 1, 0, 0, true)

	modal.Primitive = modal.grid

	return modal
}

func (modal *RunParamsModal) run() {
	params := modal.form.GetFormItem(0).(*tview.InputField).GetText()
	if modal.onRun != nil {
		modal.onRun(params)
	}
	mainPages.RemovePage(pageNameRunParams)
}

func (modal *RunParamsModal) cancel() {
	mainPages.RemovePage(pageNameRunParams)
}

func (modal *RunParamsModal) GetPrimitive() tview.Primitive {
	return modal.Primitive
}
