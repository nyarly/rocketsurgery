package main

import (
	"go/ast"
	"go/token"
	"strings"

	rs "github.com/nyarly/rocketsurgery"
	. "github.com/nyarly/rocketsurgery/shortcuts"
)

func endpointMaker(t rs.ASTTemplate, s rs.Struct, m rs.Method) ast.Decl {
	endpointFn := t.FunctionDecl("makeExampleEndpoint")
	scope := rs.ScopeWith("ctx", "req", s.Name().Name)

	anonFunc := endpointFn.Body.List[0].(*ast.ReturnStmt).Results[0].(*ast.FuncLit)
	if !rs.HasContext(m) {
		// strip context param from endpoint function
		anonFunc.Type.Params.List = anonFunc.Type.Params.List[1:]
	}

	anonFunc = rs.ReplaceIdent(anonFunc, "ExampleRequest", requestStructName(m)).(*ast.FuncLit)
	callMethod := called(m, s, scope, "ctx", "req")
	anonFunc.Body.List[1] = callMethod
	anonFunc.Body.List[2].(*ast.ReturnStmt).Results[0] = wrapResult(m, callMethod.Lhs)

	endpointFn.Body.List[0].(*ast.ReturnStmt).Results[0] = anonFunc
	endpointFn.Name = endpointMakerName(m)
	endpointFn.Type.Params = FieldList(s.Receiver())
	endpointFn.Type.Results = FieldList(TypeField(Sel(Id("endpoint"), Id("Endpoint"))))
	return endpointFn
}

func pathName(m rs.Method) string {
	return "/" + strings.ToLower(m.Name().Name)
}

func encodeFuncName(m rs.Method) *ast.Ident {
	return Id("Encode" + m.Name().Name + "Response")
}

func decodeFuncName(m rs.Method) *ast.Ident {
	return Id("Decode" + m.Name().Name + "Request")
}

// xxx make generic?
func called(m rs.Method, s rs.Struct, scope *ast.Scope, ctxName, spreadStruct string) *ast.AssignStmt {
	resNamesExpr := []ast.Expr{}
	for _, r := range m.ResultNames(scope) {
		resNamesExpr = append(resNamesExpr, ast.Expr(r))
	}

	arglist := []ast.Expr{}
	if rs.HasContext(m) {
		arglist = append(arglist, Id(ctxName))
	}
	ssid := Id(spreadStruct)
	for _, f := range requestStructFields(m).List {
		arglist = append(arglist, Sel(ssid, f.Names[0]))
	}

	return &ast.AssignStmt{
		Lhs: resNamesExpr,
		Tok: token.DEFINE,
		Rhs: []ast.Expr{
			&ast.CallExpr{
				Fun:  Sel(s.Name(), m.Name()),
				Args: arglist,
			},
		},
	}
}

func wrapResult(m rs.Method, results []ast.Expr) ast.Expr {
	kvs := []ast.Expr{}

	scope := rs.ScopeWith()
	for i, a := range m.Results() {
		kvs = append(kvs, &ast.KeyValueExpr{
			Key:   ast.NewIdent(Export(a.Distinguish(scope).Name().Name)), //xxx
			Value: results[i],
		})
	}
	return &ast.CompositeLit{
		Type: responseStructName(m),
		Elts: kvs,
	}
}

func decoderFunc(t rs.ASTTemplate, m rs.Method) ast.Decl {
	fn := t.FunctionDecl("DecodeExampleRequest")
	fn.Name = decodeFuncName(m)
	fn = rs.ReplaceIdent(fn, "ExampleRequest", requestStructName(m)).(*ast.FuncDecl)
	return fn
}

func encoderFunc(t rs.ASTTemplate, m rs.Method) ast.Decl {
	fn := t.FunctionDecl("EncodeExampleResponse")
	fn.Name = encodeFuncName(m)
	return fn
}

func endpointMakerName(m rs.Method) *ast.Ident {
	return Id("make" + m.Name().Name + "Endpoint")
}

func requestStruct(m rs.Method) ast.Decl {
	return StructDecl(requestStructName(m), requestStructFields(m))
}

func responseStruct(m rs.Method) ast.Decl {
	return StructDecl(responseStructName(m), responseStructFields(m))
}

func requestStructName(m rs.Method) *ast.Ident {
	return Id(Export(m.Name().Name) + "Request")
}

func requestStructFields(m rs.Method) *ast.FieldList {
	return rs.MappedFieldList(func(a rs.Arg) *ast.Field {
		return a.Distinguish(rs.ScopeWith()).Exported()
	}, m.NonContextParams()...)
}

func responseStructName(m rs.Method) *ast.Ident {
	return Id(Export(m.Name().Name) + "Response")
}

func responseStructFields(m rs.Method) *ast.FieldList {
	return rs.MappedFieldList(func(a rs.Arg) *ast.Field {
		return a.Distinguish(rs.ScopeWith()).Exported()
	}, m.Results()...)
}
