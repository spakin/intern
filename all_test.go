// This file provides helper routines needed by multiple tests.

package intern_test

import (
	"math/rand"
	"unicode/utf8"
)

// charSet is the list of characters from which to draw for randomly generated
// strings.
var charSet []rune

// init initializes our global state.
func init() {
	const cs = "ÂBÇDÈFGHÍJKLMÑÖPQRSTÛVWXÝZ0123456789âbçdèfghíjklmñöpqrstûvwxýz @#$*-+<>一二三"
	charSet = make([]rune, 0, utf8.RuneCountInString(cs))
	for _, r := range cs {
		charSet = append(charSet, r)
	}
}

// randomString returns a random string of a given length.
func randomString(r *rand.Rand, n int) string {
	rs := make([]rune, n)
	nc := len(charSet)
	for i := range rs {
		rs[i] = charSet[r.Intn(nc)]
	}
	return string(rs)
}

// reverseString returns a string with its characters reversed.
func reverseString(s string) string {
	rs := []rune(s)
	nr := len(rs)
	for i := nr / 2; i >= 0; i-- {
		rs[i], rs[nr-i-1] = rs[nr-i-1], rs[i]
	}
	return string(rs)
}
