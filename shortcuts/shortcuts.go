// Package shortcuts includes quick ways to do common things with AST manipulation.
// Because what is rocket surgery if you can't take some shortcuts?
// You are *encouraged* to use a dot import with this package (that is:
//     import . "github.com/nyarly/rocketsurgery/shortcuts"
// so that the functions here will be available without qualifier.
// Because what fun is using rocket surgery shortcuts if you can't self-amputate with high explosives?
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

func StructDecl(name *ast.Ident, fields *ast.FieldList) ast.Decl {
	return TypeDecl(&ast.TypeSpec{
		Name: name,
		Type: &ast.StructType{
			Fields: fields,
		},
	})
}

func TypeDecl(ts *ast.TypeSpec) ast.Decl {
	return &ast.GenDecl{
		Tok:   token.TYPE,
		Specs: []ast.Spec{ts},
	}
}

func StubFile(pkgname string) *ast.File {
	return &ast.File{
		Name:  Id(pkgname),
		Decls: []ast.Decl{},
	}
}

func InventName(scope *ast.Scope, t ast.Expr) *ast.Ident {
	n := BaseName(t)
	for try := 0; ; try++ {
		nstr := PickName(n, try)
		obj := ast.NewObj(ast.Var, nstr)
		if alt := scope.Insert(obj); alt == nil {
			return ast.NewIdent(nstr)
		}
	}
}

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

func PickName(base string, idx int) string {
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

func Export(s string) string {
	return strings.Title(s)
}

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

func TypeField(t ast.Expr) *ast.Field {
	return &ast.Field{Type: t}
}

func Field(n *ast.Ident, t ast.Expr) *ast.Field {
	return &ast.Field{
		Names: []*ast.Ident{n},
		Type:  t,
	}
}

func FieldList(list ...*ast.Field) *ast.FieldList {
	return &ast.FieldList{List: list}
}
