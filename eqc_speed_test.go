// This file measures the performance of operations on the EqC type.

package intern_test

import (
	"testing"

	"github.com/spakin/intern"
)

// BenchmarkEqCCreation measures the time needed to create a symbol.
func BenchmarkEqCCreation(b *testing.B) {
	strs := generateRandomStrings(b.N)
	syms := make([]intern.EqC, len(strs))
	b.ResetTimer()
	for i, s := range strs {
		syms[i] = intern.NewEqC(s, bMarkFunc)
	}
}

// BenchmarkCompareRandomEqCs compares a number of long, randomly generated
// strings by first mapping them to EqCs.
func BenchmarkCompareRandomEqCs(b *testing.B) {
	// Create N EqCs with random contents.
	if b.N < nComp {
		return // Nothing to do
	}
	strs := generateRandomStrings(b.N)
	syms := make([]intern.EqC, len(strs))
	intern.ForgetAllEqC()
	for i, s := range strs {
		syms[i] = intern.NewEqC(s, bMarkFunc)
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

// BenchmarkCompareSimilarEqCs compares a number of long strings that have a
// substantial prefix in common by first mapping them to EqCs.
func BenchmarkCompareSimilarEqCs(b *testing.B) {
	// Create EqCs for N mostly similar strings.
	if b.N < nComp {
		return // Nothing to do
	}
	strs := generateSimilarStrings(b.N)
	syms := make([]intern.EqC, len(strs))
	intern.ForgetAllEqC()
	for i, s := range strs {
		syms[i] = intern.NewEqC(s, bMarkFunc)
	}

	// Measure the time needed to compare each EqC to each of the first
	// nComp EqCs.
	b.ResetTimer()
	for _, s1 := range syms {
		for _, s2 := range syms[:nComp] {
			if s1 == s2 {
				Dummy++
			}
		}
	}
}
