package leader

import (
	"sync"
	"time"
)

const DefaultTimeout = 500 * time.Millisecond

type TimeoutManager struct {
	duration  time.Duration
	timer     *time.Timer
	onTimeout func()
	mu        sync.Mutex
}

func NewTimeoutManager(duration time.Duration, onTimeout func()) *TimeoutManager {
	if duration <= 0 {
		duration = DefaultTimeout
	}
	if onTimeout == nil {
		onTimeout = func() {}
	}
	return &TimeoutManager{
		duration:  duration,
		onTimeout: onTimeout,
	}
}

func (m *TimeoutManager) Start() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.stopLocked()
	m.timer = time.AfterFunc(m.duration, m.onTimeout)
}

func (m *TimeoutManager) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.timer != nil {
		m.timer.Reset(m.duration)
	}
}

func (m *TimeoutManager) Stop() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.stopLocked()
}

func (m *TimeoutManager) SetDuration(duration time.Duration) {
	if duration <= 0 {
		return
	}
	m.mu.Lock()
	m.duration = duration
	m.mu.Unlock()
}

func (m *TimeoutManager) stopLocked() {
	if m.timer != nil {
		m.timer.Stop()
		m.timer = nil
	}
}
