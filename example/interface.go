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

func httpHandler(t rs.ASTTemplate, i rs.Interface) ast.Decl {
	handlerFn := t.FunctionDecl("NewHTTPHandler")

	// does this "inlining" process merit a helper akin to replaceIdent?
	handleCalls := []ast.Stmt{}
	for _, m := range i.Methods() {
		handleCall := t.FunctionDecl("inlineHandlerBuilder").Body.List[0].(*ast.ExprStmt).X.(*ast.CallExpr)

		handleCall = rs.ReplaceLit(handleCall, `"/bar"`, `"`+m.pathName()+`"`).(*ast.CallExpr)
		handleCall = rs.ReplaceIdent(handleCall, "ExampleEndpoint", m.Name()).(*ast.CallExpr)
		handleCall = rs.ReplaceIdent(handleCall, "DecodeExampleRequest", decodeFuncName(m)).(*ast.CallExpr)
		handleCall = rs.ReplaceIdent(handleCall, "EncodeExampleResponse", encodeFuncName(m)).(*ast.CallExpr)

		handleCalls = append(handleCalls, &ast.ExprStmt{X: handleCall})
	}

	rs.PasteStmts(handlerFn.Body, 1, handleCalls)

	return handlerFn
}
