package modes

type VimMode uint8

const (
	ModeNormal VimMode = iota
	ModeInsert
)

func (m VimMode) String() string {
	switch m {
	case ModeNormal:
		return "NORMAL"
	case ModeInsert:
		return "INSERT"
	default:
		return "UNKNOWN"
	}
}
