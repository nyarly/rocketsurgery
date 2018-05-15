package main

import "go/ast"
import rs "github.com/nyarly/rocketsurgery"

func addStubStruct(root *ast.File, s rs.Struct) {
	root.Decls = append(root.Decls, s.Decl())
}

func addType(root *ast.File, typ *ast.TypeSpec) {
	root.Decls = append(root.Decls, typeDecl(typ))
}

func addMethod(root *ast.File, iface rs.Interface, meth rs.Method) {
	def := meth.definition(iface)
	root.Decls = append(root.Decls, def)
}

func addRequestStruct(root *ast.File, meth rs.Method) {
	root.Decls = append(root.Decls, requestStruct(meth))
}

func addResponseStruct(root *ast.File, meth rs.Method) {
	root.Decls = append(root.Decls, responseStruct(meth))
}

func addEndpointMaker(root *ast.File, ifc rs.Interface, meth rs.Method) {
	root.Decls = append(root.Decls, endpointMaker(meth, ifc))
}

func addEndpointsStruct(root *ast.File, ifc rs.Interface) {
	root.Decls = append(root.Decls, endpointsStruct(ifc))
}

func addHTTPHandler(root *ast.File, ifc rs.Interface) {
	root.Decls = append(root.Decls, httpHandler(ifc))
}

func addDecoder(root *ast.File, meth rs.Method) {
	root.Decls = append(root.Decls, decoderFunc(meth))
}

func addEncoder(root *ast.File, meth rs.Method) {
	root.Decls = append(root.Decls, encoderFunc(meth))
}
