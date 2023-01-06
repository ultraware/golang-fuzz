package main

import (
	"fmt"
	"go/ast"
	"os"
	"strings"
	"text/template"

	_ "embed"
)

//go:embed tmpl/native.go.tmpl
var tmplNative string

func generateGoNative(pkgName string, fname string, fuzzFunc *ast.FuncDecl) {
	tmpl, err := template.New(``).Parse(tmplNative)
	if err != nil {
		panic(err)
	}

	var input, inputType string // TODO: Move to functions, and support https://go.dev/security/fuzz/ types
	switch t := fuzzFunc.Type.Params.List[0].Type.(type) {
	case *ast.ArrayType:
		if v, ok := t.Elt.(*ast.Ident); ok {
			inputType = `[]` + v.Name
		} else {
			panic(fmt.Sprintf(`parameter ast type %T not supported`, t.Elt))
		}
	case *ast.Ident:
		inputType = t.Name
	default:
		panic(fmt.Sprintf(`parameter ast type %T not supported`, t))
	}

	switch inputType {
	case `[]byte`:
		input = `input`
	case `string`:
		input = `string(input)`
	default:
		panic(fmt.Sprintf(`parameter type %s not supported`, inputType))
	}

	fname = strings.TrimSuffix(fname, `.go`) + `_test.go`
	fuzzFile, err := os.OpenFile(fname, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o644)
	if err != nil {
		panic(err)
	}
	defer fuzzFile.Close()

	err = tmpl.Execute(fuzzFile, struct {
		PkgName   string
		FuncName  string
		CorpusDir string
		InputType string
		Input     string
	}{
		PkgName:   pkgName,
		FuncName:  fuzzFunc.Name.Name,
		CorpusDir: *corpusDir,
		InputType: inputType,
		Input:     input,
	})
	if err != nil {
		panic(err)
	}
}
