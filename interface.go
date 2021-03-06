package rocketsurgery

import (
	"go/ast"

	. "github.com/nyarly/rocketsurgery/shortcuts"
)

type (
	// An Interface describes a Go interface.
	Interface interface {
		// Methods returns a list of methods for this interface.
		Methods() []Procedure
		// Implementor returns a Struct suitable to use an an implementor of this
		// interface.
		Implementor() Struct
	}
	// because "interface" is a keyword...
	iface struct {
		name, stubname, rcvrName *ast.Ident
		methods                  []Procedure
	}
)

func (i iface) Methods() []Procedure {
	ms := []Procedure{}
	for _, m := range i.methods {
		ms = append(ms, m)
	}
	return ms
}

func (i iface) Implementor() Struct {
	s := strct{name: *i.name}
	s.name.Name = i.name.String()
	return s
}

func (i iface) stubName() *ast.Ident {
	return i.stubname
}

// xxx this should return a `struct` type, and that should have a Decl() method.
func (i iface) stubStructDecl() ast.Decl {
	return StructDecl(i.stubName(), &ast.FieldList{})
}
