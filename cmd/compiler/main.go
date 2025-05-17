package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"compiler/environment"
	"compiler/evaluator"
	"compiler/lexer"
	"compiler/parser"
)

func main() {
	// 1) Locate your source file
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("could not get working directory: %v", err)
	}
	path := filepath.Join(wd, "files", "main.blue")

	// 2) Read it all in one shot
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("failed to read %s: %v", path, err)
	}
	source := string(data)

	// 3) Lex + parse
	l := lexer.New(source)
	p := parser.New(l)
	program := p.ParseProgram()

	// (Optional) dump the AST for debugging
	astJSON, _ := json.MarshalIndent(program, "", "  ")
	fmt.Printf("AST:\n%s\n", astJSON)

	// 4) Evaluate in a fresh environment
	env := environment.NewEnvironment()
	result := evaluator.Eval(program, env)
	if result != nil {
		fmt.Printf("Result: %s\n", result.Inspect())
	} else {
		fmt.Println("Result: <nil>")
	}
}
