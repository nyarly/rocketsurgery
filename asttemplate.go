package rocketsurgery

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
)

type (
	// An ASTTemplate holds an abstract syntax tree to that parts of it can be extracted and tranformed as needed.
	ASTTemplate interface {
		// Imports lists the imports defined in the template file.
		Imports() []*ast.ImportSpec

		// FunctionDecl returns the function declaration subtree for a given name (or panics if it isn't present.)
		FunctionDecl(name string) *ast.FuncDecl
	}

	astTemplate struct {
		name string
		buf  []byte
	}
)

// LoadAST loads an io.Reader with Go code as an ASTTemplate.
func LoadAST(name string, full io.Reader) ASTTemplate {
	b := &bytes.Buffer{}
	io.Copy(b, full)

	t := astTemplate{name: name, buf: b.Bytes()}

	t.reparse()
	return t
}

func (astt astTemplate) reparse() *ast.File {
	b := bytes.NewBuffer(astt.buf)
	f, err := parser.ParseFile(token.NewFileSet(), astt.name, b, parser.DeclarationErrors)
	if err != nil {
		panic(err)
	}
	return f
}

func (astt astTemplate) Imports() []*ast.ImportSpec {
	return astt.reparse().Imports
}

func (astt astTemplate) FunctionDecl(name string) *ast.FuncDecl {
	f := astt.reparse()
	for _, decl := range f.Decls {
		if f, ok := decl.(*ast.FuncDecl); ok && f.Name.Name == name {
			return f
		}
	}
	panic(fmt.Errorf("No function called %q in %q", name, astt.name))
}
