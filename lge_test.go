// This file includes unit tests for the LGE data type.

package intern_test

import (
	"math/rand"
	"testing"

	"github.com/spakin/intern"
)

// TestPreLGEDups tests if we can create a large number of symbols for which
// duplicates are certain to occur.
func TestPreLGEDups(t *testing.T) {
	intern.ForgetAllLGE()
	const sLen = 3                        // Symbol length in characters
	const nSymbols = 1000000              // Must be greater than len(charSet) choose sLen
	prng := rand.New(rand.NewSource(910)) // Constant for reproducibility
	for i := 0; i < nSymbols; i++ {
		intern.PreLGE(randomString(prng, sLen))
	}
	_, err := intern.NewLGE("Yet another string") // Force tree construction.
	if err != nil {
		t.Fatal(err)
	}
}

// TestPreLGENoDups tests if we can create a large number of symbols for which
// duplicates are extremely unlikely to occur.
func TestPreLGENoDups(t *testing.T) {
	intern.ForgetAllLGE()
	const sLen = 50                        // Symbol length in characters
	const nSymbols = 100000                // Number of symbols to generate
	prng := rand.New(rand.NewSource(1112)) // Constant for reproducibility
	for i := 0; i < nSymbols; i++ {
		intern.PreLGE(randomString(prng, sLen))
	}
	_, err := intern.NewLGE("Yet another string") // Force tree construction.
	if err != nil {
		t.Fatal(err)
	}
}
