package rocketsurgery

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"os"

	"github.com/pkg/errors"
)

type (
	// A SourceContext ...
	SourceContext interface {
		Package() *ast.Ident
		AddImports(*ast.File, ASTTemplate)
		Interfaces() []Interface
		Types() []*ast.TypeSpec
	}

	sourceContext struct {
		pkg        *ast.Ident
		imports    []*ast.ImportSpec
		interfaces []iface
		types      []*ast.TypeSpec
	}
)

func ParsePath(path string) (SourceContext, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, errors.Wrapf(err, "error while opening %q", path)
	}
	defer file.Close()
	return ParseReader(path, file)
}

// Example:
//    func process(filename string, source io.Reader, layout rs.Transformer) (rs.Files, error) {
//    	context, err := rs.ParseReader(filename, source)
//    	if err != nil {
//    		return nil, err
//    	}
//      return layout.TransformAST(context)
//    }
func ParseReader(filename string, source io.Reader) (SourceContext, error) {
	f, err := parser.ParseFile(token.NewFileSet(), filename, source, parser.DeclarationErrors)
	if err != nil {
		return nil, errors.Wrapf(err, "parsing input %q", filename)
	}

	return ExtractContext(f)
}

func ExtractContext(f ast.Node) (SourceContext, error) {
	context := &sourceContext{}
	visitor := &parseVisitor{src: context}

	ast.Walk(visitor, f)

	return context, context.validate()
}

func (sc *sourceContext) Interfaces() []Interface {
	is := []Interface{}
	for _, i := range sc.interfaces {
		is = append(is, Interface(&i))
	}
	return is
}

func (sc *sourceContext) Types() []*ast.TypeSpec {
	return sc.types
}

func (sc *sourceContext) Package() *ast.Ident {
	return sc.pkg
}

func (sc *sourceContext) validate() error {
	if len(sc.interfaces) != 1 {
		return fmt.Errorf("found %d interfaces, expecting exactly 1", len(sc.interfaces))
	}
	for _, i := range sc.interfaces {
		for _, m := range i.methods {
			if len(m.Results()) < 1 {
				return fmt.Errorf("method %q of interface %q has no result types", m.Name(), i.name)
			}
		}
	}
	return nil
}

func (sc *sourceContext) importDecls(astt ASTTemplate) (decls []ast.Decl) {
	have := map[string]struct{}{}
	notHave := func(is *ast.ImportSpec) bool {
		if _, has := have[is.Path.Value]; has {
			return false
		}
		have[is.Path.Value] = struct{}{}
		return true
	}

	for _, is := range sc.imports {
		if notHave(is) {
			decls = append(decls, importFor(is))
		}
	}

	for _, is := range astt.Imports() {
		if notHave(is) {
			decls = append(decls, &ast.GenDecl{Tok: token.IMPORT, Specs: []ast.Spec{is}})
		}
	}

	return
}

func (sc *sourceContext) AddImports(root *ast.File, astt ASTTemplate) {
	root.Decls = append(root.Decls, sc.importDecls(astt)...)
}
