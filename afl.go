package main

import "go/ast"

func buildAFL(pkgName string, fname string, fuzzFunc *ast.FuncDecl) {
	funcName, cleanup := generateLibFuzzer(pkgName, fname, fuzzFunc)
	defer cleanup()

	command(`go-afl-build`, `-func`, funcName)
}
