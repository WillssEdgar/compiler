package main

import (
	"compiler/lexer"
	"compiler/parser"
	"fmt"
)

func main() {
	source := "varName = 3"
	l := lexer.New(source)
	p := parser.New(l)

	prog := p.ParseProgram()
	fmt.Printf("prog: %v\n", prog)
}
