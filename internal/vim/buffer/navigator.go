package buffer

import "fmt"

type Navigator struct {
	manager *Manager
	history *NavigationHistory
}

func NewNavigator(manager *Manager) *Navigator {
	navigator := &Navigator{
		manager: manager,
		history: NewNavigationHistory(10),
	}
	manager.AddListener(func(event BufferEvent) {
		if event.Type == EventActivated && event.Buffer != nil {
			navigator.history.Push(event.Buffer.ID)
		}
	})
	return navigator
}

func (n *Navigator) Next() error {
	buffers := n.manager.ListBuffers()
	if len(buffers) <= 1 {
		return nil
	}
	currentID := n.manager.ActiveID()
	index := findBufferIndex(buffers, currentID)
	if index < 0 {
		return n.manager.SetActive(buffers[0].ID)
	}
	nextIndex := (index + 1) % len(buffers)
	return n.manager.SetActive(buffers[nextIndex].ID)
}

func (n *Navigator) Previous() error {
	buffers := n.manager.ListBuffers()
	if len(buffers) <= 1 {
		return nil
	}
	currentID := n.manager.ActiveID()
	index := findBufferIndex(buffers, currentID)
	if index < 0 {
		return n.manager.SetActive(buffers[0].ID)
	}
	previousIndex := (index - 1 + len(buffers)) % len(buffers)
	return n.manager.SetActive(buffers[previousIndex].ID)
}

func (n *Navigator) GoTo(id BufferID) error {
	return n.manager.SetActive(id)
}

func (n *Navigator) GoToIndex(index int) error {
	if index <= 0 {
		return fmt.Errorf("invalid buffer index %d", index)
	}
	buffers := n.manager.ListBuffers()
	if index > len(buffers) {
		return fmt.Errorf("buffer index %d out of range", index)
	}
	return n.manager.SetActive(buffers[index-1].ID)
}

func (n *Navigator) Alternate() (BufferID, bool) {
	return n.history.Alternate()
}

func (n *Navigator) List() []*Buffer {
	return n.manager.ListBuffers()
}

func findBufferIndex(buffers []*Buffer, id BufferID) int {
	for i, buffer := range buffers {
		if buffer.ID == id {
			return i
		}
	}
	return -1
}
