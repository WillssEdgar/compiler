package parser

import (
	"fmt"
	"compiler/lexer" // adjust the module path accordingly
)

type Parser struct {
	l         *lexer.Lexer
	curToken  lexer.Token
	peekToken lexer.Token
}

// New creates a new Parser instance.
func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l: l,
	}
	// Read two tokens to initialize current and peek.
	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

// ParseProgram starts parsing and returns results.
// Later you'll populate this with your AST.
func (p *Parser) ParseProgram() {
	for p.curToken.Type != lexer.TokenEOF {
		fmt.Printf("Parsing token: %+v\n", p.curToken)
		p.nextToken()
	}
}
