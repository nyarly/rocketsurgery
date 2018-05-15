package rocketsurgery

import (
	"go/ast"
	"go/token"
)

func scopeWith(names ...string) *ast.Scope {
	scope := ast.NewScope(nil)
	for _, name := range names {
		scope.Insert(ast.NewObj(ast.Var, name))
	}
	return scope
}

type visitFn func(ast.Node, func(ast.Node))

func (fn visitFn) Visit(node ast.Node, r func(ast.Node)) Visitor {
	fn(node, r)
	return fn
}

func replaceIdent(src ast.Node, named string, with ast.Node) ast.Node {
	r := visitFn(func(node ast.Node, replaceWith func(ast.Node)) {
		switch id := node.(type) {
		case *ast.Ident:
			if id.Name == named {
				replaceWith(with)
			}
		}
	})
	return WalkReplace(r, src)
}

func replaceLit(src ast.Node, from, to string) ast.Node {
	r := visitFn(func(node ast.Node, replaceWith func(ast.Node)) {
		switch lit := node.(type) {
		case *ast.BasicLit:
			if lit.Value == from {
				replaceWith(&ast.BasicLit{Value: to})
			}
		}
	})
	return WalkReplace(r, src)
}

func typeField(t ast.Expr) *ast.Field {
	return &ast.Field{Type: t}
}

func field(n *ast.Ident, t ast.Expr) *ast.Field {
	return &ast.Field{
		Names: []*ast.Ident{n},
		Type:  t,
	}
}

func fieldList(list ...*ast.Field) *ast.FieldList {
	return &ast.FieldList{List: list}
}

func mappedFieldList(fn func(arg) *ast.Field, args ...arg) *ast.FieldList {
	fl := &ast.FieldList{List: []*ast.Field{}}
	for _, a := range args {
		fl.List = append(fl.List, fn(a))
	}
	return fl
}

func blockStmt(stmts ...ast.Stmt) *ast.BlockStmt {
	return &ast.BlockStmt{
		List: stmts,
	}
}

func pasteStmts(body *ast.BlockStmt, idx int, stmts []ast.Stmt) {
	list := body.List
	prefix := list[:idx]
	suffix := make([]ast.Stmt, len(list)-idx-1)
	copy(suffix, list[idx+1:])

	body.List = append(append(prefix, stmts...), suffix...)
}

func importFor(is *ast.ImportSpec) *ast.GenDecl {
	return &ast.GenDecl{Tok: token.IMPORT, Specs: []ast.Spec{is}}
}

func importSpec(path string) *ast.ImportSpec {
	return &ast.ImportSpec{Path: &ast.BasicLit{Kind: token.STRING, Value: `"` + path + `"`}}
}
