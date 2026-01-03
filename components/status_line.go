package components

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"github.com/lancekrogers/lazysql/app"
)

type StatusLine struct {
	*tview.TextView
}

func NewStatusLine() *StatusLine {
	line := &StatusLine{TextView: tview.NewTextView()}
	line.SetTextColor(app.Styles.TertiaryTextColor)
	return line
}

func (s *StatusLine) Info(message string) {
	s.SetTextColor(app.Styles.TertiaryTextColor)
	s.SetText(message)
}

func (s *StatusLine) Error(message string) {
	s.SetTextColor(tcell.ColorRed)
	s.SetText(message)
}
