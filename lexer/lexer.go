package lexer

import (
	"unicode/utf8"
)

// Lexer represents the lexical analyzer.
type Lexer struct {
	Input        string // the input string
	Position     int    // start of current rune
	ReadPosition int    // next reading position in bytes
	Ch           rune   // current rune
	Line         int    // current line number
	Column       int    // current column number
}

// New returns a new Lexer instance for a given input.
func New(input string) *Lexer {
	l := &Lexer{
		Input:        input,
		Position:     0,
		ReadPosition: 0,
		Line:         1,
		Column:       0,
	}
	l.readChar()
	return l
}

// readChar advances the lexer one rune at a time.
func (l *Lexer) readChar() {
	if l.ReadPosition >= len(l.Input) {
		l.Ch = 0 // EOF
		return
	}
	var width int
	l.Ch, width = utf8.DecodeRuneInString(l.Input[l.ReadPosition:])
	l.Position = l.ReadPosition
	l.ReadPosition += width

	if l.Ch == '\n' {
		l.Line++
		l.Column = 0
	} else {
		l.Column++
	}
}

// peekChar returns the next rune without advancing.
func (l *Lexer) peekChar() rune {
	if l.ReadPosition >= len(l.Input) {
		return 0
	}
	ch, _ := utf8.DecodeRuneInString(l.Input[l.ReadPosition:])
	return ch
}

// isLetter checks if a rune is a letter or underscore.
func isLetter(ch rune) bool {
	return ('a' <= ch && ch <= 'z') ||
		('A' <= ch && ch <= 'Z') ||
		ch == '_'
}

// isDigit checks if a rune is a digit.
func isDigit(ch rune) bool {
	return '0' <= ch && ch <= '9'
}

// readIdentifier reads a full identifier.
func (l *Lexer) readIdentifier() string {
	start := l.Position
	for isLetter(l.Ch) || isDigit(l.Ch) {
		l.readChar()
	}
	return l.Input[start:l.Position]
}

// readNumber reads a number token.
func (l *Lexer) readNumber() string {
	start := l.Position
	for isDigit(l.Ch) {
		l.readChar()
	}
	return l.Input[start:l.Position]
}

// lookupIdentifier classifies an identifier.
func lookupIdentifier(ident string) TokenType {
	keywords := map[string]TokenType{
		"if":     TokenKeyword,
		"else":   TokenKeyword,
		"for":    TokenKeyword,
		"func":   TokenKeyword,
		"return": TokenKeyword,
	}
	if typ, ok := keywords[ident]; ok {
		return typ
	}
	return TokenIdentifier
}

// NextToken returns the next token from input.
func (l *Lexer) NextToken() Token {
	var tok Token
	// Skip whitespace
	for l.Ch == ' ' || l.Ch == '\t' || l.Ch == '\r' || l.Ch == '\n' {
		l.readChar()
	}

	tok.Line = l.Line
	tok.Column = l.Column

	if isLetter(l.Ch) {
		lexeme := l.readIdentifier()
		tok.Lexeme = lexeme
		tok.Type = lookupIdentifier(lexeme)
		return tok
	} else if isDigit(l.Ch) {
		tok.Lexeme = l.readNumber()
		tok.Type = TokenNumber
		return tok
	}

	// Process operators and single character tokens.
	switch l.Ch {
	case '=':
		if l.peekChar() == '=' {
			ch := l.Ch
			l.readChar()
			tok.Lexeme = string(ch) + string(l.Ch)
			tok.Type = TokenOperator
		} else {
			tok.Type = TokenOperator
			tok.Lexeme = string(l.Ch)
		}
	case 0:
		tok.Type = TokenEOF
		tok.Lexeme = ""
	default:
		tok.Type = TokenIdentifier
		tok.Lexeme = string(l.Ch)
	}
	l.readChar()
	return tok
}
