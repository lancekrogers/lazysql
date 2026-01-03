package leader

import (
	"errors"
	"sync"

	"github.com/jorgerojas26/lazysql/internal/vim/whichkey"
)

var (
	ErrInactive   = errors.New("leader not active")
	ErrInvalidKey = errors.New("invalid leader key")
	ErrNoTree     = errors.New("leader key tree missing")
	ErrExecuting  = errors.New("leader is executing")
)

type TransitionResult struct {
	NewState LeaderState
	Action   func()
	Sequence []rune
	Error    error
}

type StateMachine struct {
	mu       sync.RWMutex
	state    LeaderState
	sequence []rune
	tree     *whichkey.KeyTree
	current  *whichkey.KeyNode
}

func NewStateMachine(tree *whichkey.KeyTree) *StateMachine {
	var current *whichkey.KeyNode
	if tree != nil {
		current = tree.Root
	}
	return &StateMachine{
		state:    StateIdle,
		sequence: make([]rune, 0),
		tree:     tree,
		current:  current,
	}
}

func (m *StateMachine) State() LeaderState {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.state
}

func (m *StateMachine) Sequence() []rune {
	m.mu.RLock()
	defer m.mu.RUnlock()
	sequence := make([]rune, len(m.sequence))
	copy(sequence, m.sequence)
	return sequence
}

func (m *StateMachine) Activate() bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.state != StateIdle {
		return false
	}
	m.state = StateActive
	m.sequence = m.sequence[:0]
	if m.tree != nil {
		m.current = m.tree.Root
	}
	return true
}

func (m *StateMachine) AddKey(key rune) TransitionResult {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.tree == nil {
		return TransitionResult{NewState: m.state, Error: ErrNoTree}
	}
	if m.state == StateIdle {
		return TransitionResult{NewState: m.state, Error: ErrInactive}
	}
	if m.state == StateExecuting {
		return TransitionResult{NewState: m.state, Error: ErrExecuting}
	}
	if m.current == nil {
		m.current = m.tree.Root
	}

	child, ok := m.current.Children[key]
	if !ok {
		m.resetLocked()
		return TransitionResult{NewState: m.state, Error: ErrInvalidKey}
	}

	m.sequence = append(m.sequence, key)
	m.current = child

	if child.Action != nil {
		m.state = StateExecuting
		return TransitionResult{
			NewState: m.state,
			Action:   child.Action,
			Sequence: append([]rune{}, m.sequence...),
		}
	}

	m.state = StateSequence
	return TransitionResult{
		NewState: m.state,
		Sequence: append([]rune{}, m.sequence...),
	}
}

func (m *StateMachine) Complete() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.resetLocked()
}

func (m *StateMachine) Cancel() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.resetLocked()
}

func (m *StateMachine) resetLocked() {
	m.state = StateIdle
	m.sequence = m.sequence[:0]
	if m.tree != nil {
		m.current = m.tree.Root
	}
}
