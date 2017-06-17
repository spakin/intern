package intern_test

import (
	"bufio"
	"fmt"
	"io"
	"os"

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

// Maintain a long list of symbols, remapping as necessary.
func ExampleRemapAllLGEs() {
	syms := make([]intern.LGE, 0, 10)
	rb := bufio.NewReader(os.Stdin)
	for {
		// Read a line from standard input.
		s, err := rb.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}
		s = s[:len(s)-1]

		// Map the symbol if we can.  For this example,
		// pretend we need the symbol right away and therefore
		// wouldn't benefit by pre-allocating it with
		// intern.PreLGE.
		sy, err := intern.NewLGE(s)
		if err == nil {
			syms = append(syms, sy)
			continue
		}

		// The LGE symbol table is full.  Remap all existing
		// symbols and try again.
		intern.PreLGE(s)
		m, err := intern.RemapAllLGEs()
		if err != nil {
			panic(err)
		}
		for i, sy := range syms {
			syms[i] = m[sy]
		}
		sy, err = intern.NewLGE(s)
		if err != nil {
			panic(err)
		}
		syms = append(syms, sy)
	}
	fmt.Printf("%v\n", syms)
}
