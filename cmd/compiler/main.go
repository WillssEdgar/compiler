package main

import (
	"fmt"
	"compiler/lexer"
	"compiler/parser"
)

func main() {
	source := "varName = 12345\nif x == 10"
	l := lexer.New(source)
	p := parser.New(l)
	p.ParseProgram()
	fmt.Println("Parsing complete")
}
