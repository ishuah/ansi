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

// Init initializes a Lexer with new input
func (l *Lexer) Init(input []byte) {
	l.input = input
	l.items = make(chan Item, 2)
	l.state = lexBytes
}

func (l *Lexer) Backup() {
	l.pos--
}

func (l *Lexer) Cancel(revert stateFn) stateFn {
	l.pos = l.itemStart
	return revert
}

func (l *Lexer) Emit(t ItemType) {
	l.items <- Item{T: t, Value: l.input[l.start:l.pos]}
	l.start = l.pos
}

func (l *Lexer) Next() (byte, error) {
	if l.pos >= len(l.input) {
		return 0, io.EOF
	}
	b := l.input[l.pos]
	l.pos++
	return b, nil
}

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

func (l *Lexer) Peek() (byte, error) {
	n, err := l.Next()
	l.Backup()
	return n, err
}

func (l *Lexer) Accept(valid []byte) bool {
	if next, err := l.Next(); err == nil && bytes.IndexByte(valid, next) >= 0 {
		return true
	}
	l.Backup()
	return false
}

func (l *Lexer) AcceptFn(validPredicate func(byte) bool) bool {
	if next, err := l.Next(); err == nil && validPredicate(next) {
		return true
	}
	l.Backup()
	return false
}

func (l *Lexer) AcceptRun(valid []byte) {
	for {
		next, err := l.Next()
		if err != nil || bytes.IndexByte(valid, next) < 0 {
			break
		}
	}
	l.Backup()
}

func (l *Lexer) AcceptRunFn(validPredicate func(byte) bool) {
	for {
		next, err := l.Next()
		if err != nil || !validPredicate(next) {
			break
		}
	}
	l.Backup()
}
