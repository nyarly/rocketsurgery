// Package shortcuts includes quick ways to do common things with AST manipulation.
// Because what is rocket surgery if you can't take some shortcuts?
// You are *encouraged* to use a dot import with this package (that is:
//     import . "github.com/nyarly/rocketsurgery/shortcuts"
// so that the functions here will be available without qualifier.
// Because what fun is using rocket surgery shortcuts if you can't self-amputate with high explosives?
//
// Many of the functions in this package serve to quickly build 'go/ast' structs.
package shortcuts

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"
	"unicode"
)

// Sel constructs SelectorExpr (e.g. rocketsurgery.Sel) out of Idents.
// Example:
//     Sel(Id("rocketsurgery"), Id("Sel"))
func Sel(ids ...*ast.Ident) ast.Expr {
	switch len(ids) {
	default:
		return &ast.SelectorExpr{
			X:   Sel(ids[:len(ids)-1]...),
			Sel: ids[len(ids)-1],
		}
	case 1:
		return ids[0]
	case 0:
		panic("zero ids to Sel()")
	}
}

// SameSel compares two expresions and returns true if they're the same selector or ident.
func SameSel(left, right ast.Expr) bool {
	switch l := left.(type) {
	case *ast.Ident:
		if r, is := right.(*ast.Ident); is {
			return l.Name == r.Name
		}
		return false
	case *ast.SelectorExpr:
		if r, is := right.(*ast.SelectorExpr); is {
			return SameSel(l.X, r.X)
		}
		return false
	default:
		return false
	}
}

// Id constructs an Ident from a string.
func Id(name string) *ast.Ident {
	return ast.NewIdent(name)
}

// StructDecl constructs a struct declaration node, based on a name and a list of fields.
func StructDecl(name *ast.Ident, fields *ast.FieldList) ast.Decl {
	return TypeDecl(&ast.TypeSpec{
		Name: name,
		Type: &ast.StructType{
			Fields: fields,
		},
	})
}

// TypeDecl constructs a type declaration based on a TypeSpec.
func TypeDecl(ts *ast.TypeSpec) ast.Decl {
	return &ast.GenDecl{
		Tok:   token.TYPE,
		Specs: []ast.Spec{ts},
	}
}

// StubFile creates an *ast.File for a package name.
func StubFile(pkgname string) *ast.File {
	return &ast.File{
		Name:  Id(pkgname),
		Decls: []ast.Decl{},
	}
}

// InventName uses a base expression to figure out a unique name within a scope.
// It's otherwise very easy for generated code not to build because there's already a 'foo' in scope.
func InventName(scope *ast.Scope, t ast.Expr) *ast.Ident {
	n := BaseName(t)
	for try := 0; ; try++ {
		nstr := pickName(n, try)
		obj := ast.NewObj(ast.Var, nstr)
		if alt := scope.Insert(obj); alt == nil {
			return ast.NewIdent(nstr)
		}
	}
}

// BaseName chooses a name for a variable based on its type.
func BaseName(t ast.Expr) string {
	switch tt := t.(type) {
	default:
		panic(fmt.Sprintf("don't know how to choose a base name for %T (%[1]v)", tt))
	case *ast.ArrayType:
		return "slice"
	case *ast.Ident:
		return tt.Name
	case *ast.SelectorExpr:
		return tt.Sel.Name
	}
}

func pickName(base string, idx int) string {
	if idx == 0 {
		switch base {
		default:
			return strings.Split(base, "")[0]
		case "Context":
			return "ctx"
		case "error":
			return "err"
		}
	}
	return fmt.Sprintf("%s%d", base, idx)
}

// Export takes a string and forms the "Export" form of it.
func Export(s string) string {
	return strings.Title(s)
}

// Unexport takes a string and forms the "unexported" form of it.
func Unexport(s string) string {
	first := true
	return strings.Map(func(r rune) rune {
		if first {
			first = false
			return unicode.ToLower(r)
		}
		return r
	}, s)
}

// TypeField constructs a field with a particular type.
func TypeField(t ast.Expr) *ast.Field {
	return &ast.Field{Type: t}
}

// Field builds a Field for a name and type.
func Field(n *ast.Ident, t ast.Expr) *ast.Field {
	return &ast.Field{
		Names: []*ast.Ident{n},
		Type:  t,
	}
}

// FieldList builds a list of fields, as used for function parameters, function
// results, or struct members.
func FieldList(list ...*ast.Field) *ast.FieldList {
	return &ast.FieldList{List: list}
}
