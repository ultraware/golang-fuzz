package main

import (
	"go/ast"
	"path/filepath"

	"golang.org/x/exp/slices"
)

func buildAFL(pkgName string, fname string, fuzzFunc *ast.FuncDecl) {
	funcName, cleanup := generateLibFuzzer(pkgName, fname, fuzzFunc)
	defer cleanup()

	args := getBuildArgs(aflFlags, `-func`, funcName, `-o`, getOutputFile(`afl`))
	command(`go-afl-build`, args...)
}

func runAFL(args []string) {
	if !slices.Contains(args, `-t`) {
		panic(`must specify timeout (-t) in run_args`)
	}
	if !slices.Contains(args, `-o`) {
		args = append(args, `-o`, `.`)
	}

	valid, dir := isValidCorpusDir()
	if valid && !slices.Contains(args, dir.Name()) {
		args = append(args, `-i`, dir.Name())
	} else {
		panic(`must specify valid corpus directory`)
	}

	filePath, err := filepath.Abs(getOutputFile(`afl`))
	if err != nil {
		panic(err)
	}

	args = append(args, `--`, filePath)
	execute(true, `afl-fuzz`, args...)
}
