package main

import "go/ast"

func buildAFL(pkgName string, fname string, fuzzFunc *ast.FuncDecl) {
	funcName, cleanup := generateLibFuzzer(pkgName, fname, fuzzFunc)
	defer cleanup()

	outFile := `afl`
	if *outputFile != `` {
		outFile = *outputFile
	}
	args := getBuildArgs(aflFlags, `-func`, funcName, `-o`, outFile)
	command(`go-afl-build`, args...)
}
