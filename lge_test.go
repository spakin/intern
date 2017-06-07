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
