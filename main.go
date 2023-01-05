package main

import (
	"flag"
	"fmt"
	"go/ast"
	"os"
)

var (
	funcName  = flag.String("func", "", "name of the Fuzz function")
	genCorpus = flag.Bool("corpus", true, "generate corpus dir")
	keepFile  = flag.Bool("keep", false, "keep generated fuzz file")
)

func main() {
	// Parse command line flags
	flag.Parse()
	if flag.NArg() == 0 || *funcName == "" {
		fmt.Println("Usage: go-fuzz-build -func FUNC_NAME PACKAGE_PATH")
		os.Exit(1)
	}
	packagePath := flag.Args()[0]

	// Find the Fuzz function in the package
	pkg, fname, sourceCode, fuzzFunc := findFuzzFunc(packagePath, *funcName)
	if fuzzFunc == nil {
		fmt.Printf("Fuzz function %s not found in package %s\n", *funcName, packagePath)
		os.Exit(1)
	}
	err := os.Chdir(packagePath)
	if err != nil {
		panic(err)
	}

	// Find the function passed to f.Fuzz
	fuzzCall := findFuzzCall(fuzzFunc.Body)
	if fuzzCall == nil {
		fmt.Printf("Fuzz function %s does not contain a call to f.Fuzz\n", *funcName)
		os.Exit(1)
	}

	fFuzz, ok := fuzzCall.Args[0].(*ast.FuncLit)
	if !ok {
		fmt.Printf("Expected function passed to f.Fuzz, got %T\n", fuzzCall.Args[0])
		os.Exit(1)
	}

	params := fFuzz.Type.Params.List
	if len(params) != 2 { // TODO: Support multiple parameters
		fmt.Printf("Function passed to f.Fuzz must only have one parameter")
		os.Exit(1)
	}

	fmt.Printf("Function passed to f.Fuzz has %d parameters:\n", len(params)-1)
	// Print the types of the parameters of the passed function
	for i := 1; i < len(params); i++ {
		fmt.Printf("- Parameter %d: %s\n", i, params[i].Type) // TODO: Check supported types
	}

	generateLibFuzzer(sourceCode, fFuzz)

	if *genCorpus {
		captureCorpusValues(pkg, fname, sourceCode, fFuzz)
	}
}
