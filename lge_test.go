// This file provides unit tests for the LGE data type.

package intern_test

import (
	"fmt"
	"math/rand"
	"runtime"
	"testing"

	"github.com/spakin/intern"
)

// TestPreLGEDups tests if we can create a large number of symbols for which
// duplicates are certain to occur.
func TestPreLGEDups(t *testing.T) {
	intern.ForgetAllLGEs()
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
	intern.ForgetAllLGEs()
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

// TestNewLGEFull tests that the tree does fill up and return an error if we
// don't use PreLGE.
func TestNewLGEFull(t *testing.T) {
	// Creating 64 symbols in alphabetical order should work.
	intern.ForgetAllLGEs()
	var i int
	for i = 0; i < 64; i++ {
		str := fmt.Sprintf("This is symbol #%03d.", i+1)
		_, err := intern.NewLGE(str)
		if err != nil {
			t.Fatal(err)
		}
	}

	// Creating 65 symbols in alphabetical order should fail.
	str := fmt.Sprintf("This is symbol #%03d.", i)
	_, err := intern.NewLGE(str)
	if err == nil {
		t.Fatal("NewLGE failed to return an error when its symbol table filled up")
	}
}

// TestLGEOrder ensures that LGE symbol comparisons match the corresponding
// string comparisons.
func TestLGEOrder(t *testing.T) {
	// Create a bunch of random strings.
	intern.ForgetAllLGEs()
	const sLen = 10                        // Symbol length in characters
	const nSymbols = 100                   // Number of symbols to generate
	prng := rand.New(rand.NewSource(1314)) // Constant for reproducibility
	strList := make([]string, nSymbols)
	for i := range strList {
		strList[i] = randomString(prng, sLen)
	}

	// Convert all of the strings to LGE symbols.
	for _, str := range strList {
		intern.PreLGE(str)
	}
	symList := make([]intern.LGE, nSymbols)
	for i, str := range strList {
		var err error
		symList[i], err = intern.NewLGE(str)
		if err != nil {
			t.Fatal(err)
		}
	}

	// Compare all symbols.
	for i, sym1 := range symList {
		str1 := strList[i]
		for j, sym2 := range symList {
			str2 := strList[j]
			switch {
			case sym1 < sym2 && str1 < str2:
			case sym1 == sym2 && str1 == str2:
			case sym1 > sym2 && str1 > str2:
			default:
				t.Fatalf("Strings %q and %q mapped incorrectly to LGEs %d and %d", str1, str2, sym1, sym2)
			}
		}
	}
}

// TestLGEString tests if we can convert strings to LGEs and back to strings.
func TestLGEString(t *testing.T) {
	// Prepare the test.
	const ns = 10000                       // Number of strings to generate
	strs := make([]string, ns)             // Original strings
	syms := make([]intern.LGE, ns)         // Interned strings
	prng := rand.New(rand.NewSource(1516)) // Constant for reproducibility

	// Generate a bunch of strings.
	for i := range strs {
		nc := prng.Intn(20) + 1 // Number of characters
		strs[i] = randomString(prng, nc)
	}

	// Intern each string to an LGE.
	for _, s := range strs {
		intern.PreLGE(s)
	}
	var err error
	for i, s := range strs {
		syms[i], err = intern.NewLGE(s)
		if err != nil {
			t.Fatal(err)
		}
	}

	// Ensure that converting an LGE back to a string is a lossless
	// operation.  We use fmt.Sprintf as this represents a typical way an
	// LGE might be converted to a string.
	for i, str := range strs {
		sym := syms[i]
		sStr := fmt.Sprintf("%s", sym)
		if str != sStr {
			t.Fatalf("Expected %q but saw %q", str, sStr)
		}
	}
}

// TestLGEStringMulti tests if we can convert strings to LGEs and back to
// strings.  Unlike TestLGEString, it uses PreLGEs and NewLGEMulti.
func TestLGEStringMulti(t *testing.T) {
	// Prepare the test.
	const ns = 10000                       // Number of strings to generate
	strs := make([]string, ns)             // Original strings
	prng := rand.New(rand.NewSource(1516)) // Constant for reproducibility

	// Generate a bunch of strings.
	for i := range strs {
		nc := prng.Intn(20) + 1 // Number of characters
		strs[i] = randomString(prng, nc)
	}

	// Intern each string to an LGE.
	intern.PreLGEs(strs)
	syms, err := intern.NewLGEMulti(strs)
	if err != nil {
		t.Fatal(err)
	}

	// Ensure that converting an LGE back to a string is a lossless
	// operation.  We use fmt.Sprintf as this represents a typical way an
	// LGE might be converted to a string.
	for i, str := range strs {
		sym := syms[i]
		sStr := fmt.Sprintf("%s", sym)
		if str != sStr {
			t.Fatalf("Expected %q but saw %q", str, sStr)
		}
	}
}

// TestBadLGE ensures we panic when converting an invalid LGE to a string.
func TestBadLGE(t *testing.T) {
	defer func() { _ = recover() }()
	var bad intern.LGE
	_ = bad.String() // Should panic
	t.Fatalf("Failed to catch invalid intern.LGE %d", bad)
}

// TestForgetAllLGEs ensures we panic when converting a forgotten LGE to a
// string.
func TestForgetAllLGEs(t *testing.T) {
	defer func() { _ = recover() }()
	sym, err := intern.NewLGE("old string")
	if err != nil {
		t.Fatal(err)
	}
	str := sym.String()
	intern.ForgetAllLGEs()
	str = sym.String() // Should panic
	t.Fatalf("Failed to catch invalid intern.LGE %d (%q)", sym, str)
}

// TestLGECase ensures that symbol comparisons are case-sensitive.
func TestLGECase(t *testing.T) {
	// Convert a set of strings to LGEs.
	strs := []string{
		"roadrunner",
		"Roadrunner",
		"roadRunner",
		"ROADRUNNER",
		"rOaDrUnNeR",
		"ROADrunner",
		"roadRUNNER",
	}
	syms := make([]intern.LGE, len(strs))
	var err error
	for i, s := range strs {
		syms[i], err = intern.NewLGE(s)
		if err != nil {
			t.Fatal(err)
		}
	}

	// Ensure that each symbol is equal only to itself.
	numLGE := 0
	for _, s1 := range syms {
		for _, s2 := range syms {
			if s1 == s2 {
				numLGE++
			}
		}
	}
	if numLGE != len(syms) {
		t.Fatalf("Expected %d case-sensitive comparisons but saw %d",
			len(syms), numLGE)
	}
}

// TestLGEConcurrent performs a bunch of accesses in parallel in an attempt to
// expose race conditions.
func TestLGEConcurrent(t *testing.T) {
	const symsPerThread = 1000
	nThreads := runtime.NumCPU() * 2 // Oversubscribe CPUs by a factor of 2.

	// Spawn a number of goroutines.
	begin := make(chan bool, nThreads)
	done := make(chan bool, nThreads)
	for j := 0; j < nThreads; j++ {
		go func() {
			_ = <-begin
			prng := rand.New(rand.NewSource(2021)) // Constant for reproducibility and to invite conflicts
			for i := 0; i < symsPerThread; i++ {
				nc := prng.Intn(20) + 1 // Number of characters
				intern.PreLGE(randomString(prng, nc))
			}
			prng = rand.New(rand.NewSource(2021)) // Restart from the same seed.
			for i := 0; i < symsPerThread; i++ {
				nc := prng.Intn(20) + 1 // Number of characters
				_, err := intern.NewLGE(randomString(prng, nc))
				if err != nil {
					t.Fatal(err)
				}
			}
			done <- true
		}()
	}

	// Tell all goroutines to begin then wait for them all to finish.
	for j := 0; j < nThreads; j++ {
		begin <- true
	}
	for j := 0; j < nThreads; j++ {
		_ = <-done
	}
}

// TestRemapAllLGEs tests that old strings are remapped and pending strings are
// added.
func TestRemapAllLGEs(t *testing.T) {
	// Define a list of strings and allocate space for two lists of
	// associated LGEs.
	strs := []string{
		"David J. Thouless",
		"F. Duncan M. Haldane",
		"J. Michael Kosterlitz",
		"Jean-Pierre Sauvage",
		"Sir J. Fraser Stoddart",
		"Bernard L. Feringa",
		"Yoshinori Ohsumi",
		"Bob Dylan",
		"Juan Manuel Santos",
		"Oliver Hart",
		"Bengt Holmström",
	}
	nStrs := len(strs)
	syms0 := make([]intern.LGE, nStrs/2)
	syms1 := make([]intern.LGE, nStrs)

	// Generate symbols for half the names.
	for i := 0; i < nStrs/2; i++ {
		intern.PreLGE(strs[i])
	}
	var err error
	for i := 0; i < nStrs/2; i++ {
		syms0[i], err = intern.NewLGE(strs[i])
		if err != nil {
			t.Error(err)
		}
	}

	// Preallocate the remaining symbols (with a bit of overlap to test
	// that, too), but don't actually allocate them.
	for i := nStrs/2 - 1; i < nStrs; i++ {
		intern.PreLGE(strs[i])
	}

	// Remap all old and pending symbols.
	m, err := intern.RemapAllLGEs()
	if err != nil {
		t.Error(err)
	}
	for i, s := range strs {
		syms1[i], err = intern.NewLGE(s)
		if err != nil {
			t.Error(err)
		}
	}

	// Confirm that the map is accurate.
	for i, s := range strs[:nStrs/2] {
		s0 := syms0[i]
		s1 := syms1[i]
		if m[s0] != s1 {
			t.Errorf("For %q, expected %d to map to %d but it instead mapped to %d",
				s, s0, s1, m[s0])
		}
	}

	// Confirm that all new LGEs compare with each other as expected.
	for i, a0 := range strs {
		b0 := syms1[i]
		for j, a1 := range strs {
			b1 := syms1[j]
			switch {
			case a0 == a1 && b0 == b1:
			case a0 < a1 && b0 < b1:
			case a0 > a1 && b0 > b1:
			default:
				t.Errorf("Comparison of strings %q and %q does not match comparison of LGEs %d and %d",
					a0, a1, b0, b1)
			}
		}
	}
}
