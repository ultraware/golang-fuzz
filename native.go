package main

import (
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

	inputType := getInputType(fuzzFunc.Type.Params.List[0].Type)
	inputCode, inputLen, imprts := getInputCode(inputType)

	fname = strings.TrimSuffix(fname, `.go`) + `_test.go`
	fuzzFile, err := os.OpenFile(fname, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o644)
	if err != nil {
		panic(err)
	}
	defer fuzzFile.Close()

	err = tmpl.Execute(fuzzFile, tmplData{
		PkgName:   pkgName,
		Imports:   imprts,
		FuncName:  fuzzFunc.Name.Name,
		CorpusDir: *corpusDir,
		InputType: inputType,
		InputCode: inputCode,
		InputLen:  inputLen,
	})
	if err != nil {
		panic(err)
	}
}
