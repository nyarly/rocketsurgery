package rocketsurgery

import (
	"go/ast"
	"go/token"
)

func ScopeWith(names ...string) *ast.Scope {
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

func ReplaceIdent(src ast.Node, named string, with ast.Node) ast.Node {
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

func ReplaceLit(src ast.Node, from, to string) ast.Node {
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

func MappedFieldList(fn func(Arg) *ast.Field, args ...Arg) *ast.FieldList {
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

func PasteStmts(body *ast.BlockStmt, idx int, stmts []ast.Stmt) {
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
