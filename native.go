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
	_ = createTemplate(
		tmplNative, ``+strings.TrimSuffix(fname, `.go`)+`_gofuzzbuild_test.go`,
		pkgName, fuzzFunc.Name.Name, inputType,
	)
}
