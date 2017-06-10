// This file provides helper routines needed by multiple tests.

package intern_test

import (
	"fmt"
	"math/rand"
	"strings"
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

// Dummy is used to prevent benchmarks from being treated as dead code.
var Dummy uint64

// nComp is the number of strings to compare all the others to when benchmarking.
const nComp = 1000

// bMarkFunc is a string canonicalization function to use during benchmarking.
var bMarkFunc = strings.ToUpper

// generateSimilarStrings generates a list of strings that have a substantial
// prefix in common.
func generateSimilarStrings(n int) []string {
	strs := make([]string, n)
	for i := range strs {
		strs[i] = fmt.Sprintf("String comparisons can be slow when the strings to compare have a long prefix in common.  My favorite number is %15d.", i+1)
	}
	return strs
}

// generateRandomStrings generates a list of long strings with random content.
func generateRandomStrings(n int) []string {
	prng := rand.New(rand.NewSource(1920)) // Constant for reproducibility
	strs := make([]string, n)
	for i := range strs {
		nc := prng.Intn(50) + 10 // Number of characters
		strs[i] = randomString(prng, nc)
	}
	return strs
}
