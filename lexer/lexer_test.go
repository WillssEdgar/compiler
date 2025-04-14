package lexer

import "testing"

func TestNextToken(t *testing.T) {
	input := "varName = 12345\nif x == 10"

	tests := []struct {
		expectedType   TokenType
		expectedLexeme string
	}{
		{TokenIdentifier, "varName"},
		{TokenOperator, "="},
		{TokenNumber, "12345"},
		{TokenKeyword, "if"},
		{TokenIdentifier, "x"},
		{TokenOperator, "=="},
		{TokenNumber, "10"},
		{TokenEOF, ""},
	}

	// Create a new lexer with the input.
	l := New(input)

	// Loop through each expected token and compare with lexer output.
	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - token type wrong. expected=%d, got=%d", i, tt.expectedType, tok.Type)
		}

		if tok.Lexeme != tt.expectedLexeme {
			t.Fatalf("tests[%d] - token lexeme wrong. expected=%q, got=%q", i, tt.expectedLexeme, tok.Lexeme)
		}
	}
}

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
