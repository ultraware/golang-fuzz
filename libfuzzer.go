package main

import (
	"go/ast"
	"os"
	"strings"

	_ "embed"

	"golang.org/x/exp/slices"
)

//go:embed tmpl/libfuzzer.go.tmpl
var tmplLibFuzzer string

func buildLibfFuzzer(pkgName string, fname string, fuzzFunc *ast.FuncDecl) {
	funcName, cleanup := generateLibFuzzer(pkgName, fname, fuzzFunc)
	defer cleanup()

	libFileName := createEmptyFile(`libfuzzer.*.a`)
	defer os.Remove(libFileName)                            //nolint: errcheck
	defer os.Remove(libFileName[:len(libFileName)-1] + `h`) //nolint: errcheck

	buildArgs := getBuildArgs(libfuzzerFlags, `-func`, funcName, `-o`, libFileName, `.`)
	command(`go-libfuzz-build`, buildArgs...)

	clangArgs := getBuildArgs(clangFlags, libFileName, `-o`, getOutputFile(`libfuzzer`))
	command(`clang`, clangArgs...)
}

func generateLibFuzzer(pkgName string, fname string, fuzzFunc *ast.FuncDecl) (string, func()) {
	funcName := fuzzFunc.Name.Name
	inputType := getInputType(fuzzFunc.Type.Params.List[0].Type)
	if inputType == `[]byte` { // generate not necessary, can just use original file
		return funcName, func() {}
	}

	cleanup := createTemplate(
		tmplLibFuzzer, strings.TrimSuffix(fname, `.go`)+`.*.go`,
		pkgName, funcName, inputType,
	)

	return funcName + `_`, cleanup
}

func createEmptyFile(pattern string) string {
	tmpFile, err := os.CreateTemp(`.`, pattern)
	if err == nil {
		err = tmpFile.Close()
	}
	if err != nil {
		panic(err)
	}

	return tmpFile.Name()
}

func runLibFuzzer(args []string) {
	valid, dir := isValidCorpusDir()
	if valid && !slices.Contains(args, dir.Name()) {
		args = append(args, dir.Name())
	}
	execute(false, getOutputFile(`libfuzzer`), args...)
}
