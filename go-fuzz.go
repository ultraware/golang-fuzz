package main

import "go/ast"

func buildGoFuzz(pkgName string, fname string, fuzzFunc *ast.FuncDecl) {
	funcName, cleanup := generateLibFuzzer(pkgName, fname, fuzzFunc)
	defer cleanup()

	command(`go`, `get`, `-u`, `github.com/dvyukov/go-fuzz/go-fuzz-dep`)
	command(`go-fuzz-build`, `-func`, funcName)
}
