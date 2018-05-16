package rocketsurgery

import (
	"go/ast"

	. "github.com/nyarly/rocketsurgery/shortcuts"
)

type (
	Arg interface {
		Name() *ast.Ident
		Type() ast.Expr
		Distinguish(scope *ast.Scope) Arg
		AsField() *ast.Field
		AsResult() *ast.Field
		Exported() *ast.Field
	}

	arg struct {
		name *ast.Ident
		typ  ast.Expr
	}
)

func (a arg) Name() *ast.Ident {
	return a.name
}

func (a arg) Type() ast.Expr {
	return a.typ
}

func (a arg) Distinguish(scope *ast.Scope) Arg {
	if a.name == nil || scope.Lookup(a.name.Name) != nil {
		name := InventName(scope, a.typ)
		return arg{
			name: name,
			typ:  a.typ,
		}
	}
	return a
}

// XXX this has changed behavior - be sure to Distinguish
func (a arg) AsField() *ast.Field {
	return &ast.Field{
		Names: []*ast.Ident{a.name},
		Type:  a.typ,
	}
}

func (a arg) AsResult() *ast.Field {
	return &ast.Field{
		Names: nil,
		Type:  a.typ,
	}
}

func (a arg) Exported() *ast.Field {
	return &ast.Field{
		Names: []*ast.Ident{Id(Export(a.name.Name))},
		Type:  a.typ,
	}
}
