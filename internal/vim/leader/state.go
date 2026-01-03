package leader

type LeaderState uint8

const (
	StateIdle LeaderState = iota
	StateActive
	StateSequence
	StateExecuting
)

func (s LeaderState) String() string {
	switch s {
	case StateIdle:
		return "Idle"
	case StateActive:
		return "Active"
	case StateSequence:
		return "Sequence"
	case StateExecuting:
		return "Executing"
	default:
		return "Unknown"
	}
}
