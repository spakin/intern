// This file provides various performance tests for the intern package.

package intern_test

import (
	"fmt"
	"math/rand"
	"testing"
	//"github.com/spakin/intern"
)

// Dummy is used to prevent benchmarks from being treated as dead code.
var Dummy uint64

// Compare a number of long strings that have a substantial prefix in common.
func BenchmarkCompareSimilarStrings(b *testing.B) {
	// Create N mostly similar strings.
	const nComp = 1000 // Number of strings to compare all the others to
	if b.N < nComp {
		return // Nothing to do
	}
	strs := make([]string, b.N)
	for i := range strs {
		strs[i] = fmt.Sprintf("String comparisons can be slow when the strings to compare have a long prefix in common.  My favorite number is %15d.", i+1)
	}

	// Measure the time neede to compare every string to the first string.
	b.ResetTimer()
	for _, s1 := range strs {
		for _, s2 := range strs[:nComp] {
			if s1 == s2 {
				Dummy++
			}
		}
	}
}

// Compare a number of long, randomly generated strings.
func BenchmarkCompareRandomStrings(b *testing.B) {
	// Create N strings with random contents.
	prng := rand.New(rand.NewSource(1920)) // Constant for reproducibility
	const nComp = 1000                     // Number of strings to compare all the others to
	if b.N < nComp {
		return // Nothing to do
	}
	strs := make([]string, b.N)
	for i := range strs {
		nc := prng.Intn(50) + 10 // Number of characters
		strs[i] = randomString(prng, nc)
	}

	// Measure the time neede to compare every string to the first string.
	b.ResetTimer()
	for _, s1 := range strs {
		for _, s2 := range strs[:nComp] {
			if s1 == s2 {
				Dummy++
			}
		}
	}
}
