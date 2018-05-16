package main

import (
	"go/ast"
	"path/filepath"

	rs "github.com/nyarly/rocketsurgery"
	. "github.com/nyarly/rocketsurgery/shortcuts"
)

type deflayout struct {
	tmpl      rs.ASTTemplate
	targetDir string
}

func (l deflayout) packagePath(sub string) string {
	return filepath.Join(l.targetDir, sub)
}

func (l deflayout) TransformAST(ctx rs.SourceContext) (rs.Files, error) {
	out := rs.NewOutputTree()

	endpoints := StubFile("endpoints")
	out.AddFile("endpoints/endpoints.go", endpoints)

	http := StubFile("http")
	out.AddFile("http/http.go", http)

	service := StubFile("service")
	out.AddFile("service/service.go", service)

	ctx.AddImports(endpoints, l.tmpl)
	ctx.AddImports(http, l.tmpl)
	ctx.AddImports(service, l.tmpl)

	for _, typ := range ctx.Types() {
		addType(service, typ)
	}

	for _, iface := range ctx.Interfaces() { //only one...
		s := iface.Implementor()
		addStubStruct(service, s)

		for _, meth := range iface.Methods() {
			addMethod(l.tmpl, service, iface, meth)
			addRequestStruct(endpoints, meth)
			addResponseStruct(endpoints, meth)
			addEndpointMaker(l.tmpl, endpoints, iface, meth)
		}

		addEndpointsStruct(endpoints, iface)
		addHTTPHandler(l.tmpl, http, iface)

		for _, meth := range iface.Methods() {
			addDecoder(l.tmpl, http, meth)
			addEncoder(l.tmpl, http, meth)
		}

		out.MapFiles(func(name string, f *ast.File) *ast.File {
			f = rs.Selectify(f, "service", s.Name().Name, l.packagePath("service"))
			for _, meth := range iface.Methods() {
				f = rs.Selectify(f, "endpoints", requestStructName(meth).Name, l.packagePath("endpoints"))
			}
			return f
		})
	}

	out.MapFiles(func(name string, f *ast.File) *ast.File {
		f = rs.Selectify(f, "endpoints", "Endpoints", l.packagePath("endpoints"))

		for _, typ := range ctx.Types() {
			f = rs.Selectify(f, "service", typ.Name.Name, l.packagePath("service"))
		}
		return f
	})

	return out.FormatNodes()
}
