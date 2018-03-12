package ansi

import (
	"bytes"
	"errors"
	"strconv"
)

var (
	// ErrBadControlSequence is an error definition for malformed control
	// sequences
	ErrBadControlSequence = errors.New("malformed control sequence")
	// ControlSequenceIntroducer holds the ANSI control sequence introducer '['
	ControlSequenceIntroducer = byte('[')
	// SelectGraphicRendition holds the ANSI control sequence command 'm'
	SelectGraphicRendition = byte('m')
	// CursorPosition holds the ANSI control sequence command 'H'
	CursorPosition = byte('H')
	// EraseInDisplay holds the ANSI control sequence command 'J'
	EraseInDisplay = byte('J')
	// EraseInLine holds the ANSI control sequence command 'K'
	EraseInLine = byte('K')
	// CursorUp holds the ANSI control sequence command 'A'
	CursorUp = byte('A')
	// CursorDown holds the ANSI control sequence command 'B'
	CursorDown = byte('B')
	// CursorForward holds the ANSI control sequence command 'C'
	CursorForward = byte('C')
	// CursorBack holds the ANSI control sequence command 'D'
	CursorBack = byte('D')
)

// SequenceData is a struct that describes a control sequence
type SequenceData struct {
	Prefix  byte
	Params  []int
	Inters  []byte
	Command byte
}

// ParseControlSequence takes a slice of bytes as input and returns SequenceData and an error
func ParseControlSequence(v []byte) (*SequenceData, error) {
	// Immediatly reject any malformed control sequence: it must start with the
	// escape character, and contain at least one prefix and command byte.
	if len(v) < 3 || v[0] != EscapeCode {
		return nil, ErrBadControlSequence
	}

	// Everything between the prefix and the command bytes are arguments: we
	// need to determine where parameters end and intermediate char begin.
	var i int
	end := len(v) - 1
	for i = end - 1; IsInterChar(v[i]); i-- {
	}

	// Value of i marks the separation between (semicolon-separated) parameters
	// and intermediate bytes. One catch: when no parameters are specified, we
	// want to have [][]byte{} rather than [][]byte{[]byte{}}.
	params := []int{}
	if i >= 2 {
		paramBytes := bytes.Split(v[2:i+1], []byte{';'})
		for _, param := range paramBytes {
			val, _ := strconv.Atoi(string(param))
			params = append(params, val)
		}
	}

	return &SequenceData{
		Prefix:  v[1],
		Params:  params,
		Inters:  v[i+1 : end],
		Command: v[end],
	}, nil
}
