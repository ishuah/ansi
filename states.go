package ansi

import "io"

type stateFn func(*Lexer) stateFn

// Status function lexBytes eats up raw bytes, and switches to lexing an espace
// sequence when it encouters the ESC character code.
func lexBytes(l *Lexer) stateFn {
	for l.pos < len(l.input) {
		if l.input[l.pos] == EscapeCode {
			if l.pos > l.start {
				l.emit(RawBytes)
			}
			return lexEscapeSequence
		}
		if _, err := l.next(); err == io.EOF {
			break
		}
	}

	if l.pos > l.start {
		l.emit(RawBytes)
	}
	l.emit(EOF)
	return nil
}

// Status function lexEscapeSequence lexes an sequence starting with the ESC
// character code. It may be a CSI introduced sequence, or a two char sequence.
func lexEscapeSequence(l *Lexer) stateFn {
	l.pos++ // Drop the ESC byte
	next, _ := l.peek()
	if next == '[' {
		return lexControlSequence
	} else if next >= '@' && next <= '_' {
		return lexTwoCharSequence
	}
	return nil
}

// Status function lexTwoCharSequence.
func lexTwoCharSequence(l *Lexer) stateFn {
	l.pos++ // Eat up the command character
	l.emit(TwoCharSequence)
	return lexBytes
}

// Status function lexControlSequence lexes a CSI introduced sequence. If any
// character in the sequence is out of the allowed range, it falls back to
// lexing raw bytes.
func lexControlSequence(l *Lexer) stateFn {
	// General form of a control sequence is:
	// ESC [ n1 ; n2... [trailiing intermediate characters] letter
	// A control sequence is invalid if it doesn't end with a command char: in
	// this case we revert to treat it as a bunch of raw bytes.
	l.accept([]byte{'['})
	l.acceptRunFn(isParamChar)
	l.acceptRunFn(isInterChar)
	if !l.acceptFn(isCommandChar) {
		return l.cancel(lexBytes)
	}

	l.emit(ControlSequence)
	return lexBytes
}
