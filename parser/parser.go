package parser

import (
	"compiler/ast"   // adjust the module path accordingly
	"compiler/lexer" // adjust the module path accordingly
	"strconv"
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

func (p *Parser) ParseProgram() *ast.Program {
    program := &ast.Program{}
    for p.curToken.Type != lexer.TokenEOF {
        stmt := p.parseStatement()
        if stmt != nil {
            program.Statements = append(program.Statements, stmt)
        }
        p.nextToken()
    }
    return program
}

func (p *Parser) parseStatement() ast.Statement {
    // if current token is an identifier and next is '=', parse assignment
    if p.curToken.Type == lexer.TokenIdentifier && p.peekToken.Lexeme == "=" {
        return p.parseAssignmentStatement()
    }
    // … other statements …
    return nil
}

func (p *Parser) parseAssignmentStatement() *ast.AssignmentStatement {
    stmt := &ast.AssignmentStatement{
        Name: &ast.Identifier{Value: p.curToken.Lexeme},
    }
    // consume identifier
    p.nextToken()
    // consume '='
    p.nextToken()
    // parse the right-hand expression (just a number literal for now)
    if p.curToken.Type == lexer.TokenNumber {
        val, _ := strconv.ParseInt(p.curToken.Lexeme, 10, 64)
        stmt.Value = &ast.IntegerLiteral{Value: val}
    }
    return stmt
}
