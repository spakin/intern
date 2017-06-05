// This file includes unit tests for the EqSymbol data type.

package intern_test

import (
	"fmt"
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
	const sLen = 3                       // Symbol length in characters
	const nSymbols = 1000000             // Must be greater than len(charSet) choose sLen
	prng := rand.New(rand.NewSource(12)) // Constant for reproducibility
	for i := 0; i < nSymbols; i++ {
		_ = intern.NewEqSymbol(randomString(prng, sLen))
	}
}

// TestEqSymbolString tests if we can convert strings to symbols and back to
// strings.
func TestEqSymbolString(t *testing.T) {
	// Prepare the test.
	const ns = 10000                     // Number of strings to generate
	strs := make([]string, ns)           // Original strings
	syms := make([]intern.EqSymbol, ns)  // Interned strings
	prng := rand.New(rand.NewSource(34)) // Constant for reproducibility

	// Generate a bunch of strings.
	for i := range strs {
		nc := prng.Intn(20) + 1 // Number of characters
		strs[i] = randomString(prng, nc)
	}

	// Intern each string to an EqSymbol.
	for i, s := range strs {
		syms[i] = intern.NewEqSymbol(s)
	}

	// Ensure that converting an EqSymbol back to a string is a lossless
	// operation.  We use fmt.Sprintf as this represents a typical way an
	// EqSymbol might be converted to a string.
	for i, str := range strs {
		sym := syms[i]
		if str != fmt.Sprintf("%s", sym) {
			t.Errorf("Expected %q but saw %q", str, sym)
		}
	}
}

// TestBadEqSymbol ensures we panic when converting an invalid EqSymbol to a
// string.
func TestBadEqSymbol(t *testing.T) {
	defer func() { _ = recover() }()
	var bad intern.EqSymbol
	_ = bad.String() // Should panic
	t.Errorf("Failed to catch invalid EqSymbol %d", bad)
}
