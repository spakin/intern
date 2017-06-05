// This file includes unit tests for the EqC data type.  It reuses some of the
// global state defined in eq_test.go

package intern_test

import (
	"fmt"
	"math/rand"
	"strings"
	"testing"

	"github.com/spakin/intern"
)

// TestNewEqC tests if we can create a large number of symbols for which
// duplicates are certain to occur.
func TestNewEqC(*testing.T) {
	const sLen = 3                       // Symbol length in characters
	const nSymbols = 1000000             // Must be greater than len(charSet) choose sLen
	prng := rand.New(rand.NewSource(56)) // Constant for reproducibility
	for i := 0; i < nSymbols; i++ {
		_ = intern.NewEqC(randomString(prng, sLen), strings.ToTitle)
	}
}

// TestEqCString tests if we can convert strings to symbols and back to
// strings.
func TestEqCString(t *testing.T) {
	// Prepare the test.
	const ns = 10000                     // Number of strings to generate
	strs := make([]string, ns)           // Original strings
	syms := make([]intern.EqC, ns)       // Interned strings
	prng := rand.New(rand.NewSource(78)) // Constant for reproducibility

	// Generate a bunch of strings.
	for i := range strs {
		nc := prng.Intn(20) + 1 // Number of characters
		strs[i] = randomString(prng, nc)
	}

	// Intern each string to an EqC.
	for i, s := range strs {
		syms[i] = intern.NewEqC(s, strings.ToLower)
	}

	// Ensure that converting an EqC back to a string is a lossless
	// operation.  We use fmt.Sprintf as this represents a typical way an
	// EqC might be converted to a string.
	for i, str := range strs {
		sym := syms[i]
		lstr := strings.ToLower(str)
		lsym := strings.ToLower(fmt.Sprintf("%s", sym))
		if lstr != lsym {
			t.Errorf("Expected %q but saw %q", lstr, lsym)
		}
	}
}

// TestBadEqC ensures we panic when converting an invalid EqC to a
// string.
func TestBadEqC(t *testing.T) {
	defer func() { _ = recover() }()
	var bad intern.EqC
	_ = bad.String() // Should panic
	t.Errorf("Failed to catch invalid intern.EqC %d", bad)
}

// TestEqCCase ensures that symbol comparisons are case-sensitive when used
// with strings.ToUpper.
func TestEqCCase(t *testing.T) {
	// Convert a set of strings to EqCs.
	strs := []string{
		"roadrunner",
		"Roadrunner",
		"roadRunner",
		"ROADRUNNER",
		"rOaDrUnNeR",
		"ROADrunner",
		"roadRUNNER",
	}
	syms := make([]intern.EqC, len(strs))
	for i, s := range strs {
		syms[i] = intern.NewEqC(s, strings.ToUpper)
	}

	// Ensure that each symbol is equal to all the other symbols.
	numEqC := 0
	for _, s1 := range syms {
		for _, s2 := range syms {
			if s1 == s2 {
				numEqC++
			}
		}
	}
	if numEqC != len(syms)*len(syms) {
		t.Errorf("Expected %d case-insensitive comparisons but saw %d",
			len(syms)*len(syms), numEqC)
	}
}
