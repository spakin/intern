package intern_test

import (
	"fmt"

	"github.com/spakin/intern"
)

// Find all duplicates in a list of strings.
func ExampleNewEq() {
	// Define some strings.  Note that these aren't unique.
	sList := []string{
		"Gunnar",
		"Högni",
		"Gjúki",
		"Gudrún",
		"Gotthorm",
		"Gjúki",
		"Óttar",
		"Sigurd",
		"Svanhild",
		"Jörmunrek",
		"Jónakr",
		"Hamdir",
		"Sörli",
		"Jónakr",
		"Kostbera",
		"Snævar",
		"Atli",
	}

	// Intern all symbols into a set.  Report any duplicates
	// encountered.
	seen := make(map[intern.Eq]struct{}, len(sList))
	for _, s := range sList {
		sym := intern.NewEq(s)
		if _, ok := seen[sym]; ok {
			fmt.Println(sym)
		}
		seen[sym] = struct{}{}
	}

	// Output:
	// Gjúki
	// Jónakr
}
