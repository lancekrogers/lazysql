package buffer

type NavigationHistory struct {
	stack    []BufferID
	capacity int
}

func NewNavigationHistory(capacity int) *NavigationHistory {
	if capacity <= 0 {
		capacity = 1
	}
	return &NavigationHistory{
		stack:    make([]BufferID, 0, capacity),
		capacity: capacity,
	}
}

func (h *NavigationHistory) Push(id BufferID) {
	h.stack = removeHistoryID(h.stack, id)
	h.stack = append(h.stack, id)
	if len(h.stack) > h.capacity {
		h.stack = h.stack[len(h.stack)-h.capacity:]
	}
}

func (h *NavigationHistory) Alternate() (BufferID, bool) {
	if len(h.stack) < 2 {
		return 0, false
	}
	return h.stack[len(h.stack)-2], true
}

func removeHistoryID(ids []BufferID, target BufferID) []BufferID {
	for i, id := range ids {
		if id == target {
			return append(ids[:i], ids[i+1:]...)
		}
	}
	return ids
}
