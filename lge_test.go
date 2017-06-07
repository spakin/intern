// This file includes unit tests for the LGE data type.

package intern_test

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/spakin/intern"
)

// TestPreLGEDups tests if we can create a large number of symbols for which
// duplicates are certain to occur.
func TestPreLGEDups(t *testing.T) {
	intern.ForgetAllLGE()
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
	intern.ForgetAllLGE()
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
	intern.ForgetAllLGE()
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
	intern.ForgetAllLGE()
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
			t.Error(err)
		}
	}

	// Ensure that converting an LGE back to a string is a lossless
	// operation.  We use fmt.Sprintf as this represents a typical way an
	// LGE might be converted to a string.
	for i, str := range strs {
		sym := syms[i]
		sStr := fmt.Sprintf("%s", sym)
		if str != sStr {
			t.Errorf("Expected %q but saw %q", str, sStr)
		}
	}
}

// TestBadLGE ensures we panic when converting an invalid LGE to a
// string.
func TestBadLGE(t *testing.T) {
	defer func() { _ = recover() }()
	var bad intern.LGE
	_ = bad.String() // Should panic
	t.Errorf("Failed to catch invalid intern.LGE %d", bad)
}
