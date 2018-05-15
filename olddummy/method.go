package main

import (
	"go/ast"
	"go/token"
	"strings"

	rs "github.com/nyarly/rocketsurgery"
	. "github.com/nyarly/rocketsurgery/shortcuts"
)

func endpointMaker(m rs.Method, ifc rs.Interface) ast.Decl {
	endpointFn := fetchFuncDecl("makeExampleEndpoint")
	scope := scopeWith("ctx", "req", ifc.receiverName().Name)

	anonFunc := endpointFn.Body.List[0].(*ast.ReturnStmt).Results[0].(*ast.FuncLit)
	if !m.hasContext() {
		// strip context param from endpoint function
		anonFunc.Type.Params.List = anonFunc.Type.Params.List[1:]
	}

	anonFunc = replaceIdent(anonFunc, "ExampleRequest", m.requestStructName()).(*ast.FuncLit)
	callMethod := m.called(ifc, scope, "ctx", "req")
	anonFunc.Body.List[1] = callMethod
	anonFunc.Body.List[2].(*ast.ReturnStmt).Results[0] = m.wrapResult(callMethod.Lhs)

	endpointFn.Body.List[0].(*ast.ReturnStmt).Results[0] = anonFunc
	endpointFn.Name = m.endpointMakerName()
	endpointFn.Type.Params = fieldList(ifc.reciever())
	endpointFn.Type.Results = fieldList(typeField(sel(Id("endpoint"), Id("Endpoint"))))
	return endpointFn
}

func pathName(m rs.Method) string {
	return "/" + strings.ToLower(m.name.Name)
}

func encodeFuncName(m rs.Method) *ast.Ident {
	return Id("Encode" + m.name.Name + "Response")
}

func decodeFuncName(m rs.Method) *ast.Ident {
	return Id("Decode" + m.name.Name + "Request")
}

// xxx make generic?
func called(m rs.Method, ifc rs.Interface, scope *ast.Scope, ctxName, spreadStruct string) *ast.AssignStmt {
	m.resolveStructNames()

	resNamesExpr := []ast.Expr{}
	for _, r := range m.resultNames(scope) {
		resNamesExpr = append(resNamesExpr, ast.Expr(r))
	}

	arglist := []ast.Expr{}
	if m.hasContext() {
		arglist = append(arglist, Id(ctxName))
	}
	ssid := Id(spreadStruct)
	for _, f := range m.requestStructFields().List {
		arglist = append(arglist, sel(ssid, f.Names[0]))
	}

	return &ast.AssignStmt{
		Lhs: resNamesExpr,
		Tok: token.DEFINE,
		Rhs: []ast.Expr{
			&ast.CallExpr{
				Fun:  sel(ifc.receiverName(), m.name),
				Args: arglist,
			},
		},
	}
}

func wrapResult(m rs.Method, results []ast.Expr) ast.Expr {
	kvs := []ast.Expr{}
	m.resolveStructNames()

	for i, a := range m.results {
		kvs = append(kvs, &ast.KeyValueExpr{
			Key:   ast.NewIdent(export(a.asField.Name)),
			Value: results[i],
		})
	}
	return &ast.CompositeLit{
		Type: m.responseStructName(),
		Elts: kvs,
	}
}

func decoderFunc(m rs.Method) ast.Decl {
	fn := fetchFuncDecl("DecodeExampleRequest")
	fn.Name = m.decodeFuncName()
	fn = replaceIdent(fn, "ExampleRequest", m.requestStructName()).(*ast.FuncDecl)
	return fn
}

func encoderFunc(m rs.Method) ast.Decl {
	fn := fetchFuncDecl("EncodeExampleResponse")
	fn.Name = m.encodeFuncName()
	return fn
}

func endpointMakerName(m rs.Method) *ast.Ident {
	return Id("make" + m.name.Name + "Endpoint")
}

func requestStruct(m rs.Method) ast.Decl {
	m.resolveStructNames()
	return structDecl(m.requestStructName(), m.requestStructFields())
}

func responseStruct(m rs.Method) ast.Decl {
	m.resolveStructNames()
	return structDecl(m.responseStructName(), m.responseStructFields())
}

func requestStructName(m rs.Method) *ast.Ident {
	return Id(export(m.name.Name) + "Request")
}

func requestStructFields(m rs.Method) *ast.FieldList {
	return mappedFieldList(func(a arg) *ast.Field {
		return a.exported()
	}, m.nonContextParams()...)
}

func responseStructName(m rs.Method) *ast.Ident {
	return Id(export(m.name.Name) + "Response")
}

func responseStructFields(m rs.Method) *ast.FieldList {
	return mappedFieldList(func(a arg) *ast.Field {
		return a.exported()
	}, m.results...)
}
