package token

import (
	"encoding/json"
	"fmt"
)

type TokenType int

const (
	TokenEOF TokenType = iota
	TokenIdentifier
	TokenNumber
	TokenOperator
	TokenKeyword
)

type Token struct {
	Type   TokenType `json:"tokentype"`
	Lexeme string    `json:"lexeme"`
	Line   int       `json:"line"`
	Column int       `json:"column"`
}

type Keyword int

const (
	KeywordIf Keyword = iota
	KeywordElse
	KeywordFor
	KeywordFunc
	KeywordReturn
)

func (t *Token) PrintToken() {
	jsonPrint, _ := json.MarshalIndent(t, " ", "	")
	fmt.Printf("Token: %s", jsonPrint)
}
