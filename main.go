package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"unicode"
)

var (
	funcName   = flag.String("func", "Fuzz", "name of the Fuzz function")
	corpusDir  = flag.String("corpus", "corpus", "corpus directory for native Go fuzzing")
	keepFile   = flag.Bool("keep", false, "keep generated fuzz file (always true for native)")
	printCmd   = flag.Bool("x", false, "print the commands")
	outputFile = flag.String("o", "", "output file")

	native    = flag.Bool("native", false, "generate native Go fuzzing test to run with go test -fuzz")
	libfuzzer = flag.Bool("libfuzzer", false, "build for libFuzzer")
	gofuzz    = flag.Bool("gofuzz", false, "build for go-fuzz")
	afl       = flag.Bool("afl", false, "build for AFL++")
	all       = flag.Bool("all", false, "build for all supported fuzzing engines")
)

func main() {
	// Parse command line args
	packagePath := parseArgs()

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

	params := fuzzFunc.Type.Params.List
	if len(params) != 1 { // TODO: Support multiple parameters
		fmt.Printf("Fuzz function %s must only have one parameter\n", fname)
		os.Exit(1)
	}

	if *all || *native {
		fmt.Println("Generating Go native fuzzing test ...")
		generateGoNative(pkg.Name, fname, fuzzFunc)
	}
	if *all || *libfuzzer {
		fmt.Println("\nBuilding libFuzzer binarty ...")
		buildLibfFuzzer(pkg.Name, fname, fuzzFunc)
	}
	if *all || *gofuzz {
		if pkg.Name == `main` {
			fmt.Println("\nPackage main not supported by go-fuzz")
		} else {
			fmt.Println("\nBuilding go-fuzz binarty ...")
			buildGoFuzz(pkg.Name, fname, fuzzFunc)
		}
	}
	if *all || *afl {
		fmt.Println("\nBuilding AFL++ binary ...")
		buildAFL(pkg.Name, fname, fuzzFunc)
	}
}

func parseArgs() string {
	flag.Parse()
	if *funcName == "" {
		fmt.Println("Usage: go-fuzz-build [options] PACKAGE_PATH")
		flag.PrintDefaults()
		os.Exit(1)
	}

	// Check if Fuzz function is formatted FuzzXxx
	if !strings.HasPrefix(*funcName, `Fuzz`) || (*funcName != `Fuzz` && !unicode.IsUpper(rune((*funcName)[4]))) {
		fmt.Printf("Fuzz function %s must be formatted as FuzzXxx\n", *funcName)
		os.Exit(1)
	}

	if *outputFile != `` && *all {
		fmt.Println(`Must specify a fuzzer when using -o`)
		os.Exit(1)
	}

	if flag.NArg() == 0 {
		return `.`
	}
	return flag.Args()[0]
}
