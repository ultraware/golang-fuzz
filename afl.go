package main

import (
	"fmt"
	"go/ast"
	"os"
	"os/exec"

	_ "embed"
)

func buildAFL(pkgName string, fname string, fuzzFunc *ast.FuncDecl) {
	funcName := fuzzFunc.Name.Name
	inputType := getInputType(fuzzFunc.Type.Params.List[0].Type)
	if inputType != `[]byte` {
		funcName += `_`
		fuzzFile := generateLibFuzzer(pkgName, fname, fuzzFunc)
		if !*keepFile {
			defer os.Remove(fuzzFile)
		}
	}

	b, err := exec.Command(`go-afl-build`, `-func`, funcName).CombinedOutput()
	if err != nil {
		fmt.Println(string(b))
		panic(err)
	}
}
