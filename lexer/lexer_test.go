package lexer

import (
	"testing"
)

func TestIsLetter(t *testing.T) {
	tests := []struct {
		input rune
		want  bool
	}{
		{'a', true},
		{'z', true},
		{'A', true},
		{'Z', true},
		{'_', true},
		{'1', false},
		{'$', false},
		{' ', false},
		{'%', false},
	}

	for _, tt := range tests {
		got := isLetter(tt.input)
		if got != tt.want {
			t.Errorf("isLetter(%q) = %v; want %v", tt.input, got, tt.want)
		}
	}
}

func TestIsLetterSimple(t *testing.T) {
	expected := true

	got := isLetter('a')

	if got != expected {
		t.Errorf("isLetter(a) = %v; want %v", got, expected)
	}
}

func TestNewLexer(t *testing.T) {
	input := "varName = 3"
	l := New(input)

	expected := &Lexer{
		Input:        input,
		Position:     0,
		ReadPosition: 1,
		Line:         1,
		Column:       1,
		Ch:           'v',
	}

	tests := []struct {
		name string
		got  any
		want any
	}{
		{"Input", l.Input, expected.Input},
		{"Position", l.Position, expected.Position},
		{"ReadPosition", l.ReadPosition, expected.ReadPosition},
		{"Line", l.Line, expected.Line},
		{"Column", l.Column, expected.Column},
		{"Ch", l.Ch, expected.Ch},
	}

	for _, tt := range tests {
		if tt.got != tt.want {
			t.Errorf("%s mismatch: got=%#v, want=%#v", tt.name, tt.got, tt.want)
		}
	}
}

func TestReadChar(t *testing.T) {
	input := "varName = 3"

	l := New(input)
	l.readChar()

	tests := []struct {
		name string
		got  any
		want any
	}{
		{"Position", l.Position, 1},
		{"ReadPosition", l.ReadPosition, 2},
		{"Line", l.Line, 1},
		{"Column", l.Column, 2},
		{"Ch", l.Ch, 'a'},
	}
	for _, tt := range tests {
		if tt.got != tt.want {
			t.Errorf("%s mismatch: got=%#v, want=%#v", tt.name, tt.got, tt.want)
		}
	}
}

func TestPeekChar(t *testing.T) {
	input := "varName = 3"
	l := New(input)
	nextChar := l.peekChar()

	if nextChar != 'a' {
		t.Errorf("Next Char is %c not 'a'", nextChar)
	}
}

func TestReadIdentifier(t *testing.T) {
	input := "varName = 3"
	l := New(input)
	output := l.readIdentifier()

	if output != "varName" {
		t.Errorf("Identifier expected: varName, Identifier recieved: %s", output)
	}
}

func TestReadNumber(t *testing.T) {
	input := "3932"
	l := New(input)
	output := l.readNumber()

	if output != "3932" {
		t.Errorf("Number expected: 3932, Number recieved: %s", output)
	}
}

func TestLookupIdentifier(t *testing.T) {
	input := "if"
	output := lookupIdentifier(input)

	if output != 4 {
		t.Errorf("TokenType expected: 4, Number recieved: %d", output)
	}
}


func TestNextToken(t *testing.T) {
	input := "varName = 3"
	l := New(input)
	output := l.NextToken()

	if output.Lexeme != "varName" {
		t.Errorf("Token expected: 'varName', Token recieved: %s", output.Lexeme)
	}
}
