package buffer

import "github.com/gdamore/tcell/v2"

type TabStyles struct {
	Active    tcell.Style
	Normal    tcell.Style
	Dirty     tcell.Style
	Separator tcell.Style
}

var DefaultTabStyles = TabStyles{
	Active:    tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorBlue),
	Normal:    tcell.StyleDefault.Foreground(tcell.ColorGray),
	Dirty:     tcell.StyleDefault.Foreground(tcell.ColorYellow),
	Separator: tcell.StyleDefault.Foreground(tcell.ColorGray),
}
