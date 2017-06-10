// This file includes unit tests for the LGEC data type.

package intern_test

import (
	"fmt"
	"math/rand"
	"sort"
	"strings"
	"testing"

	"github.com/spakin/intern"
)

// TestPreLGECDups tests if we can create a large number of symbols for which
// duplicates are certain to occur.
func TestPreLGECDups(t *testing.T) {
	intern.ForgetAllLGEC()
	const sLen = 3                        // Symbol length in characters
	const nSymbols = 1000000              // Must be greater than len(charSet) choose sLen
	prng := rand.New(rand.NewSource(910)) // Constant for reproducibility
	for i := 0; i < nSymbols; i++ {
		intern.PreLGEC(randomString(prng, sLen), strings.ToTitle)
	}
	_, err := intern.NewLGEC("Yet another string", strings.ToTitle) // Force tree construction.
	if err != nil {
		t.Fatal(err)
	}
}

// TestPreLGECNoDups tests if we can create a large number of symbols for which
// duplicates are extremely unlikely to occur.
func TestPreLGECNoDups(t *testing.T) {
	intern.ForgetAllLGEC()
	const sLen = 50                        // Symbol length in characters
	const nSymbols = 100000                // Number of symbols to generate
	prng := rand.New(rand.NewSource(1112)) // Constant for reproducibility
	for i := 0; i < nSymbols; i++ {
		intern.PreLGEC(randomString(prng, sLen), strings.ToTitle)
	}
	_, err := intern.NewLGEC("Yet another string", strings.ToTitle) // Force tree construction.
	if err != nil {
		t.Fatal(err)
	}
}

// TestNewLGECFull tests that the tree does fill up and return an error if we
// don't use PreLGEC.
func TestNewLGECFull(t *testing.T) {
	// Creating 64 symbols in alphabetical order should work.
	intern.ForgetAllLGEC()
	var i int
	for i = 0; i < 64; i++ {
		str := fmt.Sprintf("This is symbol #%03d.", i+1)
		_, err := intern.NewLGEC(str, strings.ToUpper)
		if err != nil {
			t.Fatal(err)
		}
	}

	// Creating 65 symbols in alphabetical order should fail.
	str := fmt.Sprintf("This is symbol #%03d.", i)
	_, err := intern.NewLGEC(str, strings.ToUpper)
	if err == nil {
		t.Fatal("NewLGEC failed to return an error when its symbol table filled up")
	}
}

// TestLGECOrder ensures that LGEC symbol comparisons match the corresponding
// string comparisons.
func TestLGECOrder(t *testing.T) {
	// Create a bunch of random strings.
	intern.ForgetAllLGEC()
	f := strings.ToLower                   // Function to apply
	const sLen = 10                        // Symbol length in characters
	const nSymbols = 100                   // Number of symbols to generate
	prng := rand.New(rand.NewSource(1314)) // Constant for reproducibility
	strList := make([]string, nSymbols)
	for i := range strList {
		strList[i] = randomString(prng, sLen)
	}

	// Convert all of the strings to LGEC symbols.
	for _, str := range strList {
		intern.PreLGEC(str, f)
	}
	symList := make([]intern.LGEC, nSymbols)
	for i, str := range strList {
		var err error
		symList[i], err = intern.NewLGEC(str, f)
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
			case sym1 < sym2 && f(str1) < f(str2):
			case sym1 == sym2 && f(str1) == f(str2):
			case sym1 > sym2 && f(str1) > f(str2):
			default:
				t.Fatalf("Strings %q and %q mapped incorrectly to LGECs %d and %d", str1, str2, sym1, sym2)
			}
		}
	}
}

