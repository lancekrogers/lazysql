package buffer

import (
	"errors"
	"fmt"
	"sync"
)

var ErrBufferNotFound = errors.New("buffer not found")

type EventType uint8

const (
	EventCreated EventType = iota
	EventActivated
	EventModified
	EventClosed
	EventSaved
)

type BufferEvent struct {
	Type   EventType
	Buffer *Buffer
}

type Manager struct {
	mu            sync.RWMutex
	buffers       map[BufferID]*Buffer
	order         []BufferID
	active        BufferID
	nextID        BufferID
	untitledCount int
	listeners     []func(BufferEvent)
}

func NewManager() *Manager {
	return &Manager{
		buffers: make(map[BufferID]*Buffer),
		order:   make([]BufferID, 0),
		nextID:  1,
	}
}

func (m *Manager) Create(name string) *Buffer {
	m.mu.Lock()
	if name == "" {
		m.untitledCount++
		name = fmt.Sprintf("Untitled-%d", m.untitledCount)
	}
	id := m.nextID
	m.nextID++
	buf := newBuffer(id, name)
	m.buffers[id] = buf
	m.order = append(m.order, id)
	m.active = id
	listeners := append([]func(BufferEvent){}, m.listeners...)
	m.mu.Unlock()

	notify(listeners, BufferEvent{Type: EventCreated, Buffer: buf})
	notify(listeners, BufferEvent{Type: EventActivated, Buffer: buf})
	return buf
}

func (m *Manager) Get(id BufferID) (*Buffer, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	buf, ok := m.buffers[id]
	if !ok {
		return nil, ErrBufferNotFound
	}
	return buf, nil
}

func (m *Manager) ListBuffers() []*Buffer {
	m.mu.RLock()
	defer m.mu.RUnlock()
	buffers := make([]*Buffer, 0, len(m.order))
	for _, id := range m.order {
		if buf, ok := m.buffers[id]; ok {
			buffers = append(buffers, buf)
		}
	}
	return buffers
}

func (m *Manager) Active() *Buffer {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.buffers[m.active]
}

func (m *Manager) ActiveID() BufferID {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.active
}

func (m *Manager) SetActive(id BufferID) error {
	m.mu.Lock()
	buf, ok := m.buffers[id]
	if !ok {
		m.mu.Unlock()
		return ErrBufferNotFound
	}
	m.active = id
	listeners := append([]func(BufferEvent){}, m.listeners...)
	m.mu.Unlock()

	notify(listeners, BufferEvent{Type: EventActivated, Buffer: buf})
	return nil
}

func (m *Manager) UpdateContent(id BufferID, content string) error {
	m.mu.Lock()
	buf, ok := m.buffers[id]
	if !ok {
		m.mu.Unlock()
		return ErrBufferNotFound
	}
	buf.SetContent(content)
	listeners := append([]func(BufferEvent){}, m.listeners...)
	m.mu.Unlock()

	notify(listeners, BufferEvent{Type: EventModified, Buffer: buf})
	return nil
}

func (m *Manager) MarkSaved(id BufferID) error {
	m.mu.Lock()
	buf, ok := m.buffers[id]
	if !ok {
		m.mu.Unlock()
		return ErrBufferNotFound
	}
	buf.MarkSaved()
	listeners := append([]func(BufferEvent){}, m.listeners...)
	m.mu.Unlock()

	notify(listeners, BufferEvent{Type: EventSaved, Buffer: buf})
	return nil
}

func (m *Manager) Close(id BufferID) error {
	m.mu.Lock()
	buf, ok := m.buffers[id]
	if !ok {
		m.mu.Unlock()
		return ErrBufferNotFound
	}
	delete(m.buffers, id)
	m.order = removeID(m.order, id)
	if m.active == id {
		m.active = m.nextActive(id)
	}
	listeners := append([]func(BufferEvent){}, m.listeners...)
	m.mu.Unlock()

	notify(listeners, BufferEvent{Type: EventClosed, Buffer: buf})
	if m.active != 0 {
		if activeBuf, err := m.Get(m.active); err == nil {
			notify(listeners, BufferEvent{Type: EventActivated, Buffer: activeBuf})
		}
	}
	return nil
}

func (m *Manager) GetDirtyBuffers() []*Buffer {
	m.mu.RLock()
	defer m.mu.RUnlock()
	dirty := make([]*Buffer, 0)
	for _, buf := range m.buffers {
		if buf.Dirty {
			dirty = append(dirty, buf)
		}
	}
	return dirty
}

func (m *Manager) AddListener(fn func(BufferEvent)) {
	if fn == nil {
		return
	}
	m.mu.Lock()
	m.listeners = append(m.listeners, fn)
	m.mu.Unlock()
}

func (m *Manager) nextActive(closedID BufferID) BufferID {
	if len(m.order) == 0 {
		return 0
	}
	for i, id := range m.order {
		if id == closedID {
			if i < len(m.order)-1 {
				return m.order[i+1]
			}
			if i-1 >= 0 {
				return m.order[i-1]
			}
		}
	}
	return m.order[len(m.order)-1]
}

func removeID(ids []BufferID, target BufferID) []BufferID {
	for i, id := range ids {
		if id == target {
			return append(ids[:i], ids[i+1:]...)
		}
	}
	return ids
}

func notify(listeners []func(BufferEvent), event BufferEvent) {
	for _, listener := range listeners {
		listener(event)
	}
}
