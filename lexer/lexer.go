package lexer

import (
	"compiler/token"
	"fmt"
	"unicode/utf8"
)

type Lexer struct {
	Input        string `json:"Input"`
	Position     int    `json:"Position"`
	ReadPosition int    `json:"ReadPosition"`
	Ch           rune   `json:"Ch"`
	Line         int    `json:"Line"`
	Column       int    `json:"Column"`
}

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

func (l *Lexer) readChar() {
	if l.ReadPosition >= len(l.Input) {
		l.Position = l.ReadPosition
		l.ReadPosition++
		l.Ch = 0
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

func (l *Lexer) peekChar() rune {
	if l.ReadPosition >= len(l.Input) {
		return 0
	}
	ch, _ := utf8.DecodeRuneInString(l.Input[l.ReadPosition:])
	return ch
}

func isLetter(ch rune) bool {
	return ('a' <= ch && ch <= 'z') ||
		('A' <= ch && ch <= 'Z') ||
		ch == '_'
}

func isDigit(ch rune) bool {
	return '0' <= ch && ch <= '9'
}

func (l *Lexer) readIdentifier() string {
	start := l.Position
	for isLetter(l.Ch) || isDigit(l.Ch) {
		l.readChar()
	}
	return l.Input[start:l.Position]
}

func (l *Lexer) readNumber() string {
	start := l.Position
	for isDigit(l.Ch) {
		l.readChar()
	}
	return l.Input[start:l.Position]
}

func lookupIdentifier(ident string) token.TokenType {
	keywords := map[string]token.TokenType{
		"let":     token.TokenKeyword,
		"if":      token.TokenKeyword,
		"else":    token.TokenKeyword,
		"for":     token.TokenKeyword,
		"func":    token.TokenKeyword,
		"return":  token.TokenKeyword,
		"Integer": token.TokenKeyword,
		"String":  token.TokenKeyword,
	}
	if typ, ok := keywords[ident]; ok {
		return typ
	}
	return token.TokenIdentifier
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token
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
		tok.Type = token.TokenNumber
		return tok
	}

	switch l.Ch {
	case '=':
		if l.peekChar() == '=' {
			ch := l.Ch
			l.readChar()
			tok.Lexeme = string(ch) + string(l.Ch)
			tok.Type = token.TokenOperator
		} else {
			tok.Type = token.TokenOperator
			tok.Lexeme = string(l.Ch)
		}
	case '+':
		tok.Type = token.TokenOperator
		tok.Lexeme = string(l.Ch)
	case '{':
		tok.Type = token.TokenLBrace
		tok.Lexeme = "{"
	case '}':
		tok.Type = token.TokenRBrace
		tok.Lexeme = "}"
	case '(':
		tok.Type = token.TokenLParen
		tok.Lexeme = "("
	case ')':
		tok.Type = token.TokenRParen
		tok.Lexeme = ")"
	case ';':
		tok.Type = token.TokenSemicolon
		tok.Lexeme = string(l.Ch)
	case ',':
		tok.Type = token.TokenComma
		tok.Lexeme = string(l.Ch)
	case 0:
		tok.Type = token.TokenEOF
		tok.Lexeme = ""
	default:
		tok.Type = token.TokenIdentifier
		tok.Lexeme = string(l.Ch)
	}
	l.readChar()
	return tok
}

func (l *Lexer) PrintLexer() {
	fmt.Printf("{\n\tInput: %s\n\tPosition: %d\n\tReadPosition: %d\n\tCh: %c\n\tLine: %d\n\tColumn: %d\n}\n",
		l.Input,
		l.Position,
		l.ReadPosition,
		l.Ch,
		l.Line,
		l.Column,
	)
}
