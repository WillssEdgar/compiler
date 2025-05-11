package main

import (
	"compiler/environment"
	"compiler/evaluator"
	"compiler/lexer"
	"compiler/parser"
	"encoding/json"
	"fmt"
)

func main() {
	source := "let varName = 3 + 3;"
	l := lexer.New(source)
	p := parser.New(l)

	prog := p.ParseProgram()
	fmt.Printf("prog: %#v\n", prog)
	parsedprogram, _ := json.MarshalIndent(prog, " ", "	")
	fmt.Printf("program: %s", parsedprogram)
	env := environment.NewEnvironment()
	result := evaluator.Eval(prog, env)

	fmt.Printf("\nResult: %s\n", result.Inspect())

	if val, ok := env.Get("varName"); ok {
		fmt.Printf("varName = %s\n", val.Inspect())
	}
}
