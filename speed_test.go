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

// nComp is the number of strings to compare all the others to.
const nComp = 1000

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

// Compare a number of long strings that have a substantial prefix in common.
func BenchmarkCompareSimilarStrings(b *testing.B) {
	// Create N mostly similar strings.
	if b.N < nComp {
		return // Nothing to do
	}
	strs := generateSimilarStrings(b.N)

	// Measure the time needed to compare each string to each of the first
	// nComp string.
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
	if b.N < nComp {
		return // Nothing to do
	}
	strs := generateRandomStrings(b.N)

	// Measure the time needed to compare every string to the first string.
	b.ResetTimer()
	for _, s1 := range strs {
		for _, s2 := range strs[:nComp] {
			if s1 == s2 {
				Dummy++
			}
		}
	}
}
