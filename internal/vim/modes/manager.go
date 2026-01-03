package modes

import "sync"

type ModeManager struct {
	mu        sync.RWMutex
	mode      VimMode
	listeners []func(VimMode)
}

func NewManager() *ModeManager {
	return &ModeManager{
		mode:      ModeNormal,
		listeners: make([]func(VimMode), 0),
	}
}

func (m *ModeManager) Mode() VimMode {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.mode
}

func (m *ModeManager) EnterInsert() bool {
	m.mu.Lock()
	if m.mode != ModeNormal {
		m.mu.Unlock()
		return false
	}
	m.mode = ModeInsert
	listeners := append([]func(VimMode){}, m.listeners...)
	m.mu.Unlock()

	for _, fn := range listeners {
		fn(ModeInsert)
	}
	return true
}

func (m *ModeManager) EnterNormal() bool {
	m.mu.Lock()
	if m.mode == ModeNormal {
		m.mu.Unlock()
		return false
	}
	m.mode = ModeNormal
	listeners := append([]func(VimMode){}, m.listeners...)
	m.mu.Unlock()

	for _, fn := range listeners {
		fn(ModeNormal)
	}
	return true
}

func (m *ModeManager) AddListener(fn func(VimMode)) {
	if fn == nil {
		return
	}
	m.mu.Lock()
	m.listeners = append(m.listeners, fn)
	m.mu.Unlock()
}
