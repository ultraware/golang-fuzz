package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"unicode"
)

var (
	funcName  = flag.String("func", "Fuzz", "name of the Fuzz function")
	corpusDir = flag.String("corpus", "corpus", "corpus directory for native Go fuzzing")
	keepFile  = flag.Bool("keep", false, "keep generated fuzz file (always true for native)")

	native    = flag.Bool("native", false, "generate native Go fuzzing test to run with go test -fuzz")
	libfuzzer = flag.Bool("libfuzzer", false, "build for libFuzzer")
	gofuzz    = flag.Bool("gofuzz", false, "build for go-fuzz")
	afl       = flag.Bool("afl", false, "build for AFL++")
	all       = flag.Bool("all", false, "build for all supported fuzzing engines")
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

	if pkg.Name == `main` && (*all || *gofuzz) {
		fmt.Println(`Package main not supported by go-fuzz`)
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
		generateGoNative(pkg.Name, fname, fuzzFunc)
	}
	if *all || *libfuzzer {
		buildLibfFuzzer(pkg.Name, fname, fuzzFunc)
	}
	if *all || *gofuzz {
		buildGoFuzz(pkg.Name, fname, fuzzFunc)
	}
	// if *all || *afl {
	// 	buildAFL(pkg.Name, fname, fuzzFunc)
	// }
}
