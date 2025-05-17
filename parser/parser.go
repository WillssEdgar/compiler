package parser

import (
	"compiler/ast"
	"compiler/lexer"
	"compiler/token"
	"encoding/json"
	"fmt"
	"strconv"
)

const (
	_ int = iota
	LOWEST
	SUM     // + or -
	PRODUCT // * or /
	PREFIX  // -X or !X
	// CALL  // func(X)
)

var precedences = map[string]int{
	"+": SUM,
	"-": SUM,
	"*": PRODUCT,
	"/": PRODUCT,
}

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

type Parser struct {
	L         *lexer.Lexer `json:"l"`
	CurToken  token.Token  `json:"curToken"`
	PeekToken token.Token  `json:"peekToken"`

	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		L:              l,
		prefixParseFns: make(map[token.TokenType]prefixParseFn),
		infixParseFns:  make(map[token.TokenType]infixParseFn),
	}

	p.registerPrefix(token.TokenIdentifier, p.parseIdentifier)
	p.registerPrefix(token.TokenNumber, p.parseIntegerLiteral)
	p.registerPrefix(token.TokenOperator, p.parsePrefixExpression)

	p.registerInfix(token.TokenOperator, p.parseInfixExpression)
	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) registerPrefix(tt token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tt] = fn
}

func (p *Parser) registerInfix(tt token.TokenType, fn infixParseFn) {
	p.infixParseFns[tt] = fn
}

func (p *Parser) nextToken() {
	p.CurToken = p.PeekToken
	p.PeekToken = p.L.NextToken()
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	fmt.Printf("lexeme: %s", p.CurToken.Lexeme)
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
	case token.TokenKeyword:
		switch p.CurToken.Lexeme {
		case "let":
			return p.parseLetStatement()
		case "return":
			println("i am in return")
			return p.parseReturnStatement()
		case "func":
			return p.parseFunctionDeclaration()
		}
	case token.TokenRBrace:
		print("help i am here")
		return nil
	}
	return nil
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	p.nextToken()
	name := &ast.Identifier{Value: p.CurToken.Lexeme}

	p.nextToken()
	p.nextToken()

	value := p.parseExpression(LOWEST)

	if p.PeekToken.Lexeme == ";" {
		p.nextToken()
	}

	return &ast.LetStatement{
		Assignment: ast.AssignmentStatement{
			Name:  name,
			Value: value,
		},
	}
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.CurToken}
	p.nextToken() // move past 'return'
	stmt.ReturnValue = p.parseExpression(LOWEST)
	if p.PeekToken.Lexeme == ";" {
		p.nextToken()
	}
	return stmt
}

func (p *Parser) parseFunctionDeclaration() *ast.FunctionalLiteral {
	fl := &ast.FunctionalLiteral{Token: p.CurToken}

	p.nextToken()
	fl.ReturnType = p.CurToken.Lexeme

	p.nextToken()
	fl.FunctionName = &ast.Identifier{Value: p.CurToken.Lexeme}

	if !p.expectPeek(token.TokenLParen) {
		return nil
	}
	fl.Parameters = p.parseFunctionParameters()

	if !p.expectPeek(token.TokenLBrace) {
		return nil
	}
	fl.Body = p.parseBlockStatement()

	return fl
}

func (p *Parser) parseFunctionParameters() []*ast.Identifier {
	identifiers := []*ast.Identifier{}
	p.nextToken()
	if p.CurToken.Type == token.TokenRParen {
		return identifiers
	}

	for {
		ident := &ast.Identifier{Value: p.CurToken.Lexeme}
		identifiers = append(identifiers, ident)
		if p.PeekToken.Type != token.TokenComma {
			break
		}
		p.nextToken() // ,
		p.nextToken() // next identifier
	}

	if !p.expectPeek(token.TokenRParen) {
		return nil
	}
	return identifiers
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	bs := &ast.BlockStatement{Token: p.CurToken}
	bs.Statements = []ast.Statement{}

	p.nextToken() // consume '{'

	for p.CurToken.Type != token.TokenRBrace && p.CurToken.Type != token.TokenEOF {
		stmt := p.parseStatement()
		if stmt != nil {
			bs.Statements = append(bs.Statements, stmt)
		}
		p.nextToken()
	}
	return bs
}

// helper to check & advance
func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.PeekToken.Type == t {
		p.nextToken()
		return true
	}
	return false
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.CurToken.Type]
	if prefix == nil {
		return nil
	}
	leftExp := prefix()

	for p.PeekToken.Lexeme != ";" && precedence < p.peekPrecedence() {
		// p.PeekToken.PrintToken()
		infix := p.infixParseFns[p.PeekToken.Type]
		if infix == nil {
			return leftExp
		}
		p.nextToken()
		leftExp = infix(leftExp)
	}
	return leftExp
}

func (p *Parser) peekPrecedence() int {
	if prec, ok := precedences[p.PeekToken.Lexeme]; ok {
		return prec
	}
	return LOWEST
}

func (p *Parser) curPrecedence() int {
	if prec, ok := precedences[p.CurToken.Lexeme]; ok {
		return prec
	}
	return LOWEST
}

func (p *Parser) parseAssignmentStatement() *ast.AssignmentStatement {
	name := &ast.Identifier{Value: p.CurToken.Lexeme}

	p.nextToken()
	p.nextToken()

	value := p.parseExpression(LOWEST)

	if p.PeekToken.Lexeme == ";" {
		p.nextToken()
	}

	return &ast.AssignmentStatement{
		Name:  name,
		Value: value,
	}
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Value: p.CurToken.Lexeme}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{}
	val, err := strconv.ParseInt(p.CurToken.Lexeme, 10, 64)
	if err != nil {
		return nil
	}
	lit.Value = val
	return lit
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	println("I am here")
	expr := &ast.PrefixExpression{
		Operator: p.CurToken.Lexeme,
	}
	p.nextToken()
	expr.Right = p.parseExpression(PREFIX)
	return expr
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	exp := &ast.InfixExpression{
		Left:     left,
		Operator: p.CurToken.Lexeme,
	}
	prec := p.curPrecedence()
	p.nextToken()
	exp.Right = p.parseExpression(prec)
	return exp
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.PeekToken.Type == t
}

func (p *Parser) PrintParser() {
	jsonP, _ := json.MarshalIndent(p, " ", "	")
	fmt.Printf("\nParser: %s\n", jsonP)
}
