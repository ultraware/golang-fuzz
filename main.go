package main

import (
	"flag"
	"fmt"
	"os"
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
		fmt.Println("Usage: go-fuzz-build -func FUNC_NAME PACKAGE_PATH")
		os.Exit(1)
	}
	packagePath := flag.Args()[0]

	// Find the Fuzz function in the package
	pkg, fname, _, fuzzFunc := findFuzzFunc(packagePath, *funcName)
	if fuzzFunc == nil {
		fmt.Printf("Fuzz function %s not found in package %s\n", *funcName, packagePath)
		os.Exit(1)
	}
	err := os.Chdir(packagePath)
	if err != nil {
		panic(err)
	}

	params := fuzzFunc.Type.Params.List
	if len(params) != 1 { // TODO: Support multiple parameters
		fmt.Printf("Fuzz function %s must only have one parameter\n", fname)
		os.Exit(1)
	}

	// TODO: Check supported types?
	generateGoNative(pkg.Name, fname, fuzzFunc)
}
