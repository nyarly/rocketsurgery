package rocketsurgery

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
)

type (
	ASTTemplate interface {
		Imports() []*ast.ImportSpec
		FunctionDecl(name string) *ast.FuncDecl
	}

	astTemplate struct {
		name string
		file *ast.File
	}
)

// full, err := ASTTemplates.Open("full.go")

func LoadAST(name string, full io.Reader) ASTTemplate {
	f, err := parser.ParseFile(token.NewFileSet(), name, full, parser.DeclarationErrors)
	if err != nil {
		panic(err)
	}
	return astTemplate{name: name, file: f}
}

func (astt astTemplate) Imports() []*ast.ImportSpec {
	return astt.file.Imports
}

func (astt astTemplate) FunctionDecl(name string) *ast.FuncDecl {
	for _, decl := range astt.file.Decls {
		if f, ok := decl.(*ast.FuncDecl); ok && f.Name.Name == name {
			return f
		}
	}
	panic(fmt.Errorf("No function called %q in %q", name, astt.name))
}
