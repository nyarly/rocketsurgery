package rocketsurgery

import (
	"go/ast"

	. "github.com/nyarly/rocketsurgery/shortcuts"
)

type (
	Method interface {
		Name() *ast.Ident
		Distinguished(scope *ast.Scope) Method
		Params() []Arg
		Results() []Arg
		ResultNames(scope *ast.Scope) []*ast.Ident
		Definition(s Struct, astt ASTTemplate, sourceName string) ast.Decl
		NonContextParams() []Arg
	}

	method struct {
		name            *ast.Ident
		params          []Arg
		results         []Arg
		structsResolved bool
	}
)

func (m method) Name() *ast.Ident {
	return m.name
}

func (m method) Distinguished(scope *ast.Scope) Method {
	name := m.name
	if scope.Lookup(m.name.Name) != nil {
		name = InventName(scope, m.name)
	}
	nm := method{name: name}

	scope = ast.NewScope(nil)
	for _, p := range m.params {
		nm.params = append(nm.params, p.Distinguish(scope))
	}
	scope = ast.NewScope(nil)
	for _, r := range m.results {
		nm.results = append(nm.results, r.Distinguish(scope))
	}
	return nm
}

func (m method) Params() []Arg {
	as := []Arg{}
	for _, a := range m.params {
		as = append(as, a)
	}
	return as
}

func (m method) Results() []Arg {
	as := []Arg{}
	for _, a := range m.results {
		as = append(as, a)
	}
	return as
}

// Definition produces a method declaration for this Method, with Struct as the
// receiver, and the body taken from that sourceName in astt.

// One possible refinement would be to align the template function with this
// method. There's no effort made (yet) to get the body of the template
// function to agree with this method. But maybe a like number of parameters could be replaced so that
//
//    func Template(x,y,z int) (a []int) {
//      return []int{x,y,z}
//    }
//
// could become
//
//    func (s struct) method(tom, dick, harry string) (a []string) {
//      return []string{tom,dick,harry}
//    }
//
// That'd be cool, right? It doesn't happen yet.
func (m method) Definition(s Struct, astt ASTTemplate, sourceName string) ast.Decl {
	notImpl := astt.FunctionDecl(sourceName)

	notImpl.Name = m.name
	notImpl.Recv = FieldList(s.Receiver())
	scope := ScopeWith(notImpl.Recv.List[0].Names[0].Name)
	notImpl.Type.Params = m.funcParams(scope)
	notImpl.Type.Results = m.funcResults()

	return notImpl
}

func (m method) funcResults() *ast.FieldList {
	return MappedFieldList(func(a Arg) *ast.Field {
		return a.AsResult()
	}, m.results...)
}

func (m method) funcParams(scope *ast.Scope) *ast.FieldList {
	parms := &ast.FieldList{}
	if HasContext(m) {
		parms.List = []*ast.Field{{
			Names: []*ast.Ident{ast.NewIdent("ctx")},
			Type:  Sel(Id("context"), Id("Context")),
		}}
		scope.Insert(ast.NewObj(ast.Var, "ctx"))
	}
	parms.List = append(parms.List, MappedFieldList(func(a Arg) *ast.Field {
		return a.Distinguish(scope).AsField()
	}, m.NonContextParams()...).List...)
	return parms
}

// Seems too specialized...
func (m method) NonContextParams() []Arg {
	if HasContext(m) {
		return m.params[1:]
	}
	return m.params
}

func (m method) ResultNames(scope *ast.Scope) []*ast.Ident {
	ids := []*ast.Ident{}
	for _, rz := range m.results {
		ids = append(ids, rz.Distinguish(scope).Name())
	}
	return ids
}

func HasContext(m Method) bool {
	if len(m.Params()) < 1 {
		return false
	}
	return SameSel(m.Params()[0].Type(), Sel(Id("context"), Id("Context")))
}
