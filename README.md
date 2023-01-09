# go-fuzz-build

A tool for generating/building fuzzing tests for various fuzzing engines without having to implement different types of Fuzz functions. The following fuzzing engines are supported:

- [Go native fuzzing engine](https://go.dev/security/fuzz/) (`go test -fuzz`)
- [dvyukov's go-fuzz](https://github.com/dvyukov/go-fuzz) using go-fuzz-build
- [LLVM libFuzzer](https://llvm.org/docs/LibFuzzer.html) using [go-libfuzz-build](https://github.com/elwint/go-libfuzz-build) (fork of [go114-fuzz-build](https://github.com/mdempsky/go114-fuzz-build))
- [AFL++](https://aflplus.plus/) using [go-afl-build](https://github.com/elwint/go-afl-build) (experimental)

## The Fuzz function

First, create an exported fuzzing function formatted as `FuzzXxx`. Unlike native fuzzing functions, this function must NOT be placed in a test file (`_test.go`). The Fuzz function must only have one parameter. The parameter type must be supported by the native Go fuzzing engine (<https://go.dev/security/fuzz/>). See the example folder for an implementation example.

## Usage

To use go-fuzz-build, run the following command:

```
go-fuzz-build [options] PACKAGE_PATH
```

Where `PACKAGE_PATH` is the path to the Go package containing the Fuzz function to be tested.

## Options

The following options are available:

- `-func`: the name of the Fuzz function (default: "Fuzz")
- `-corpus`: the corpus directory for native Go fuzzing (default: "corpus")
- `-keep`: keep generated fuzz file (always true for native)

Fuzzing engines:

- `-native`: generate native Go fuzzing test to run with go test -fuzz
- `-libfuzzer`: build for libFuzzer
- `-gofuzz`: build for go-fuzz
- `-afl`: build for AFL++
- `-all`: build for all supported fuzzing engines

## Notes

- Package main is not supported by go-fuzz.
