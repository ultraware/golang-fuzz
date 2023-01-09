package main

import (
	"fmt"
	"go/ast"
	"os"
	"os/exec"
	"strings"
	"text/template"

	_ "embed"
)

//go:embed tmpl/libfuzzer.go.tmpl
var tmplLibFuzzer string

func buildLibfFuzzer(pkgName string, fname string, fuzzFunc *ast.FuncDecl) {
	funcName := fuzzFunc.Name.Name
	inputType := getInputType(fuzzFunc.Type.Params.List[0].Type)
	if inputType != `[]byte` {
		funcName += `_`
		fuzzFile := generateLibFuzzer(pkgName, fname, fuzzFunc)
		if !*keepFile {
			defer os.Remove(fuzzFile)
		}
	}

	libFile, err := os.CreateTemp(`.`, `libfuzzer.*.a`)
	if err == nil {
		err = libFile.Close()
	}
	if err != nil {
		panic(err)
	}
	libFileName := libFile.Name()
	defer os.Remove(libFileName)
	defer os.Remove(libFileName[:len(libFileName)-1] + `h`)

	// fmt.Println(`go-libfuzz-build`, `-func`, funcName, `-o`, libFile.Name(), `.`)
	b, err := exec.Command(`go-libfuzz-build`, `-func`, funcName, `-o`, libFile.Name(), `.`).CombinedOutput()
	if err != nil {
		fmt.Println(string(b))
		panic(err)
	}

	b, err = exec.Command(`clang`, `-fsanitize=fuzzer`, libFile.Name(), `-o`, `libfuzzer`).CombinedOutput()
	if err != nil {
		fmt.Println(string(b))
		panic(err)
	}
}

// TODO: Fix dupl
func generateLibFuzzer(pkgName string, fname string, fuzzFunc *ast.FuncDecl) string {
	tmpl, err := template.New(``).Parse(tmplLibFuzzer)
	if err != nil {
		panic(err)
	}

	inputType := getInputType(fuzzFunc.Type.Params.List[0].Type)
	inputCode, inputLen, imprts := getInputCode(inputType)

	fuzzFile, err := os.CreateTemp(`.`, strings.TrimSuffix(fname, `.go`)+`.*.go`)
	if err != nil {
		panic(err)
	}
	defer fuzzFile.Close()

	err = tmpl.Execute(fuzzFile, tmplData{
		PkgName:   pkgName,
		Imports:   imprts,
		FuncName:  fuzzFunc.Name.Name,
		CorpusDir: *corpusDir,
		InputType: inputType,
		InputCode: inputCode,
		InputLen:  inputLen,
	})
	if err != nil {
		panic(err)
	}

	return fuzzFile.Name()
}
