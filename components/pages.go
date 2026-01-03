package components

import (
	"github.com/rivo/tview"

	"github.com/lancekrogers/lazysql/app"
)

var mainPages *tview.Pages

func MainPages() *tview.Pages {
	mainPages = tview.NewPages()
	mainPages.SetBackgroundColor(app.Styles.PrimitiveBackgroundColor)
	mainPages.AddPage(pageNameConnections, NewConnectionPages().Grid, true, true)
	return mainPages
}
