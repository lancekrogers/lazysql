package colors

import "github.com/gdamore/tcell/v2"

type Scheme struct {
	Database   tcell.Color
	Schema     tcell.Color
	Table      tcell.Color
	View       tcell.Color
	Column     tcell.Color
	PrimaryKey tcell.Color
	ForeignKey tcell.Color
	Index      tcell.Color
	Function   tcell.Color
	Procedure  tcell.Color
	Trigger    tcell.Color
	Sequence   tcell.Color
	Type       tcell.Color
	Extension  tcell.Color
	Selected   tcell.Color
	Dimmed     tcell.Color
}

var DefaultScheme = Scheme{
	Database:   tcell.ColorAqua,
	Schema:     tcell.ColorYellow,
	Table:      tcell.ColorBlue,
	View:       tcell.ColorGreen,
	Column:     tcell.ColorWhite,
	PrimaryKey: tcell.ColorGold,
	ForeignKey: tcell.ColorTeal,
	Index:      tcell.ColorPurple,
	Function:   tcell.ColorOrange,
	Procedure:  tcell.ColorOrange,
	Trigger:    tcell.ColorRed,
	Sequence:   tcell.ColorTeal,
	Type:       tcell.ColorSilver,
	Extension:  tcell.ColorFuchsia,
	Selected:   tcell.ColorNavy,
	Dimmed:     tcell.ColorGray,
}

var DarkScheme = Scheme{
	Database:   tcell.ColorAqua,
	Schema:     tcell.ColorYellow,
	Table:      tcell.ColorBlue,
	View:       tcell.ColorGreen,
	Column:     tcell.ColorWhite,
	PrimaryKey: tcell.ColorGold,
	ForeignKey: tcell.ColorTeal,
	Index:      tcell.ColorPurple,
	Function:   tcell.ColorOrange,
	Procedure:  tcell.ColorOrange,
	Trigger:    tcell.ColorRed,
	Sequence:   tcell.ColorTeal,
	Type:       tcell.ColorSilver,
	Extension:  tcell.ColorFuchsia,
	Selected:   tcell.ColorNavy,
	Dimmed:     tcell.ColorGray,
}

var LightScheme = Scheme{
	Database:   tcell.ColorTeal,
	Schema:     tcell.ColorOlive,
	Table:      tcell.ColorNavy,
	View:       tcell.ColorDarkGreen,
	Column:     tcell.ColorBlack,
	PrimaryKey: tcell.ColorMaroon,
	ForeignKey: tcell.ColorDarkCyan,
	Index:      tcell.ColorPurple,
	Function:   tcell.ColorDarkOrange,
	Procedure:  tcell.ColorDarkOrange,
	Trigger:    tcell.ColorRed,
	Sequence:   tcell.ColorDarkCyan,
	Type:       tcell.ColorGray,
	Extension:  tcell.ColorDarkMagenta,
	Selected:   tcell.ColorSilver,
	Dimmed:     tcell.ColorGray,
}

var HighContrastScheme = Scheme{
	Database:   tcell.ColorWhite,
	Schema:     tcell.ColorWhite,
	Table:      tcell.ColorWhite,
	View:       tcell.ColorWhite,
	Column:     tcell.ColorWhite,
	PrimaryKey: tcell.ColorWhite,
	ForeignKey: tcell.ColorWhite,
	Index:      tcell.ColorWhite,
	Function:   tcell.ColorWhite,
	Procedure:  tcell.ColorWhite,
	Trigger:    tcell.ColorWhite,
	Sequence:   tcell.ColorWhite,
	Type:       tcell.ColorWhite,
	Extension:  tcell.ColorWhite,
	Selected:   tcell.ColorBlack,
	Dimmed:     tcell.ColorWhite,
}
