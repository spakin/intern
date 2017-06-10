// This file provides the LGE type, which represents strings that can be
// compared for less than, greater than, or equal to another string.

package intern

import (
	"sort"
)

// An LGE is a string that has been interned to an integer.  An LGE supports
// less than, greater than, and equal to comparisons (<, <=, >, >=, ==, !=)
// with other LGEs.
type LGE uint64

// lge maintains all the state needed to manipulate LGEs.
var lge state

// init initializes our global state.
func init() {
	lge.forgetAll()
}

// PreLGE provides advance notice of a string that will be interned using
// NewSymbolLGE.  A provided function canonicalizes the string.  Batching up a
// large number of PreLGE calls before calling NewSymbolLGE helps avoid
// running out of symbols that are properly comparable with all other symbols.
func PreLGE(s string) {
	lge.Lock()
	lge.pending = append(lge.pending, s)
	lge.Unlock()
}

// NewLGE maps a string to an LGE symbol.  It guarantees that two equal strings
// will always map to the same LGE.  However, it is possible that the package
// cannot accomodate a particular string, in which case NewLGE returns a
// non-nil error.  Pre-allocate as many LGEs as possible using PreLGE to reduce
// the likelihood of this happening.
func NewLGE(s string) (LGE, error) {
	var err error
	st := &lge
	st.Lock()
	defer st.Unlock()

	// Flush any pending symbols.
	err = st.flushPending()
	if err != nil {
		return 0, err
	}

	// Insert the new symbol.
	sym, err := st.assignSymbol(s, id, true)
	return LGE(sym), err
}

// String converts an LGE back to a string.  It panics if given an LGE that was
// not created using NewLGE.
func (s LGE) String() string {
	return lge.toString(uint64(s), "LGE")
}

// LGESlice is a slice of LGEs that implements sort.Interface.
type LGESlice []LGE

// Len returns the length of an LGESlice.
func (ls LGESlice) Len() int { return len(ls) }

// Less reports whether one element of an LGESlice is less than another.
func (ls LGESlice) Less(i, j int) bool { return ls[i] < ls[j] }

// Swap swaps two elements of an LGESLice.
func (ls LGESlice) Swap(i, j int) { ls[i], ls[j] = ls[j], ls[i] }

// Sort sorts an LGESlice in ascending order.
func (ls LGESlice) Sort() { sort.Sort(ls) }

// Search searches for x in a sorted LGESlice and returns the index as
// specified by sort.Search.  The return value is the index to insert x if x is
// not present.  (It can be len(ls).)  The slice must be sorted in ascending
// order.
func (ls LGESlice) Search(x LGE) int {
	return sort.Search(len(ls), func(i int) bool { return ls[i] >= x })
}
