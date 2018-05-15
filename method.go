package rocketsurgery

import (
	"go/ast"

	. "github.com/nyarly/rocketsurgery/shortcuts"
)

type (
	Method interface {
		Name() *ast.Ident
		Params() []Arg
		Results() []Arg
		Definition(s Struct, astt ASTTemplate, sourceName string) ast.Decl
	}

	method struct {
		name            *ast.Ident
		params          []arg
		results         []arg
		structsResolved bool
	}
)

func (m method) Name() *ast.Ident {
	return m.name
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
	scope := scopeWith(notImpl.Recv.List[0].Names[0].Name)
	notImpl.Type.Params = m.funcParams(scope)
	notImpl.Type.Results = m.funcResults()

	return notImpl
}

func (m method) funcResults() *ast.FieldList {
	return mappedFieldList(func(a arg) *ast.Field {
		return a.result()
	}, m.results...)
}

func (m method) funcParams(scope *ast.Scope) *ast.FieldList {
	parms := &ast.FieldList{}
	if m.hasContext() {
		parms.List = []*ast.Field{{
			Names: []*ast.Ident{ast.NewIdent("ctx")},
			Type:  Sel(Id("context"), Id("Context")),
		}}
		scope.Insert(ast.NewObj(ast.Var, "ctx"))
	}
	parms.List = append(parms.List, mappedFieldList(func(a arg) *ast.Field {
		return a.field(scope)
	}, m.nonContextParams()...).List...)
	return parms
}

func (m method) nonContextParams() []arg {
	if m.hasContext() {
		return m.params[1:]
	}
	return m.params
}

func (m method) hasContext() bool {
	if len(m.params) < 1 {
		return false
	}
	carg := m.params[0].typ
	// ugh. this is maybe okay for the one-off, but a general case for matching
	// types would be helpful
	if sel, is := carg.(*ast.SelectorExpr); is && sel.Sel.Name == "Context" {
		if id, is := sel.X.(*ast.Ident); is && id.Name == "context" {
			return true
		}
	}
	return false
}

func (m method) resolveStructNames() {
	if m.structsResolved {
		return
	}
	m.structsResolved = true
	scope := ast.NewScope(nil)
	for i, p := range m.params {
		p.asField = p.chooseName(scope)
		m.params[i] = p
	}
	scope = ast.NewScope(nil)
	for i, r := range m.results {
		r.asField = r.chooseName(scope)
		m.results[i] = r
	}
}

func (m method) resultNames(scope *ast.Scope) []*ast.Ident {
	ids := []*ast.Ident{}
	for _, rz := range m.results {
		ids = append(ids, rz.chooseName(scope))
	}
	return ids
}
