// This file measures the performance of operations on the Eq type.

package intern_test

import (
	"math/rand"
	"testing"

	"github.com/spakin/intern"
)

// BenchmarkEqCreation measures the time needed to create a symbol.
func BenchmarkEqCreation(b *testing.B) {
	strs := generateRandomStrings(b.N)
	syms := make([]intern.Eq, len(strs))
	b.ResetTimer()
	for i, s := range strs {
		syms[i] = intern.NewEq(s)
	}
}

// BenchmarkCompareRandomEqs compares a number of long, randomly generated
// strings by first mapping them to Eqs.
func BenchmarkCompareRandomEqs(b *testing.B) {
	// Create N Eqs with random contents.
	if b.N < nComp {
		return // Nothing to do
	}
	strs := generateRandomStrings(b.N)
	syms := make([]intern.Eq, len(strs))
	intern.ForgetAllEqs()
	for i, s := range strs {
		syms[i] = intern.NewEq(s)
	}

	// Measure the time needed to compare each string to each of the first
	// nComp strings.
	b.ResetTimer()
	for _, s1 := range syms {
		for _, s2 := range syms[:nComp] {
			if s1 == s2 {
				Dummy++
			}
		}
	}
}

// BenchmarkCompareSimilarEqs compares a number of long strings that have a
// substantial prefix in common by first mapping them to Eqs.
func BenchmarkCompareSimilarEqs(b *testing.B) {
	// Create Eqs for N mostly similar strings.
	if b.N < nComp {
		return // Nothing to do
	}
	strs := generateSimilarStrings(b.N)
	syms := make([]intern.Eq, len(strs))
	intern.ForgetAllEqs()
	for i, s := range strs {
		syms[i] = intern.NewEq(s)
	}

	// Measure the time needed to compare each Eq to each of the first
	// nComp Eqs.
	b.ResetTimer()
	for _, s1 := range syms {
		for _, s2 := range syms[:nComp] {
			if s1 == s2 {
				Dummy++
			}
		}
	}
}

// BenchmarkMergeEqMaps measures the performance of retrieving a number of Eqs
// from a map.
func BenchmarkMergeEqMaps(b *testing.B) {
	// Populate two maps.
	prng := rand.New(rand.NewSource(2223)) // Constant for reproducibility
	const sLen = 20                        // Symbol length in characters
	type Empty struct{}
	m1 := make(map[intern.Eq]Empty, b.N)
	m2 := make(map[intern.Eq]Empty, b.N)
	for i := 0; i < b.N; i++ {
		s := randomString(prng, sLen)
		m1[intern.NewEq(s)] = Empty{}
		s = randomString(prng, sLen)
		m2[intern.NewEq(s)] = Empty{}
	}

	// Start the clock then merge the two maps into a third.
	m3 := make(map[intern.Eq]Empty, 2*b.N)
	b.ResetTimer()
	for k := range m1 {
		m3[k] = Empty{}
	}
	for k := range m2 {
		m3[k] = Empty{}
	}
}
