package parser

import (
	"compiler/ast"
	"compiler/lexer"
	"compiler/token"
	"encoding/json"
	"fmt"
	"strconv"
)

type Parser struct {
	L         *lexer.Lexer `json:"l"`
	CurToken  token.Token  `json:"curToken"`
	PeekToken token.Token  `json:"peekToken"`
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		L: l,
	}
	p.PrintParser()
	p.nextToken()
	p.PrintParser()
	p.nextToken()
	p.PrintParser()
	return p
}

func (p *Parser) nextToken() {
	p.CurToken = p.PeekToken
	p.PeekToken = p.L.NextToken()
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	for p.CurToken.Type != token.TokenEOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}
	return program
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.CurToken.Type {
	case token.TokenIdentifier:
		if p.PeekToken.Lexeme == "=" {
			return p.parseAssignmentStatement()
		}
	}
	return nil
}

func (p *Parser) parseAssignmentStatement() *ast.AssignmentStatement {
	stmt := &ast.AssignmentStatement{
		Name: &ast.Identifier{Value: p.CurToken.Lexeme},
	}
	p.nextToken()
	p.nextToken()
	p.PrintParser()
	if p.CurToken.Type == token.TokenNumber {
		val, _ := strconv.ParseInt(p.CurToken.Lexeme, 10, 64)
		stmt.Value = &ast.IntegerLiteral{Value: val}
	}
	return stmt
}

func (p *Parser) PrintParser() {
	jsonP, _ := json.MarshalIndent(p, " ", "	")
	fmt.Printf("\nParser: %s\n", jsonP)
}
