package modes

import "github.com/rivo/tview"

type ModeIndicator struct {
	*tview.TextView
	modeManager *ModeManager
}

func NewModeIndicator(mm *ModeManager) *ModeIndicator {
	indicator := &ModeIndicator{
		TextView:    tview.NewTextView(),
		modeManager: mm,
	}
	indicator.SetTextAlign(tview.AlignCenter)
	indicator.SetDynamicColors(true)
	indicator.update(mm.Mode())

	mm.AddListener(indicator.update)

	return indicator
}

func (i *ModeIndicator) update(mode VimMode) {
	text := "-- UNKNOWN --"
	color := "[gray]"

	switch mode {
	case ModeNormal:
		text = "-- NORMAL --"
		color = "[green]"
	case ModeInsert:
		text = "-- INSERT --"
		color = "[blue]"
	}

	i.SetText(color + text + "[-]")
}
