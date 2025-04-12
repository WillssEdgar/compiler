package main

import (
	"encoding/json"
	"fmt"
	"unicode/utf8"
)

type TokenType int

const (
	TokenEOF TokenType = iota
	TokenIdentifier
	TokenNumber
	TokenOperator
	// ... add more token types as needed.
)

type Token struct {
	Type    TokenType `json:"tokentype"`
	Lexeme  string `json:"lexeme"`
	Line    int `json:"line"`
	Column  int `json:"column"`
}

type Lexer struct {
	Input        string `json:"input"`
	Position     int `json:"position"`
	ReadPosition int `json:"readPosition"`
	Ch           rune `json:"ch"`
	Line         int `json:"line"`
	Column       int `json:"column"`
}

func NewLexer(input string) *Lexer {
	l := &Lexer{Input: input, Line: 1, Column: 0}

	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.Position >= len(l.Input) {
		l.Ch = 0
		return
	}
	l.Ch, _ = utf8.DecodeRuneInString(l.Input[l.Position:])
	l.Position = l.ReadPosition
	l.ReadPosition += utf8.RuneLen(l.Ch)

	if l.Ch == '\n' {
		l.Line++
		l.Column = 0
	} else {
		l.Column++
	}
}

func (l *Lexer) NextToken() Token {

	json, err := json.Marshal(l)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("L: %s\n", json)
	var tok Token

	for l.Ch == ' ' || l.Ch == '\t' || l.Ch == '\r' || l.Ch == '\n' {
		l.readChar()
	}

	tok.Line = l.Line
	tok.Column = l.Column

	switch l.Ch {
	case '=':
		tok = Token{Type: TokenOperator, Lexeme: string(l.Ch), Line: l.Line, Column: l.Column}
	case 0:
		tok = Token{Type: TokenEOF, Lexeme: "", Line: l.Line, Column: l.Column}
	default:
		tok = Token{Type: TokenIdentifier, Lexeme: string(l.Ch), Line: l.Line, Column: l.Column}
	}

	l.readChar()
	return tok
}

func main() {
	source := "a = 3\nb = a + 4"
	lexer := NewLexer(source)

	t := []Token{}

	for {
		token := lexer.NextToken()
		if token.Type == TokenEOF {
			break
		}
		fmt.Printf("Token: %+v\n", token)
		t = append(t, token)
	}

	jsont, _ := json.Marshal(t)
	fmt.Printf("Token List: %s", string(jsont))
}
