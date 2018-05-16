package main

import (
	"go/ast"

	rs "github.com/nyarly/rocketsurgery"
	. "github.com/nyarly/rocketsurgery/shortcuts"
)

func addStubStruct(root *ast.File, s rs.Struct) {
	root.Decls = append(root.Decls, s.Decl())
}

func addType(root *ast.File, typ *ast.TypeSpec) {
	root.Decls = append(root.Decls, TypeDecl(typ))
}

func addMethod(at rs.ASTTemplate, root *ast.File, ifc rs.Interface, meth rs.Method) {
	s := ifc.Implementor()
	def := meth.Definition(s, at, "ExampleEndpoint")
	root.Decls = append(root.Decls, def)
}

func addRequestStruct(root *ast.File, meth rs.Method) {
	root.Decls = append(root.Decls, requestStruct(meth))
}

func addResponseStruct(root *ast.File, meth rs.Method) {
	root.Decls = append(root.Decls, responseStruct(meth))
}

func addEndpointMaker(at rs.ASTTemplate, root *ast.File, ifc rs.Interface, meth rs.Method) {
	s := ifc.Implementor()
	root.Decls = append(root.Decls, endpointMaker(at, s, meth))
}

func addEndpointsStruct(root *ast.File, ifc rs.Interface) {
	root.Decls = append(root.Decls, endpointsStruct(ifc))
}

func addHTTPHandler(at rs.ASTTemplate, root *ast.File, ifc rs.Interface) {
	root.Decls = append(root.Decls, httpHandler(at, ifc))
}

func addDecoder(at rs.ASTTemplate, root *ast.File, meth rs.Method) {
	root.Decls = append(root.Decls, decoderFunc(at, meth))
}

func addEncoder(at rs.ASTTemplate, root *ast.File, meth rs.Method) {
	root.Decls = append(root.Decls, encoderFunc(at, meth))
}
