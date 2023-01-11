# go-fuzz-build

A tool for generating/building fuzzing tests for various fuzzing engines without having to implement different types of Fuzz functions. The following fuzzing engines are supported:

- [Go native fuzzing engine](https://go.dev/security/fuzz/) (`go test -fuzz`)
- [dvyukov's go-fuzz](https://github.com/dvyukov/go-fuzz) using go-fuzz-build
- [LLVM libFuzzer](https://llvm.org/docs/LibFuzzer.html) using [go-libfuzz-build](https://github.com/elwint/go-libfuzz-build) (fork of [go114-fuzz-build](https://github.com/mdempsky/go114-fuzz-build))
- [AFL++](https://aflplus.plus/) using [go-afl-build](https://github.com/elwint/go-afl-build) (experimental)

## The Fuzz function

In order to use `go-fuzz-build` to generate and build your fuzz tests, you will need to create an exported fuzzing function in your package. The name of this function should be formatted as `FuzzXxx`. 

It is important to note that, unlike a native fuzzing function, this function should not be placed in a test file (`_test.go`) and should only have one parameter. This parameter should be a type that is supported by the native Go fuzzing engine (<https://go.dev/security/fuzz/>).

A `Fuzz` function can be implemented as follows:

```go
func Fuzz(data []byte) int { // data can be any type supported by the native Go fuzzing engine
    // Your test logic goes here
    // ...

    /*
    The fuzz function should return an integer. This can be different than 0 to improve fuzzing performance.
    Returning 0 means that the input is accepted and may be added to the corpus.
    Returning -1 will cause libFuzzer or go-fuzz to not add that input to the corpus, regardless of coverage.
    Returning 1 will cause go-fuzz to increase priority of the given input.
    Fuzzing engines that do not support the returning value will treat it the same as returning 0.
    */
    return 0
}
```

See the example folder for a fuzzing test example.

## Usage

To use go-fuzz-build, run the following command:

```
go-fuzz-build [options] PACKAGE_PATH
```

Where `PACKAGE_PATH` is the path to the Go package containing the Fuzz function.

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
