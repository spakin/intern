// This file shows how to implement sort.Interface for LGE slices.

package intern_test

import (
	"fmt"
	"sort"

	"github.com/spakin/intern"
)

// LGESlice is a slice of LGEs that implements sort.Interface.
type LGESlice []intern.LGE

// Len returns the length of an LGESlice.
func (ls LGESlice) Len() int { return len(ls) }

// Less reports whether one element of an LGESlice is less than another.
func (ls LGESlice) Less(i, j int) bool { return ls[i] < ls[j] }

// Swap swaps two elements of an LGESLice.
func (ls LGESlice) Swap(i, j int) { ls[i], ls[j] = ls[j], ls[i] }

// Sort a list of strings by interning them to LGE symbols.
func ExamplePreLGE() {
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
	syms := make(LGESlice, len(sList))
	for i, s := range sList {
		l, err := intern.NewLGE(s)
		if err != nil {
			panic(err)
		}
		syms[i] = l
	}

	// Sort the LGESlice and output the result.
	sort.Sort(syms)
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
