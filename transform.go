package rocketsurgery

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/token"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/davecgh/go-spew/spew"
	. "github.com/nyarly/rocketsurgery/shortcuts"

	"golang.org/x/tools/imports"
)

type (
	Files map[string]io.Reader

	Transformer interface {
		TransformAST(SourceContext) (Files, error)
	}
)

// GetGopath gets the set Go path, or else returns an absolute path of the default path (i.e. "~/.go")
func GetGopath() string {
	gopath, set := os.LookupEnv("GOPATH")
	if !set {
		return filepath.Join(os.Getenv("HOME"), "go")
	}
	return gopath
}

func ImportPath(targetDir, gopath string) (string, error) {
	if !filepath.IsAbs(targetDir) {
		return "", fmt.Errorf("%q is not an absolute path", targetDir)
	}

	for _, dir := range filepath.SplitList(gopath) {
		abspath, err := filepath.Abs(dir)
		if err != nil {
			continue
		}
		srcPath := filepath.Join(abspath, "src")

		res, err := filepath.Rel(srcPath, targetDir)
		if err != nil {
			continue
		}
		if strings.Index(res, "..") == -1 {
			return res, nil
		}
	}
	return "", fmt.Errorf("%q is not in GOPATH (%s)", targetDir, gopath)

}

// Selectify ensures that a particular identName is considered to be a part of pkgName, imported from importPath
//     xxx needs and example
func Selectify(file *ast.File, pkgName, identName, importPath string) *ast.File {
	if file.Name.Name == pkgName {
		return file
	}

	selector := Sel(Id(pkgName), Id(identName))
	var did bool
	if file, did = selectifyIdent(identName, file, selector); did {
		addImport(file, importPath)
	}
	return file
}

type selIdentFn func(ast.Node, func(ast.Node)) Visitor

func (f selIdentFn) Visit(node ast.Node, r func(ast.Node)) Visitor {
	return f(node, r)
}

func selectifyIdent(identName string, file *ast.File, selector ast.Expr) (*ast.File, bool) {
	var replaced bool
	var r selIdentFn
	r = selIdentFn(func(node ast.Node, replaceWith func(ast.Node)) Visitor {
		switch id := node.(type) {
		case *ast.SelectorExpr:
			return nil
		case *ast.Ident:
			if id.Name == identName {
				replaced = true
				replaceWith(selector)
			}
		}
		return r
	})
	return WalkReplace(r, file).(*ast.File), replaced
}

func formatNode(fname string, node ast.Node) (*bytes.Buffer, error) {
	if file, is := node.(*ast.File); is {
		sort.Stable(sortableDecls(file.Decls))
	}
	outfset := token.NewFileSet()
	buf := &bytes.Buffer{}
	err := format.Node(buf, outfset, node)
	if err != nil {
		return nil, err
	}
	imps, err := imports.Process(fname, buf.Bytes(), nil)
	if err != nil {
		return nil, err
	}
	return bytes.NewBuffer(imps), nil
}

type sortableDecls []ast.Decl

func (sd sortableDecls) Len() int {
	return len(sd)
}

func (sd sortableDecls) Less(i int, j int) bool {
	switch left := sd[i].(type) {
	case *ast.GenDecl:
		switch right := sd[j].(type) {
		default:
			return left.Tok == token.IMPORT
		case *ast.GenDecl:
			return left.Tok == token.IMPORT && right.Tok != token.IMPORT
		}
	}
	return false
}

func (sd sortableDecls) Swap(i int, j int) {
	sd[i], sd[j] = sd[j], sd[i]
}

// XXX debug
func spewDecls(f *ast.File) {
	for _, d := range f.Decls {
		switch dcl := d.(type) {
		default:
			spew.Dump(dcl)
		case *ast.GenDecl:
			spew.Dump(dcl.Tok)
		case *ast.FuncDecl:
			spew.Dump(dcl.Name.Name)
		}
	}
}

func addImport(root *ast.File, path string) {
	for _, d := range root.Decls {
		if imp, is := d.(*ast.GenDecl); is && imp.Tok == token.IMPORT {
			for _, s := range imp.Specs {
				if s.(*ast.ImportSpec).Path.Value == `"`+path+`"` {
					return // already have one
					// xxx aliased imports?
				}
			}
		}
	}
	root.Decls = append(root.Decls, importFor(importSpec(path)))
}
