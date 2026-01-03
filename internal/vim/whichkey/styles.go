package whichkey

import "github.com/gdamore/tcell/v2"

// Styles controls overlay colors.
type Styles struct {
	BorderColor      tcell.Color
	TitleColor       tcell.Color
	KeyColor         tcell.Color
	DescriptionColor tcell.Color
	GroupColor       tcell.Color
	BackgroundColor  tcell.Color
}

var DefaultStyles = Styles{
	BorderColor:      tcell.ColorWhite,
	TitleColor:       tcell.ColorWhite,
	KeyColor:         tcell.ColorYellow,
	DescriptionColor: tcell.ColorWhite,
	GroupColor:       tcell.ColorGreen,
	BackgroundColor:  tcell.ColorBlack,
}
