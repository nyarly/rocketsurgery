package main

import (
	"go/ast"

	rs "github.com/nyarly/rocketsurgery"
)

type flat struct {
	tmpl rs.ASTTemplate
}

func (f flat) TransformAST(ctx rs.SourceContext) (rs.Files, error) {
	root := &ast.File{
		Name:  ctx.Package(),
		Decls: []ast.Decl{},
	}

	ctx.AddImports(root, f.tmpl)

	for _, typ := range ctx.Types() {
		addType(root, typ)
	}

	for _, iface := range ctx.Interfaces() { //only one...
		s := iface.Implementor()
		addStubStruct(root, s)

		for _, meth := range iface.Methods() {
			addMethod(f.tmpl, root, iface, meth)
			addRequestStruct(root, meth)
			addResponseStruct(root, meth)
			addEndpointMaker(f.tmpl, root, iface, meth)
		}

		addEndpointsStruct(root, iface)
		addHTTPHandler(f.tmpl, root, iface)

		for _, meth := range iface.Methods() {
			addDecoder(f.tmpl, root, meth)
			addEncoder(f.tmpl, root, meth)
		}
	}

	out := rs.NewOutputTree()
	out.AddFile("gokit.go", root)

	return out.FormatNodes()
}
