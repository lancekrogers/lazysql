package components

import (
	"errors"
	"fmt"

	"github.com/rivo/tview"

	"github.com/lancekrogers/lazysql/app"
	"github.com/lancekrogers/lazysql/internal/vim/buffer"
)

type BufferPicker struct {
	manager *buffer.Manager
}

func NewBufferPicker(manager *buffer.Manager) *BufferPicker {
	return &BufferPicker{manager: manager}
}

func (p *BufferPicker) Show(buffers []*buffer.Buffer) error {
	if p.manager == nil {
		return errors.New("buffer manager not configured")
	}
	if mainPages == nil {
		return errors.New("main pages not configured")
	}

	list := tview.NewList()
	list.SetBorder(true)
	list.SetTitle(" Buffers ")
	list.SetMainTextColor(app.Styles.PrimaryTextColor)
	list.SetSecondaryTextColor(app.Styles.InverseTextColor)
	list.ShowSecondaryText(true)

	previousFocus := app.App.GetFocus()

	for _, buf := range buffers {
		if buf == nil {
			continue
		}
		name := buf.Name
		if buf.Dirty {
			name = fmt.Sprintf("%s [+]", name)
		}
		secondary := buf.FilePath
		if secondary == "" {
			secondary = "unsaved"
		}
		id := buf.ID
		list.AddItem(name, secondary, 0, func() {
			_ = p.manager.SetActive(id)
			mainPages.RemovePage(pageNameBufferPicker)
			if previousFocus != nil {
				app.App.SetFocus(previousFocus)
			}
		})
	}

	list.SetDoneFunc(func() {
		mainPages.RemovePage(pageNameBufferPicker)
		if previousFocus != nil {
			app.App.SetFocus(previousFocus)
		}
	})

	if mainPages.HasPage(pageNameBufferPicker) {
		mainPages.RemovePage(pageNameBufferPicker)
	}
	mainPages.AddPage(pageNameBufferPicker, list, true, true)
	app.App.SetFocus(list)
	return nil
}
