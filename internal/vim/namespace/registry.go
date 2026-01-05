package namespace

import (
	"fmt"
	"sync"

	"github.com/lancekrogers/lazysql/internal/vim/leader"
)

type Registry struct {
	mu         sync.RWMutex
	leader     *leader.Registry
	namespaces map[rune]Namespace
}

func NewRegistry(leaderRegistry *leader.Registry) *Registry {
	return &Registry{
		leader:     leaderRegistry,
		namespaces: make(map[rune]Namespace),
	}
}

func (r *Registry) Register(ns Namespace) error {
	if ns == nil {
		return fmt.Errorf("namespace is nil")
	}
	prefix := ns.Prefix()
	if prefix == 0 {
		return fmt.Errorf("namespace prefix required")
	}

	r.mu.Lock()
	if _, exists := r.namespaces[prefix]; exists {
		r.mu.Unlock()
		return fmt.Errorf("namespace prefix %c already registered", prefix)
	}
	r.namespaces[prefix] = ns
	r.mu.Unlock()

	if r.leader != nil {
		if err := r.leader.RegisterGroup([]rune{prefix}, ns.Name()); err != nil {
			return err
		}
		for _, cmd := range ns.Commands() {
			seq := append([]rune{prefix}, cmd.Sequence...)
			if err := r.leader.Register(leader.Command{
				Sequence:    seq,
				Description: cmd.Description,
				Handler:     cmd.Handler,
				Category:    ns.Name(),
			}); err != nil {
				return err
			}
		}
	}
	return nil
}

func (r *Registry) Get(prefix rune) (Namespace, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	ns, ok := r.namespaces[prefix]
	return ns, ok
}
