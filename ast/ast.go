// ast/ast.go
package ast

import (
    "bytes"
    "fmt"
)

// Node is the common interface for all AST nodes.
type Node interface {
    TokenLiteral() string
    String() string
}

// Statement represents statements (e.g., assignments, return, etc.).
type Statement interface {
    Node
    statementNode()
}

// Expression represents expressions (identifiers, literals, infix, …).
type Expression interface {
    Node
    expressionNode()
}

// Program is the root node of every AST we produce.
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

// AssignmentStatement ::= Identifier "=" Expression
type AssignmentStatement struct {
    Name  *Identifier
    Value Expression
}

func (as *AssignmentStatement) statementNode()       {}
func (as *AssignmentStatement) TokenLiteral() string { return as.Name.Value }
func (as *AssignmentStatement) String() string {
    return fmt.Sprintf("%s = %s", as.Name.String(), as.Value.String())
}

// Identifier represents a variable name.
type Identifier struct {
    Value string
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Value }
func (i *Identifier) String() string       { return i.Value }

// IntegerLiteral is a numeric literal, e.g. "3"
type IntegerLiteral struct {
    Value int64
}

func (il *IntegerLiteral) expressionNode()      {}
func (il *IntegerLiteral) TokenLiteral() string { return fmt.Sprintf("%d", il.Value) }
func (il *IntegerLiteral) String() string       { return il.TokenLiteral() }
