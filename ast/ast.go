package ast

import (
	"bytes"
	"compiler/token"
	"fmt"
)

// Node is the interface for all AST nodes.
type Node interface {
	TokenLiteral() string
	String() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}
	return ""
}

func (p *Program) String() string {
	var out bytes.Buffer
	for _, s := range p.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}

type LetStatement struct {
	Assignment AssignmentStatement
}

func (ls *LetStatement) statementNode()       {}
func (ls *LetStatement) TokenLiteral() string { return "let" }
func (ls *LetStatement) String() string {
	return fmt.Sprintf("let %s;", ls.Assignment.String())
}

type AssignmentStatement struct {
	Name  *Identifier
	Value Expression
}

func (as *AssignmentStatement) statementNode()       {}
func (as *AssignmentStatement) TokenLiteral() string { return as.Name.Value }
func (as *AssignmentStatement) String() string {
	return fmt.Sprintf("%s = %s", as.Name.String(), as.Value.String())
}

type Identifier struct {
	Token token.Token
	Value string
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Value }
func (i *Identifier) String() string       { return i.Value }

// IntegerLiteral is a number.
type IntegerLiteral struct {
	Value int64
}

func (il *IntegerLiteral) expressionNode()      {}
func (il *IntegerLiteral) TokenLiteral() string { return fmt.Sprintf("%d", il.Value) }
func (il *IntegerLiteral) String() string       { return il.TokenLiteral() }

// PrefixExpression e.g. -x
type PrefixExpression struct {
	Operator string     // e.g. "-"
	Right    Expression // the sub‐expression
}

func (pe *PrefixExpression) expressionNode()      {}
func (pe *PrefixExpression) TokenLiteral() string { return pe.Operator }
func (pe *PrefixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")
	return out.String()
}

// InfixExpression e.g. x + y
type InfixExpression struct {
	Left     Expression // left‐hand side
	Operator string     // e.g. "+"
	Right    Expression // right‐hand side
}

func (ie *InfixExpression) expressionNode()      {}
func (ie *InfixExpression) TokenLiteral() string { return ie.Operator }
func (ie *InfixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString(" " + ie.Operator + " ")
	out.WriteString(ie.Right.String())
	out.WriteString(")")
	return out.String()
}

type FunctionStatement struct {
	Literal *FunctionalLiteral
}

func (fs *FunctionStatement) statementNode()       {}
func (fs *FunctionStatement) TokenLiteral() string { return fs.Literal.Token.Lexeme }
func (fs *FunctionStatement) String() string       { return fs.Literal.String() }

type FunctionalLiteral struct {
	Token        token.Token
	ReturnType   string
	FunctionName *Identifier
	Parameters   []*Identifier
	Body         *BlockStatement
}

func (fl *FunctionalLiteral) statementNode()       {}
func (fl *FunctionalLiteral) TokenLiteral() string { return fl.Token.Lexeme }
func (fl *FunctionalLiteral) String() string {
	var out bytes.Buffer
	out.WriteString("func ")
	out.WriteString(fl.ReturnType)
	out.WriteString(" ")
	out.WriteString(fl.FunctionName.String())
	out.WriteString("(")
	for i, p := range fl.Parameters {
		if i > 0 {
			out.WriteString(", ")
		}
		out.WriteString(p.String())
	}
	out.WriteString(") ")
	out.WriteString(fl.Body.String())
	return out.String()
}

type BlockStatement struct {
	Token      token.Token
	Statements []Statement
}

func (bs *BlockStatement) expressionNode()      {}
func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Lexeme }
func (bs *BlockStatement) String() string {
	var out bytes.Buffer
	out.WriteString("{\n")
	for _, s := range bs.Statements {
		out.WriteString(s.String())
		out.WriteString("\n")
	}
	out.WriteString("}")
	return out.String()
}

type ReturnStatement struct {
	Token       token.Token
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode()       {}
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Lexeme }
func (rs *ReturnStatement) String() string {
	return fmt.Sprintf("return %s;", rs.ReturnValue.String())
}

type CallExpression struct {
	Token     token.Token // the '(' token
	Function  Expression  // Identifier or FunctionLiteral
	Arguments []Expression
}

func (ce *CallExpression) expressionNode()      {}
func (ce *CallExpression) TokenLiteral() string { return ce.Token.Lexeme }

func (ce *CallExpression) String() string {
	var out bytes.Buffer

	// function name or literal
	out.WriteString(ce.Function.String())
	out.WriteString("(")

	// each argument, comma-separated
	for i, arg := range ce.Arguments {
		if i > 0 {
			out.WriteString(", ")
		}
		out.WriteString(arg.String())
	}

	out.WriteString(")")
	return out.String()
}
