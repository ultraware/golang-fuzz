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

	defaultOutputFile := strings.TrimSuffix(fname, `.go`) + `_golangfuzz_test.go`

	_ = createTemplate(
		tmplNative, getOutputFile(defaultOutputFile),
		pkgName, fuzzFunc.Name.Name, inputType,
	)
}

func runGoNative(pkgName string, fname string, fuzzFunc *ast.FuncDecl, args []string) {
	args = append([]string{`test`, `-fuzz`, fuzzFunc.Name.Name + `_`}, args...)
	execute(true, `go`, args...)
}
