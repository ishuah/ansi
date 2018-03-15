package ansi

import (
	"bytes"
	"io"
)

// The Lexer implementation is entirely taken from Rob Pike's "Lexical Scanning
// in Go" talk (https://www.youtube.com/watch?v=HxaD_trXwRE).

// Lexer implements reading and tokenizing bytes
type Lexer struct {
	items     chan Item
	input     []byte
	start     int
	pos       int
	itemStart int
	state     stateFn
}

// Init loads a Lexer instance with new input
func (l *Lexer) Init(input []byte) {
	l.input = input
	l.items = make(chan Item, 2)
	l.state = lexBytes
}

func (l *Lexer) backup() {
	l.pos--
}

func (l *Lexer) cancel(revert stateFn) stateFn {
	l.pos = l.itemStart
	return revert
}

func (l *Lexer) emit(t ItemType) {
	l.items <- Item{T: t, Value: l.input[l.start:l.pos]}
	l.start = l.pos
}

func (l *Lexer) next() (byte, error) {
	if l.pos >= len(l.input) {
		return 0, io.EOF
	}
	b := l.input[l.pos]
	l.pos++
	return b, nil
}

// NextItem returns the the next token from the source. Returns
// EOF at the end of the source.
func (l *Lexer) NextItem() Item {
	for {
		select {
		case item := <-l.items:
			return item
		default:
			if l.state == nil {
				return Item{EOF, nil}
			}
			l.itemStart = l.pos
			l.state = l.state(l)
		}
	}
}

func (l *Lexer) peek() (byte, error) {
	n, err := l.next()
	l.backup()
	return n, err
}

func (l *Lexer) accept(valid []byte) bool {
	if next, err := l.next(); err == nil && bytes.IndexByte(valid, next) >= 0 {
		return true
	}
	l.backup()
	return false
}

func (l *Lexer) acceptFn(validPredicate func(byte) bool) bool {
	if next, err := l.next(); err == nil && validPredicate(next) {
		return true
	}
	l.backup()
	return false
}

func (l *Lexer) acceptRunFn(validPredicate func(byte) bool) {
	for {
		next, err := l.next()
		if err != nil || !validPredicate(next) {
			break
		}
	}
	l.backup()
}
