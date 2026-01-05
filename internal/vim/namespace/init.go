package namespace

import "github.com/lancekrogers/lazysql/internal/vim/leader"

func Initialize(registry *leader.Registry, namespaces ...Namespace) (*Registry, error) {
	nsRegistry := NewRegistry(registry)
	for _, ns := range namespaces {
		if err := nsRegistry.Register(ns); err != nil {
			return nil, err
		}
	}
	return nsRegistry, nil
}
