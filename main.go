package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"os"
	"strings"
	"unicode"
)

var (
	funcName  = flag.String("func", "Fuzz", "name of the Fuzz function")
	corpusDir = flag.String("corpus", "corpus", "corpus directory for native Go fuzzing") // TODO: If not exists, show warning not adding any corpa
	keepFile  = flag.Bool("keep", false, "keep generated fuzz file")                      // TODO: Always true for native
)

func main() {
	// Parse command line flags
	flag.Parse()
	if flag.NArg() == 0 || *funcName == "" {
		fmt.Println("Usage: go-fuzz-build [options] PACKAGE_PATH")
		flag.PrintDefaults()
		os.Exit(1)
	}
	packagePath := flag.Args()[0]

	*funcName = strings.TrimPrefix(*funcName, `_`)
	if !strings.HasPrefix(*funcName, `Fuzz`) || (*funcName != `Fuzz` && !unicode.IsUpper(rune((*funcName)[4]))) {
		fmt.Printf("Fuzz function %s must be formatted as FuzzXxx\n", *funcName)
		os.Exit(1)
	}

	// Find the Fuzz function in the package
	pkg, fname, fuzzFunc := findFuzzFunc(packagePath, *funcName)
	if fuzzFunc == nil {
		fmt.Printf("Fuzz function %s not found in package %s\n", *funcName, packagePath)
		os.Exit(1)
	}
	err := os.Chdir(packagePath)
	if err != nil {
		panic(err)
	}

	prefixFuzzFuncName(fname, fuzzFunc)

	params := fuzzFunc.Type.Params.List
	if len(params) != 1 { // TODO: Support multiple parameters
		fmt.Printf("Fuzz function %s must only have one parameter\n", fname)
		os.Exit(1)
	}

	generateGoNative(pkg.Name, fname, fuzzFunc)
	// generateLibFuzzer(pkg.Name, fname, fuzzFunc)
}

func prefixFuzzFuncName(fname string, fuzzFunc *ast.FuncDecl) {
	if strings.HasPrefix(fuzzFunc.Name.Name, `_`) {
		fuzzFunc.Name.Name = strings.TrimPrefix(fuzzFunc.Name.Name, `_`)
		return
	}

	b, err := os.ReadFile(fname)
	if err != nil {
		panic(err)
	}

	b = bytes.Replace(b, []byte(`func `+fuzzFunc.Name.Name), []byte(`func _`+fuzzFunc.Name.Name), 1)
	err = os.WriteFile(fname, b, 0)
	if err != nil {
		panic(err)
	}
}
