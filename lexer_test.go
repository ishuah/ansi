package ansi

import (
	"bytes"
	"testing"
)

func TestLexer(t *testing.T) {
	l := Lexer{}
	testString := []byte("test string")

	l.Init(testString)

	item := l.NextItem()
	if bytes.Compare(item.Value, testString) != 0 {
		t.Fatalf("Wrong character sequence: got %q, expected %q", item.Value, testString)
	}
	if item.T.String() != ItemTypeName[RawBytes] {
		t.Fatalf("Wrong item type: got %s, expected %s", item.T.String(), ItemTypeName[RawBytes])
	}

	item = l.NextItem()
	if item.T.String() != ItemTypeName[EOF] {
		t.Fatalf("Wrong item type: got %s, expected %s", item.T.String(), ItemTypeName[EOF])
	}

	testString = []byte("\x1b[J")
	l = Lexer{}
	l.Init(testString)
	item = l.NextItem()

	if item.T.String() != ItemTypeName[ControlSequence] {
		t.Fatalf("Wrong item type: got %s, expected %s", item.T.String(), ItemTypeName[ControlSequence])
	}

	testString = []byte("\x1b[?")
	l = Lexer{}
	l.Init(testString)

	item = l.NextItem()

	if bytes.Compare(item.Value, testString) != 0 {
		t.Fatalf("Wrong item type: got %q, expected %q", item.Value, testString)
	}
}
