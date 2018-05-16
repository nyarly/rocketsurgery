package rocketsurgery

import (
	"go/ast"

	. "github.com/nyarly/rocketsurgery/shortcuts"
)

type (
	// An Arg represents an argument of method or function - that is, a parameter
	// or result.
	Arg interface {
		// Name returns an Ident for use as the name of the Arg.
		Name() *ast.Ident
		// Name returns an Expr for use as the type of the Arg.
		Type() ast.Expr
		// Distinguish returns an Arg whose name is unique within the given Scope.
		Distinguish(scope *ast.Scope) Arg
		// AsField returns a Field node with both name and type.
		AsField() *ast.Field
		// AsResult returns a Field node with a blank name.
		AsResult() *ast.Field
		// Exported returns a Field with the name Exported.
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
