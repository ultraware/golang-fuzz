package main

import (
	"testing"
	"unicode/utf8"
)

var inited bool

func init() {
	inited = true
}

func FuzzReverse(f *testing.F) {
	testcases := []string{"Hello, world", " ", "!12345"}
	for _, tc := range testcases {
		f.Add(tc) // Use f.Add to provide a seed corpus
	}
	f.Fuzz(func(t *testing.T, orig string) {
		if inited != true {
			panic(`not inited`)
		}

		rev, err1 := Reverse(orig)
		if err1 != nil {
			return //fuzz:rej
		}
		doubleRev, err2 := Reverse(rev)
		if err2 != nil {
			return
		}
		if !checkEqual(orig, doubleRev) {
			t.Errorf("Before: %q, after: %q", orig, doubleRev)
		}
		if utf8.ValidString(orig) && !utf8.ValidString(rev) {
			t.Errorf("Reverse produced invalid UTF-8 string %q", rev)
		}
		return //fuzz:prio
	})
}

func checkEqual(a, b string) bool {
	return a == b
}
