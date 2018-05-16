package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path"

	rs "github.com/nyarly/rocketsurgery"
	"github.com/pkg/errors"
)

// go get github.com/nyarly/inlinefiles
//go:generate inlinefiles --package=main --vfs=ASTTemplates ./templates ast_templates.go

func usage() string {
	return fmt.Sprintf("Usage: %s <filename> (try -h)", os.Args[0])
}

var (
	help       = flag.Bool("h", false, "print this help")
	layoutkind = flag.String("repo-layout", "default", "default, flat...")
	outdirrel  = flag.String("target-dir", ".", "base directory to emit into")
	//contextOmittable = flag.Bool("allow-no-context", false, "allow service methods to omit context parameter")
)

func helpText() {
	fmt.Println("USAGE")
	fmt.Println("  kitgen [flags] path/to/service.go")
	fmt.Println("")
	fmt.Println("FLAGS")
	flag.PrintDefaults()
}

func main() {
	flag.Parse()

	if *help {
		helpText()
		os.Exit(0)
	}

	outdir := *outdirrel
	if !path.IsAbs(*outdirrel) {
		wd, err := os.Getwd()
		if err != nil {
			log.Fatalf("error getting current working directory: %v", err)
		}
		outdir = path.Join(wd, *outdirrel)
	}

	var layout rs.Transformer
	switch *layoutkind {
	default:
		log.Fatalf("Unrecognized layout kind: %q - try 'default' or 'flat'", *layoutkind)
	case "default":
		gopath := rs.GetGopath()
		importBase, err := rs.ImportPath(outdir, gopath)
		if err != nil {
			log.Fatal(err)
		}
		layout = deflayout{targetDir: importBase, tmpl: fullAST()}
	case "flat":
		layout = flat{tmpl: fullAST()}
	}

	if len(os.Args) < 2 {
		log.Fatal(usage())
	}
	filename := flag.Arg(0)
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("error while opening %q: %v", filename, err)
	}

	tree, err := process(filename, file, layout)
	if err != nil {
		log.Fatal(err)
	}

	err = splat(outdir, tree)
	if err != nil {
		log.Fatal(err)
	}
}

func fullAST() rs.ASTTemplate {
	astfile, err := ASTTemplates.Open("full.go")
	if err != nil {
		log.Fatal(err)
	}
	return rs.LoadAST("full.go", astfile)
}

func process(filename string, source io.Reader, layout rs.Transformer) (rs.Files, error) {
	context, err := rs.ParseReader(filename, source)
	if err != nil {
		return nil, errors.Wrapf(err, "parsing input code")
	}
	tree, err := layout.TransformAST(context)
	return tree, errors.Wrapf(err, "generating AST")
}

/*
	buf, err := formatNode(dest)
	if err != nil {
		return nil, errors.Wrapf(err, "formatting")
	}
	return buf, nil
}
*/

func splat(dir string, tree rs.Files) error {
	for fn, buf := range tree {
		if err := splatFile(path.Join(dir, fn), buf); err != nil {
			return err
		}
	}
	return nil
}

func splatFile(target string, buf io.Reader) error {
	err := os.MkdirAll(path.Dir(target), os.ModePerm)
	if err != nil {
		return errors.Wrapf(err, "Couldn't create directory for %q", target)
	}
	f, err := os.Create(target)
	if err != nil {
		return errors.Wrapf(err, "Couldn't create file %q", target)
	}
	defer f.Close()
	_, err = io.Copy(f, buf)
	return errors.Wrapf(err, "Error writing data to file %q", target)
}
