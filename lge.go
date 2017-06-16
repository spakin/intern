// This file provides the LGE type, which represents strings that can be
// compared for less than, greater than, or equal to another string.

package intern

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
// NewLGE.  Batching up a large number of PreLGE calls before calling NewLGE
// helps avoid running out of symbols that are properly comparable with all
// other symbols.
func PreLGE(s string) {
	lge.Lock()
	lge.pending = append(lge.pending, s)
	lge.Unlock()
}

// NewLGE maps a string to an LGE symbol.  It guarantees that two equal strings
// will always map to the same LGE.  However, it is possible that the package
// cannot accommodate a particular string, in which case NewLGE returns a
// non-nil error.  Pre-allocate as many LGEs as possible using PreLGE to reduce
// the likelihood of that happening.
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
	sym, err := st.assignSymbol(s, true)
	return LGE(sym), err
}

// String converts an LGE back to a string.  It panics if given an LGE that was
// not created using NewLGE.
func (s LGE) String() string {
	return lge.toString(uint64(s), "LGE")
}

// ForgetAllLGEs discards all existing mappings from strings to LGEs so the
// associated memory can be reclaimed.  Use only when you know for sure that no
// previously mapped LGEs will subsequently be used.
func ForgetAllLGEs() {
	lge.forgetAll()
}
