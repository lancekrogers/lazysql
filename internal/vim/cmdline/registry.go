package cmdline

type Registry struct {
	commands map[string]Command
	aliases  map[string]string
}

func NewRegistry() *Registry {
	return &Registry{
		commands: make(map[string]Command),
		aliases:  make(map[string]string),
	}
}

func (r *Registry) Register(cmd Command) {
	if cmd == nil {
		return
	}
	r.commands[cmd.Name()] = cmd
}

func (r *Registry) RegisterAlias(alias string, name string) {
	if alias == "" || name == "" {
		return
	}
	r.aliases[alias] = name
}

func (r *Registry) Lookup(name string) Command {
	if name == "" {
		return nil
	}
	if alias, ok := r.aliases[name]; ok {
		name = alias
	}
	return r.commands[name]
}
