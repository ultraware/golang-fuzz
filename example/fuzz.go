package main

import (
	"fmt"
	"unicode/utf8"

	_ "embed"
)

var inited bool

func init() {
	inited = true
}

func Fuzz(orig string) int {
	if inited != true {
		panic(`not inited`)
	}

	rev, err1 := Reverse(orig)
	if err1 != nil {
		return -1
	}
	doubleRev, err2 := Reverse(rev)
	if err2 != nil {
		return 0
	}
	if !checkEqual(orig, doubleRev) {
		panic(fmt.Sprintf("Before: %q, after: %q", orig, doubleRev))
	}
	if utf8.ValidString(orig) && !utf8.ValidString(rev) {
		panic(fmt.Sprintf("Reverse produced invalid UTF-8 string %q", rev))
	}

	return 1
}

func checkEqual(a, b string) bool {
	return a == b
}
