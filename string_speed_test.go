// This file measures string performance as a baseline.

package intern_test

import (
	"testing"
)

// BenchmarkCompareRandomStrings compares a number of long, randomly generated
// strings.
func BenchmarkCompareRandomStrings(b *testing.B) {
	// Create N strings with random contents.
	if b.N < nComp {
		return // Nothing to do
	}
	strs := generateRandomStrings(b.N)

	// Measure the time needed to compare each string to each of the first
	// nComp strings.
	b.ResetTimer()
	for _, s1 := range strs {
		for _, s2 := range strs[:nComp] {
			if s1 == s2 {
				Dummy++
			}
		}
	}
}

// BenchmarkCompareSimilarStrings compares a number of long strings that have a
// substantial prefix in common.
func BenchmarkCompareSimilarStrings(b *testing.B) {
	// Create N mostly similar strings.
	if b.N < nComp {
		return // Nothing to do
	}
	strs := generateSimilarStrings(b.N)

	// Measure the time needed to compare each string to each of the first
	// nComp strings.
	b.ResetTimer()
	for _, s1 := range strs {
		for _, s2 := range strs[:nComp] {
			if s1 == s2 {
				Dummy++
			}
		}
	}
}
