package main

import (
	"bufio"
	"compiler/environment"
	"compiler/evaluator"
	"compiler/lexer"
	"compiler/parser"
	"encoding/json"
	"fmt"
	"log"
	"os"
)

func main() {
	file, err := os.Open("../../files/main.blue")
	if err != nil {
		log.Fatalf("Error: %s", err)
	}
	// source := "let varName = 3 + 3;"

	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		source := scanner.Text()
		l := lexer.New(source)
		p := parser.New(l)

		prog := p.ParseProgram()
		parsedprogram, _ := json.MarshalIndent(prog, " ", "	")
		fmt.Printf("program: %s", parsedprogram)
		env := environment.NewEnvironment()
		result := evaluator.Eval(prog, env)

		fmt.Printf("\nResult: %s\n", result.Inspect())

		if val, ok := env.Get("varName"); ok {
			fmt.Printf("varName = %s\n", val.Inspect())
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatalf("Error: %s", err)
	}
}
