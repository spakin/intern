// This file includes unit tests for the LGE data type.

package intern_test

import (
	"math/rand"
	"testing"

	"github.com/spakin/intern"
)

// TestPreLGE tests if we can create a large number of symbols for which
// duplicates are certain to occur.
func TestPreLGE(t *testing.T) {
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
