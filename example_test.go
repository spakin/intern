package intern_test

import (
	"fmt"

	"github.com/spakin/intern"
)

func ExamplePreLGE() {
	// Define some strings.  Note that these aren't unique.
	sList := []string{
		"Gunnar",
		"Högni",
		"Gjúki",
		"Gudrún",
		"Gotthorm",
		"Gjúki",
		"Óttar",
	}

	// Indicate our intention to intern all of the above.
	for _, s := range sList {
		intern.PreLGE(s)
	}

	// Intern all symbols into a set.  Report any duplicates
	// encountered.
	seen := make(map[intern.LGE]struct{}, len(sList))
	for _, s := range sList {
		sym, err := intern.NewLGE(s)
		if err != nil {
			panic(err)
		}
		if _, ok := seen[sym]; ok {
			fmt.Println(sym)
		}
		seen[sym] = struct{}{}
	}

	// Output:
	// Gjúki
}
