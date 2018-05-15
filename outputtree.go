package rocketsurgery

import (
	"go/ast"

	//. "github.com/nyarly/rocketsurgery/shortcuts"
	"github.com/pkg/errors"
)

type (
	OutputTree interface {
		AddFile(path string, file *ast.File)
		MapFiles(func(path string, file *ast.File) *ast.File)
		FormatNodes() (Files, error)
	}

	outputTree map[string]*ast.File
)

func NewOutputTree() OutputTree {
	return make(outputTree)
}

func (ot outputTree) AddFile(path string, file *ast.File) {
	ot[path] = file
}

func (ot outputTree) MapFiles(fn func(string, *ast.File) *ast.File) {
	for n, f := range ot {
		ot[n] = fn(n, f)
	}
}

func (ot outputTree) FormatNodes() (Files, error) {
	res := Files{}
	var err error
	for fn, node := range ot {
		res[fn], err = formatNode(fn, node)
		if err != nil {
			return nil, errors.Wrapf(err, "formatNodes")
		}
	}
	return res, nil
}
