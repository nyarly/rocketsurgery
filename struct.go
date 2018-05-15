package rocketsurgery

import (
	"go/ast"

	. "github.com/nyarly/rocketsurgery/shortcuts"
)

type (
	Struct interface {
		Name() *ast.Ident
		Decl() ast.Decl
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
	return field(s.receiverName(), &s.name)
}

// XXX doesn't seem quite right
func (s strct) receiverName() *ast.Ident {
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
