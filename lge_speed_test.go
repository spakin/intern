// This file measures the performance of operations on the LGE type.

package intern_test

import (
	"testing"

	"github.com/spakin/intern"
)

// BenchmarkLGECreation measures the time needed to create a symbol.
func BenchmarkLGECreation(b *testing.B) {
	strs := generateRandomStrings(b.N)
	syms := make([]intern.LGE, len(strs))
	b.ResetTimer()
	for _, s := range strs {
		intern.PreLGE(s)
	}
	var err error
	for i, s := range strs {
		syms[i], err = intern.NewLGE(s)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkCompareRandomLGEs compares a number of long, randomly generated
// strings by first mapping them to LGEs.
func BenchmarkCompareRandomLGEs(b *testing.B) {
	// Create N LGEs with random contents.
	if b.N < nComp {
		return // Nothing to do
	}
	strs := generateRandomStrings(b.N)
	syms := make([]intern.LGE, len(strs))
	intern.ForgetEverything()
	for _, s := range strs {
		intern.PreLGE(s)
	}
	var err error
	for i, s := range strs {
		syms[i], err = intern.NewLGE(s)
		if err != nil {
			b.Fatal(err)
		}
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

// BenchmarkCompareSimilarLGEs compares a number of long strings that have a
// substantial prefix in common by first mapping them to LGEs.
func BenchmarkCompareSimilarLGEs(b *testing.B) {
	// Create LGEs for N mostly similar strings.
	if b.N < nComp {
		return // Nothing to do
	}
	strs := generateSimilarStrings(b.N)
	syms := make([]intern.LGE, len(strs))
	intern.ForgetEverything()
	for _, s := range strs {
		intern.PreLGE(s)
	}
	var err error
	for i, s := range strs {
		syms[i], err = intern.NewLGE(s)
		if err != nil {
			b.Fatal(err)
		}
	}

	// Measure the time needed to compare each LGE to each of the first
	// nComp LGEs.
	b.ResetTimer()
	for _, s1 := range syms {
		for _, s2 := range syms[:nComp] {
			if s1 < s2 {
				Dummy++
			}
		}
	}
}
