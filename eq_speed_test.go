// This file measures the performance of operations on the Eq type.

package intern_test

import (
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
	intern.ForgetEverything()
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
	intern.ForgetEverything()
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