// TestLGECString tests if we can convert strings to LGECs and back to strings.
func TestLGECString(t *testing.T) {
	// Prepare the test.
	f := strings.ToUpper                   // Function to apply
	const ns = 10000                       // Number of strings to generate
	strs := make([]string, ns)             // Original strings
	syms := make([]intern.LGEC, ns)        // Interned strings
	prng := rand.New(rand.NewSource(1516)) // Constant for reproducibility

	// Generate a bunch of strings.
	for i := range strs {
		nc := prng.Intn(20) + 1 // Number of characters
		strs[i] = randomString(prng, nc)
	}

	// Intern each string to an LGEC.
	for _, s := range strs {
		intern.PreLGEC(s, f)
	}
	var err error
	for i, s := range strs {
		syms[i], err = intern.NewLGEC(s, f)
		if err != nil {
			t.Fatal(err)
		}
	}

	// Ensure that converting an LGEC back to a string is a lossless
	// operation.  We use fmt.Sprintf as this represents a typical way an
	// LGEC might be converted to a string.
	for i, str := range strs {
		sym := syms[i]
		sStr := fmt.Sprintf("%s", sym)
		if f(str) != f(sStr) {
			t.Fatalf("Expected %q but saw %q", str, sStr)
		}
	}
}

// TestBadLGEC ensures we panic when converting an invalid LGEC to a string.
func TestBadLGEC(t *testing.T) {
	defer func() { _ = recover() }()
	var bad intern.LGEC
	_ = bad.String() // Should panic
	t.Fatalf("Failed to catch invalid intern.LGEC %d", bad)
}

// TestLGECCase ensures that symbol comparisons are case-sensitive.
func TestLGECCase(t *testing.T) {
	// Convert a set of strings to LGECs.
	strs := []string{
		"roadrunner",
		"Roadrunner",
		"roadRunner",
		"ROADRUNNER",
		"rOaDrUnNeR",
		"ROADrunner",
		"roadRUNNER",
		"Coyote",
	}
	syms := make([]intern.LGEC, len(strs))
	var err error
	f := strings.ToUpper
	for i, s := range strs {
		syms[i], err = intern.NewLGEC(s, f)
		if err != nil {
			t.Fatal(err)
		}
	}

	// Ensure that each symbol is equal only to itself.
	numLGEC := 0
	for _, s1 := range syms {
		for _, s2 := range syms {
			if s1 == s2 {
				numLGEC++
			}
		}
	}
	expected := (len(syms)-1)*(len(syms)-1) + 1
	if numLGEC != expected {
		t.Fatalf("Expected %d case-insensitive comparisons but saw %d",
			expected, numLGEC)
	}
}

// TestSortLGECs generates a bunch of random LGECs and sorts them.
func TestSortLGECs(t *testing.T) {
	// Prepare the test.
	f := reverseString                     // Function to apply
	const ns = 1000                        // Number of strings to generate
	strs := make(sort.StringSlice, ns)     // Original strings
	syms := make(intern.LGECSlice, ns)     // Interned strings
	prng := rand.New(rand.NewSource(1718)) // Constant for reproducibility

	// Generate a bunch of strings.
	for i := range strs {
		nc := prng.Intn(20) + 1 // Number of characters
		strs[i] = randomString(prng, nc)
	}
	strs[5] = strs[ns-5] // Ensure at least one duplicate entry.

	// Intern each string to an LGEC.
	for _, s := range strs {
		intern.PreLGEC(s, f)
	}
	var err error
	for i, s := range strs {
		syms[i], err = intern.NewLGEC(s, f)
		if err != nil {
			t.Fatal(err)
		}
	}

	// Sort the list of LGECs and ensure that it's sorted.
	syms.Sort()
	for i, s := range syms[:syms.Len()-1] {
		if s > syms[i+1] {
			t.Fatalf("Symbols %q (%d) and %q (%d) are out of order",
				s, i, syms[i+1], i+1)
		}
	}

	// Sort the list of strings and ensure that it matches the
	// sorted list of LGECs.
	for i, str := range strs {
		strs[i] = f(str)
	}
	strs.Sort()
	for i, str := range strs {
		sym := syms[i]
		if f(str) != sym.String() {
			t.Fatalf("Sorted arrays don't match (%q != %q)", f(str), sym.String())
		}
	}
}
