package codegen

import (
    "compiler/ast"
    "llvm.org/llvm/bindings/go/llvm"
)

type CodeGenerator struct {
    module *llvm.Module
    builder llvm.Builder
    symtab map[string]llvm.Value  // maps var names → alloca pointers
}

func New(moduleName string) *CodeGenerator {
    llvm.InitializeAllTargetInfos()
    llvm.InitializeAllTargets()
    llvm.InitializeAllTargetMCs()
    llvm.InitializeAllAsmParsers()
    llvm.InitializeAllAsmPrinters()

    mod := llvm.NewModule(moduleName)
    builder := llvm.NewBuilder()

    return &CodeGenerator{
        module:  mod,
        builder: builder,
        symtab:  map[string]llvm.Value{},
    }
}

func (cg *CodeGenerator) Generate(p *ast.Program) llvm.Module {
    // Create a 'main' function
    fnType := llvm.FunctionType(llvm.VoidType(), nil, false)
    mainFn := llvm.AddFunction(cg.module, "main", fnType)
    entry := llvm.AddBasicBlock(mainFn, "entry")
    cg.builder.SetInsertPointAtEnd(entry)

    // Walk statements
    for _, stmt := range p.Statements {
        cg.genStatement(stmt)
    }

    cg.builder.CreateRetVoid()
    return *cg.module
}

func (cg *CodeGenerator) genStatement(s ast.Statement) {
    switch st := s.(type) {
    case *ast.AssignmentStatement:
        cg.genAssignment(st)
    // ... other statement types ...
    }
}

func (cg *CodeGenerator) genAssignment(stmt *ast.AssignmentStatement) {
    // Evaluate the RHS
    val := cg.genExpression(stmt.Value)

    // If we haven't seen the variable, allocate it at entry
    ptr, ok := cg.symtab[stmt.Name.Value]
    if !ok {
        ptr = cg.builder.CreateAlloca(llvm.Int32Type(), stmt.Name.Value)
        cg.symtab[stmt.Name.Value] = ptr
    }

    // Store the value into the variable
    cg.builder.CreateStore(val, ptr)
}

func (cg *CodeGenerator) genExpression(expr ast.Expression) llvm.Value {
    switch e := expr.(type) {
    case *ast.IntegerLiteral:
        return llvm.ConstInt(llvm.Int32Type(), uint64(e.Value), false)
    case *ast.Identifier:
        // load from previously allocated variable
        ptr := cg.symtab[e.Value]
        return cg.builder.CreateLoad(llvm.Int32Type(), ptr, e.Value+"_load")
    // ... other expr kinds ...
    }
    return llvm.ConstNull(llvm.Int32Type())
}
