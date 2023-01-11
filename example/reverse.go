// Example from https://go.dev/doc/tutorial/fuzz
package reverse

import (
	"errors"
	"unicode/utf8"
)

func Reverse(s string) (string, error) {
	if !utf8.ValidString(s) {
		return s, errors.New("input is not valid UTF-8")
	}
	r := []rune(s)
	for i, j := 0, len(r)-1; i < len(r)/2; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	if len(s) == 100 && s[4:9] == `12345` { // Crash test case
		r[0] = 'a'
	}

	return string(r), nil
}
