package main

import (
	"go/ast"
	"strings"

	_ "embed"
)

//go:embed tmpl/native.go.tmpl
var tmplNative string

func generateGoNative(pkgName string, fname string, fuzzFunc *ast.FuncDecl) {
	inputType := getInputType(fuzzFunc.Type.Params.List[0].Type)

	outFile := *outputFile
	if outFile == `` {
		outFile = strings.TrimSuffix(fname, `.go`) + `_gofuzzbuild_test.go`
	}

	_ = createTemplate(
		tmplNative, outFile,
		pkgName, fuzzFunc.Name.Name, inputType,
	)
}
