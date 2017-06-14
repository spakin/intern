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

// Sort a list of strings by interning them to LGE symbols.
func ExampleLGESlice_Sort() {
	// Define some strings.
	sList := []string{
		"Gerontius",
		"Reginard",
		"Hildigrim",
		"Eglantine",
		"Diamond",
		"Adamanta",
		"Sigismond",
		"Adalgrim",
		"Flambard",
		"Paladin",
		"Peregrin",
		"Pimpernel",
		"Everard",
		"Ferdibrand",
		"Pervinca",
		"Lalia",
		"Ferdinand",
		"Isembard",
		"Isembold",
		"Hildifons",
		"Isengrim",
		"Faramir",
		"Isengar",
		"Pearl",
		"Goldilocks",
		"Fortinbras",
		"Isumbras",
		"Bandobras",
		"Adelard",
		"Hildigard",
		"Hildibrand",
		"Rosa",
	}

	// Indicate our intent to intern all of the strings.
	for _, s := range sList {
		intern.PreLGE(s)
	}

	// Intern each string into an LGE and store it in an LGESlice.
	syms := make(intern.LGESlice, len(sList))
	for i, s := range sList {
		l, err := intern.NewLGE(s)
		if err != nil {
			panic(err)
		}
		syms[i] = l
	}

	// Sort the LGESlice and output the result.
	syms.Sort()
	for _, s := range syms {
		fmt.Println(s)
	}

	// Output:
	// Adalgrim
	// Adamanta
	// Adelard
	// Bandobras
	// Diamond
	// Eglantine
	// Everard
	// Faramir
	// Ferdibrand
	// Ferdinand
	// Flambard
	// Fortinbras
	// Gerontius
	// Goldilocks
	// Hildibrand
	// Hildifons
	// Hildigard
	// Hildigrim
	// Isembard
	// Isembold
	// Isengar
	// Isengrim
	// Isumbras
	// Lalia
	// Paladin
	// Pearl
	// Peregrin
	// Pervinca
	// Pimpernel
	// Reginard
	// Rosa
	// Sigismond
}
