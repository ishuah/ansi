package ansi

const (
	// EscapeCode holds the ANSI escape character '\x1B'
	EscapeCode = '\x1B'
)

func isParamChar(b byte) bool {
	return b >= 0x30 && b <= 0x3F
}

func isInterChar(b byte) bool {
	return b >= 0x20 && b <= 0x2F
}

func isCommandChar(b byte) bool {
	return b >= 0x40 && b <= 0x7E
}
