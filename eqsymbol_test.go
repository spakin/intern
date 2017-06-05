// This file includes unit tests for the EqSymbol data type.

package intern_test

import (
	"math/rand"
	"testing"
	"unicode/utf8"

	"github.com/spakin/intern"
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

// TestNewEqSymbol tests if we can create a large number of symbols for which
// duplicates are certain to occur.
func TestNewEqSymbol(*testing.T) {
	const sLen = 3                         // Symbol length in characters
	const nSymbols = 1000000               // Must be greater than len(charSet) choose sLen
	prng := rand.New(rand.NewSource(1234)) // Constant for reproducibility
	for i := 0; i < nSymbols; i++ {
		_ = intern.NewEqSymbol(randomString(prng, sLen))
	}
}
