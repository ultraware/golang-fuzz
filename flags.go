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
	corpusDir  = flag.String("corpus", "corpus", "corpus directory (optional)")
	keepFile   = flag.Bool("keep", false, "keep generated fuzz file (always true for native)")
	printCmd   = flag.Bool("x", false, "print the commands")
	outputFile = flag.String("o", "", "output file")
	verbose    = flag.Bool("v", false, "show verbose output")
	runFuzzer  = flag.Bool("run", false, "run fuzzer after building")

	native = flag.Bool("native", false, "generate native Go fuzzing test to run with go test -fuzz")

	libfuzzer      = flag.Bool("libfuzzer", false, "build for libFuzzer")
	libfuzzerFlags = flag.String("libfuzzer.flags", "", "additional go-libfuzz-build flags")
	libfuzzerList  = flag.Bool("libfuzzer.list", false, "list go-libfuzz-build flags")
	clangFlags     = flag.String("libfuzzer.clangflags", "-g -O1 -fsanitize=fuzzer", "clang build flags")

	gofuzz      = flag.Bool("gofuzz", false, "build for go-fuzz")
	gofuzzFlags = flag.String("gofuzz.flags", "", "additional go-fuzz-build flags")
	gofuzzList  = flag.Bool("gofuzz.list", false, "list go-fuzz-build flags")

	afl      = flag.Bool("afl", false, "build for AFL++")
	aflFlags = flag.String("afl.flags", "", "additional go-afl-build flags")
	aflList  = flag.Bool("afl.list", false, "list go-afl-build flags")

	all = flag.Bool("all", false, "build for all supported fuzzing engines")
)

func parseArgs() []string {
	flag.Parse()

	if *libfuzzerList || *gofuzzList || *aflList {
		listBuildFlags()
		os.Exit(0)
	}

	checkOneFuzzerFlag(`-run`, *runFuzzer)
	checkOneFuzzerFlag(`-o`, *outputFile != ``)

	if *funcName == "" {
		fmt.Println("Usage: go-fuzz-build [options] PACKAGE_PATH [run_args]")
		flag.PrintDefaults()
		os.Exit(1)
	}

	// Check if Fuzz function is formatted FuzzXxx
	if !strings.HasPrefix(*funcName, `Fuzz`) || (*funcName != `Fuzz` && !unicode.IsUpper(rune((*funcName)[4]))) {
		fmt.Printf("Fuzz function %s must be formatted as FuzzXxx\n", *funcName)
		os.Exit(1)
	}

	if flag.NArg() == 0 {
		return []string{`.`}
	}
	return flag.Args()
}

func listBuildFlags() {
	*verbose = true
	if *libfuzzerList {
		command(`go-libfuzz-build`, `-help`)
	}
	if *gofuzzList {
		command(`go-fuzz-build`, `-help`)
	}
	if *aflList {
		command(`go-afl-build`, `-help`)
	}
}

func checkOneFuzzerFlag(flagName string, condition bool) {
	if *all && condition {
		fmt.Println(`Must specify a fuzzer when using ` + flagName)
		os.Exit(1)
	}
}

func getBuildArgs(flags *string, args ...string) []string {
	additionalFlags := strings.Fields(*flags)
	if len(additionalFlags) == 0 {
		return args
	}

	return append(additionalFlags, args...)
}

func getOutputFile(defaultFile string) string {
	if *outputFile != `` {
		return *outputFile
	}

	return defaultFile
}

func isValidCorpusDir() (bool, os.FileInfo) {
	f, err := os.Stat(*corpusDir)
	if err == nil && f.IsDir() {
		return true, f
	}

	return false, nil
}
