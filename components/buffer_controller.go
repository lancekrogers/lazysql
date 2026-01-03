package components

import "github.com/lancekrogers/lazysql/internal/vim/buffer"

type BufferController struct {
	manager  *buffer.Manager
	editor   *SQLEditor
	suppress bool
}

func NewBufferController(manager *buffer.Manager, editor *SQLEditor) *BufferController {
	controller := &BufferController{
		manager: manager,
		editor:  editor,
	}
	if manager == nil || editor == nil {
		return controller
	}

	editor.SetChangedFunc(func() {
		controller.onEditorChange(editor.GetText())
	})

	manager.AddListener(func(event buffer.BufferEvent) {
		controller.onBufferEvent(event)
	})

	if manager.ActiveID() == 0 {
		manager.Create("")
	} else if buf := manager.Active(); buf != nil {
		controller.syncFromBuffer(buf)
	}

	return controller
}

func (c *BufferController) onEditorChange(text string) {
	if c.suppress || c.manager == nil {
		return
	}
	activeID := c.manager.ActiveID()
	if activeID == 0 {
		return
	}
	_ = c.manager.UpdateContent(activeID, text)
}

func (c *BufferController) onBufferEvent(event buffer.BufferEvent) {
	if c.editor == nil || c.manager == nil || event.Buffer == nil {
		return
	}
	if event.Buffer.ID != c.manager.ActiveID() {
		return
	}
	switch event.Type {
	case buffer.EventActivated, buffer.EventModified:
		c.syncFromBuffer(event.Buffer)
	}
}

func (c *BufferController) syncFromBuffer(buf *buffer.Buffer) {
	if c.editor == nil || buf == nil {
		return
	}
	if c.editor.GetText() == buf.Content {
		return
	}
	c.suppress = true
	c.editor.SetText(buf.Content, true)
	c.suppress = false
}
