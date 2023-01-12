package main

import "go/ast"

func buildAFL(pkgName string, fname string, fuzzFunc *ast.FuncDecl) {
	funcName, cleanup := generateLibFuzzer(pkgName, fname, fuzzFunc)
	defer cleanup()

	outFile := `afl`
	if *outputFile != `` {
		outFile = *outputFile
	}

	command(`go-afl-build`, `-func`, funcName, `-o`, outFile)
}
