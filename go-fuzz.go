package main

import (
	"go/ast"
)

func buildGoFuzz(pkgName string, fname string, fuzzFunc *ast.FuncDecl) {
	funcName, cleanup := generateLibFuzzer(pkgName, fname, fuzzFunc)
	defer cleanup()

	args := getBuildArgs(gofuzzFlags, `-func`, funcName)
	if *outputFile != `` {
		args = append(args, `-o`, *outputFile)
	}

	command(`go`, `get`, `-u`, `github.com/dvyukov/go-fuzz/go-fuzz-dep`)
	command(`golang-fuzz`, args...)
}

func runGoFuzz(args []string) {
	if *corpusDir != `corpus` {
		panic(`corpus directory must be called corpus`)
	}
	execute(true, `go-fuzz`, args...)
}
