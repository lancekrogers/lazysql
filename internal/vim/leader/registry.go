package leader

import (
	"errors"
	"fmt"
	"sync"

	"github.com/lancekrogers/lazysql/internal/vim/whichkey"
)

type Command struct {
	Sequence    []rune
	Description string
	Handler     func() error
	Category    string
}

type Registry struct {
	mu      sync.RWMutex
	tree    *whichkey.KeyTree
	entries map[string]Command
	onError func(error)
}

func NewRegistry() *Registry {
	return &Registry{
		tree:    whichkey.NewKeyTree(),
		entries: make(map[string]Command),
	}
}

func (r *Registry) Tree() *whichkey.KeyTree {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.tree
}

func (r *Registry) SetErrorHandler(fn func(error)) {
	r.mu.Lock()
	r.onError = fn
	r.mu.Unlock()
}

func (r *Registry) RegisterGroup(sequence []rune, description string) error {
	if len(sequence) == 0 {
		return errors.New("group sequence required")
	}
	r.mu.Lock()
	r.tree.AddGroup(sequence, description)
	r.mu.Unlock()
	return nil
}

func (r *Registry) Register(cmd Command) error {
	if len(cmd.Sequence) == 0 {
		return errors.New("command sequence required")
	}
	if cmd.Handler == nil {
		return errors.New("command handler required")
	}

	key := string(cmd.Sequence)

	r.mu.Lock()
	if _, exists := r.entries[key]; exists {
		r.mu.Unlock()
		return fmt.Errorf("leader command already registered: %s", key)
	}
	wrapped := func() {
		if err := cmd.Handler(); err != nil {
			r.mu.RLock()
			onError := r.onError
			r.mu.RUnlock()
			if onError != nil {
				onError(err)
			}
		}
	}
	r.tree.AddCommand(cmd.Sequence, cmd.Description, wrapped)
	r.entries[key] = cmd
	r.mu.Unlock()
	return nil
}

func (r *Registry) Lookup(sequence []rune) (Command, bool) {
	key := string(sequence)
	r.mu.RLock()
	defer r.mu.RUnlock()
	cmd, ok := r.entries[key]
	return cmd, ok
}
