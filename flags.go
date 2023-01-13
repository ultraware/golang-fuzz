package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"unicode"
)

func parseArgs() []string {
	flag.Parse()

	if *listFlags {
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
	checkOneFuzzerFlag(`-listflags`, true)

	*verbose = true
	if *libfuzzer {
		command(`go-libfuzz-build`, `-help`)
	}
	if *gofuzz {
		command(`go-fuzz-build`, `-help`)
	}
	if *afl {
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
