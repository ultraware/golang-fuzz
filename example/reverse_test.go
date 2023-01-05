package main

import (
	"testing"
	"unicode/utf8"
)

func FuzzReverse(f *testing.F) {
	testcases := []string{"Hello, world", " ", "!12345"}
	for _, tc := range testcases {
		f.Add(tc, `abc`) // Use f.Add to provide a seed corpus
	}
	f.Fuzz(func(t *testing.T, orig string, aaaa string) {
		rev, err1 := Reverse(orig)
		if err1 != nil {
			return //fuzz:rej
		}
		doubleRev, err2 := Reverse(rev)
		if err2 != nil {
			return
		}
		if orig != doubleRev {
			t.Errorf("Before: %q, after: %q", orig, doubleRev)
		}
		if utf8.ValidString(orig) && !utf8.ValidString(rev) {
			t.Errorf("Reverse produced invalid UTF-8 string %q", rev)
		}
		return //fuzz:prio
		// data := []byte(orig + aaaa)
		// hash := sha1.Sum(data)
		// fileName := hex.EncodeToString(hash[:])

		// f, err := os.OpenFile(path.Join(`corpus`, fileName), os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o644)
		// if err != nil {
		// 	panic(err)
		// }
		// defer f.Close()

		// _, err = f.Write(data)
		// if err != nil {
		// 	panic(err)
		// }
	})
}
