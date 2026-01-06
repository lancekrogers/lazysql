package components

import (
	"github.com/rivo/tview"

	"github.com/lancekrogers/lazysql/app"
)

type ConnectionIndicator struct {
	*tview.TextView
}

func NewConnectionIndicator() *ConnectionIndicator {
	view := tview.NewTextView()
	view.SetTextColor(app.Styles.TertiaryTextColor)
	view.SetTextAlign(tview.AlignRight)
	return &ConnectionIndicator{TextView: view}
}

func (c *ConnectionIndicator) SetConnection(name string) {
	if c == nil {
		return
	}
	label := name
	if label == "" {
		label = "Connected"
	}
	c.SetText("Conn: " + label)
}

func (c *ConnectionIndicator) SetDisconnected() {
	if c == nil {
		return
	}
	c.SetText("Conn: Disconnected")
}
