package whichkey

func formatKeyLabel(key rune) string {
	switch key {
	case '\n':
		return "<Enter>"
	case '\t':
		return "<Tab>"
	case ' ':
		return "<Space>"
	default:
		return string(key)
	}
}
