package main

import (
	"go/ast"

	rs "github.com/nyarly/rocketsurgery"
	. "github.com/nyarly/rocketsurgery/shortcuts"
)

func endpointsStruct(i rs.Interface) ast.Decl {
	fl := &ast.FieldList{}
	for _, m := range i.Methods() {
		fl.List = append(fl.List, &ast.Field{Names: []*ast.Ident{m.Name()}, Type: Sel(Id("endpoint"), Id("Endpoint"))})
	}
	return StructDecl(Id("Endpoints"), fl)
}

func httpHandler(i rs.Interface) ast.Decl {
	handlerFn := fetchFuncDecl("NewHTTPHandler")

	// does this "inlining" process merit a helper akin to replaceIdent?
	handleCalls := []ast.Stmt{}
	for _, m := range i.Methods() {
		handleCall := fetchFuncDecl("inlineHandlerBuilder").Body.List[0].(*ast.ExprStmt).X.(*ast.CallExpr)

		handleCall = replaceLit(handleCall, `"/bar"`, `"`+m.pathName()+`"`).(*ast.CallExpr)
		handleCall = replaceIdent(handleCall, "ExampleEndpoint", m.Name()).(*ast.CallExpr)
		handleCall = replaceIdent(handleCall, "DecodeExampleRequest", decodeFuncName(m)).(*ast.CallExpr)
		handleCall = replaceIdent(handleCall, "EncodeExampleResponse", encodeFuncName(m)).(*ast.CallExpr)

		handleCalls = append(handleCalls, &ast.ExprStmt{X: handleCall})
	}

	pasteStmts(handlerFn.Body, 1, handleCalls)

	return handlerFn
}
