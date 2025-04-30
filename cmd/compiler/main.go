package main

import (
    "fmt"
    "compiler/lexer"
    "compiler/parser"
    "compiler/codegen"
	  "github.com/go-llvm/llvm/bindings/go/llvm"
)

func main() {
    source := "varName = 3"
    l := lexer.New(source)
    p := parser.New(l)
    prog := p.ParseProgram()  // *ast.Program

    cg := codegen.New("my_module")
    mod := cg.Generate(prog)

    fmt.Println(mod.String())  // print the LLVM IR
}
