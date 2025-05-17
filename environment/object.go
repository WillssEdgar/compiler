package environment

import (
	"compiler/ast"
	"fmt"
)

type ObjectType string

const (
	INTEGER_OBJ      = "INTEGER"
	NULL_OBJ         = "NULL"
	ERROR_OBJ        = "ERROR"
	FUNCTION_OBJ     = "FUNCTION"
	RETURN_VALUE_OBJ = "RETURN_VALUE"
)

type Object interface {
	Type() ObjectType
	Inspect() string
}

type Integer struct {
	Value int64
}

func (i *Integer) Type() ObjectType { return INTEGER_OBJ }
func (i *Integer) Inspect() string  { return fmt.Sprintf("%d", i.Value) }

type Null struct{}

func (n *Null) Type() ObjectType { return NULL_OBJ }
func (n *Null) Inspect() string  { return "null" }

type Error struct {
	Message string
}

func (e *Error) Type() ObjectType { return ERROR_OBJ }
func (e *Error) Inspect() string  { return "ERROR: " + e.Message }

type Function struct {
	Literal *ast.FunctionalLiteral
	Env     *Environment
}

func (f *Function) Type() ObjectType {
	return FUNCTION_OBJ
}
func (f *Function) Inspect() string {
	// you could print its signature & body:
	return fmt.Sprintf("%s %s %s %s",
		f.Literal.Token.Lexeme,       // "func"
		f.Literal.ReturnType,         // e.g. "Integer"
		f.Literal.FunctionName.Value, // function name
		f.Literal.Body.String(),      // the "{ … }"
	)
}

type ReturnValue struct {
	Value Object
}

func (rv *ReturnValue) Type() ObjectType {
	return RETURN_VALUE_OBJ
}

func (rv *ReturnValue) Inspect() string {
	return rv.Value.Inspect()
}
