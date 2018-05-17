package rocketsurgery

import (
	"go/ast"

	. "github.com/nyarly/rocketsurgery/shortcuts"
)

type (
	// Struct represents a Go struct.
	Struct interface {
		// Name returns an Ident for the name of the struct.
		Name() *ast.Ident
		// Name returns an Ident for the name of the struct.
		ReceiverName() *ast.Ident
		// Decl returns a Decl node declaring the struct.
		Decl() ast.Decl
		// Receiver returns a Field node suitable to use this struct as the receiver of a method.
		Receiver() *ast.Field
	}

	strct struct {
		name     ast.Ident
		rcvrName *ast.Ident
		methods  []Method
	}
)

// xxx this should return a `struct` type, and that should have a Decl() method.
func (s strct) Decl() ast.Decl {
	return StructDecl(&s.name, &ast.FieldList{})
}

func (s strct) Name() *ast.Ident {
	return &s.name
}

func (s strct) Receiver() *ast.Field {
	return Field(s.ReceiverName(), &s.name)
}

// XXX doesn't seem quite right
func (s strct) ReceiverName() *ast.Ident {
	if s.rcvrName != nil {
		return s.rcvrName
	}
	scope := ast.NewScope(nil)
	for _, meth := range s.methods {
		for _, arg := range meth.Params() {
			if name := arg.Name(); name != nil {
				scope.Insert(ast.NewObj(ast.Var, name.Name))
			}
		}
		for _, arg := range meth.Results() {
			if name := arg.Name(); name != nil {
				scope.Insert(ast.NewObj(ast.Var, name.Name))
			}
		}
	}
	return Id(Unexport(InventName(scope, &s.name).Name))
}
