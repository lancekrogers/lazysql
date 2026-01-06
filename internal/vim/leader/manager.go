package leader

import (
	"time"

	"github.com/gdamore/tcell/v2"

	"github.com/lancekrogers/lazysql/internal/vim/whichkey"
)

type StatusReporter interface {
	Info(message string)
	Error(message string)
}

type Manager struct {
	machine   *StateMachine
	timeout   *TimeoutManager
	overlay   *whichkey.WhichKeyOverlay
	status    StatusReporter
	leaderKey rune
}

func NewManager(registry *Registry, overlay *whichkey.WhichKeyOverlay, timeout time.Duration, status StatusReporter) *Manager {
	var tree *whichkey.KeyTree
	if registry != nil {
		tree = registry.Tree()
	}
	manager := &Manager{
		machine:   NewStateMachine(tree),
		overlay:   overlay,
		status:    status,
		leaderKey: '\\',
	}
	manager.timeout = NewTimeoutManager(timeout, manager.onTimeout)
	if manager.overlay != nil && tree != nil {
		manager.overlay.SetKeyTree(tree)
		manager.overlay.SetLeaderPrefix(manager.leaderKey)
	}
	if registry != nil {
		registry.SetErrorHandler(func(err error) {
			if err != nil && manager.status != nil {
				manager.status.Error(err.Error())
			}
		})
	}
	return manager
}

func (m *Manager) Activate() {
	if m.machine == nil {
		return
	}
	if !m.machine.Activate() {
		return
	}
	if m.overlay != nil {
		m.overlay.Hide()
	}
	m.timeout.Start()
}

func (m *Manager) HandleEvent(event *tcell.EventKey) bool {
	if event == nil || m.machine == nil {
		return false
	}

	if m.machine.State() == StateIdle {
		if event.Key() == tcell.KeyRune && event.Rune() == m.leaderKey {
			m.Activate()
			return true
		}
		return false
	}

	switch event.Key() {
	case tcell.KeyEscape:
		m.Cancel()
		return true
	case tcell.KeyBackspace, tcell.KeyBackspace2:
		m.Cancel()
		return true
	}

	keyRune, ok := leaderRune(event)
	if !ok {
		return true
	}

	result := m.machine.AddKey(keyRune)
	if result.Error != nil {
		m.Cancel()
		return true
	}

	if result.Action != nil {
		m.timeout.Stop()
		if m.overlay != nil {
			m.overlay.Hide()
		}
		result.Action()
		m.machine.Complete()
		return true
	}

	m.timeout.Reset()
	if m.overlay != nil && m.overlay.IsVisible() {
		m.overlay.Show()
	}
	return true
}

func (m *Manager) Cancel() {
	if m.machine == nil {
		return
	}
	m.machine.Cancel()
	m.timeout.Stop()
	if m.overlay != nil {
		m.overlay.Hide()
	}
}

func (m *Manager) onTimeout() {
	if m.machine == nil || m.overlay == nil {
		return
	}
	state := m.machine.State()
	if state == StateActive || state == StateSequence {
		m.overlay.Show()
	}
}

func leaderRune(event *tcell.EventKey) (rune, bool) {
	if event == nil {
		return 0, false
	}
	if event.Key() == tcell.KeyRune {
		return event.Rune(), true
	}
	if event.Key() == tcell.KeyEnter {
		return '\n', true
	}
	return 0, false
}
