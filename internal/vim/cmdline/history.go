package cmdline

type CommandHistory struct {
	commands []string
	cursor   int
	capacity int
}

func NewCommandHistory(capacity int) *CommandHistory {
	if capacity <= 0 {
		capacity = 1
	}
	return &CommandHistory{
		commands: make([]string, 0, capacity),
		cursor:   0,
		capacity: capacity,
	}
}

func (h *CommandHistory) Add(command string) {
	if command == "" {
		return
	}
	if len(h.commands) > 0 && h.commands[len(h.commands)-1] == command {
		h.cursor = len(h.commands)
		return
	}
	h.commands = append(h.commands, command)
	if len(h.commands) > h.capacity {
		h.commands = h.commands[len(h.commands)-h.capacity:]
	}
	h.cursor = len(h.commands)
}

func (h *CommandHistory) Previous() (string, bool) {
	if len(h.commands) == 0 || h.cursor <= 0 {
		return "", false
	}
	h.cursor--
	return h.commands[h.cursor], true
}

func (h *CommandHistory) Next() (string, bool) {
	if len(h.commands) == 0 {
		return "", false
	}
	if h.cursor >= len(h.commands)-1 {
		h.cursor = len(h.commands)
		return "", false
	}
	h.cursor++
	return h.commands[h.cursor], true
}

func (h *CommandHistory) ResetCursor() {
	h.cursor = len(h.commands)
}
