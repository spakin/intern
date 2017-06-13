// This file measures string performance as a baseline.

package intern_test

import (
	"math/rand"
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

// BenchmarkMergeStringMaps measures the performance of retrieving a number of
// strings from a map.
func BenchmarkMergeStringMaps(b *testing.B) {
	// Populate two maps.
	prng := rand.New(rand.NewSource(2223)) // Constant for reproducibility
	const sLen = 30                        // Symbol length in characters
	type Empty struct{}
	m1 := make(map[string]Empty, b.N)
	m2 := make(map[string]Empty, b.N)
	for i := 0; i < b.N; i++ {
		s := randomString(prng, sLen)
		m1[s] = Empty{}
		s = randomString(prng, sLen)
		m2[s] = Empty{}
	}

	// Start the clock then merge the two maps into a third.
	m3 := make(map[string]Empty, 2*b.N)
	b.ResetTimer()
	for k := range m1 {
		m3[k] = Empty{}
	}
	for k := range m2 {
		m3[k] = Empty{}
	}
}
