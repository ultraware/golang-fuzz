package main

import (
	"fmt"
	"go/ast"
	"os"
)

func main() {
	// Parse command line args
	args := parseArgs()

	// Find the Fuzz function in the package
	pkg, fname, fuzzFunc := findFuzzFunc(args[0], *funcName)
	if fuzzFunc == nil {
		fmt.Printf("Fuzz function %s not found in package %s\n", *funcName, args[0])
		os.Exit(1)
	}

	err := os.Chdir(args[0])
	if err != nil {
		panic(err)
	}

	params := fuzzFunc.Type.Params.List
	if len(params) != 1 { // TODO: Support multiple parameters
		fmt.Printf("Fuzz function %s must only have one parameter\n", fname)
		os.Exit(1)
	}

	build(pkg.Name, fname, fuzzFunc)
	if *runFuzzer {
		run(pkg.Name, fname, fuzzFunc, args[1:])
	}
}

func build(pkgName, fname string, fuzzFunc *ast.FuncDecl) {
	if *all || *native {
		fmt.Println("Generating Go native fuzzing test ...")
		generateGoNative(pkgName, fname, fuzzFunc)
	}
	if *all || *libfuzzer {
		fmt.Println("\nBuilding libFuzzer binarty ...")
		buildLibfFuzzer(pkgName, fname, fuzzFunc)
	}
	if *all || *gofuzz {
		if pkgName == `main` {
			fmt.Println("\nPackage main not supported by go-fuzz")
		} else {
			fmt.Println("\nBuilding go-fuzz binarty ...")
			buildGoFuzz(pkgName, fname, fuzzFunc)
		}
	}
	if *all || *afl {
		fmt.Println("\nBuilding AFL++ binary ...")
		buildAFL(pkgName, fname, fuzzFunc)
	}
}

func run(pkgName, fname string, fuzzFunc *ast.FuncDecl, args []string) {
	switch {
	case *native:
		fmt.Println("\nRunning Go native fuzzing test ...")
		runGoNative(pkgName, fname, fuzzFunc, args)
	case *libfuzzer:
		fmt.Println("\nRunning libFuzzer ...")
		runLibFuzzer(args)
	case *gofuzz && pkgName != `main`:
		fmt.Println("\nRunning go-fuzz ...")
		runGoFuzz(args)
	case *afl:
		fmt.Println("\nRunning AFL++ ...")
		runAFL(args)
	}
}
