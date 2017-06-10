// This file provides the LGEC type, which represents strings that can be
// compared for less than, greater than, or equal to another string.  Unlike
// LGE, LGEC symbols are first canonicalized using a program-provided
// canonicalization function.

package intern

// An LGEC is a string that has been interned to an integer after being
// canonicalized using a program-provided transformation function.  An LGEC
// supports less than, greater than, and equal to comparisons (<, <=, >, >=,
// ==, !=) with other LGECs.
//
// It is strongly recommended that programs alias LGEC once for each
// transformation function.  This will help the compiler catch program errors
// if strings interned using different transformation functions are compared.
type LGEC uint64

// lgec maintains all the state needed to manipulate LGECs.
var lgec state

// init initializes our global state.
func init() {
	lgec.forgetAll()
}

// PreLGEC provides advance notice of a string that will be interned using
// NewSymbolLGEC.  A provided function canonicalizes the string.  Batching up a
// large number of PreLGEC calls before calling NewSymbolLGEC helps avoid
// running out of symbols that are properly comparable with all other symbols.
func PreLGEC(s string, f func(string) string) {
	lgec.Lock()
	lgec.pending = append(lgec.pending, f(s))
	lgec.Unlock()
}

// NewLGEC maps a string to an LGEC symbol.  It guarantees that two strings
// that compare equal after being passed through a function f will return the
// same LGEC.  However, it is possible that the package cannot accomodate a
// particular string, in which case NewLGEC returns a non-nil error.
// Pre-allocate as many LGECs as possible using PreLGEC to reduce the
// likelihood of this happening.
func NewLGEC(s string, f func(string) string) (LGEC, error) {
	var err error
	st := &lgec
	st.Lock()
	defer st.Unlock()

	// Flush any pending symbols.
	err = st.flushPending()
	if err != nil {
		return 0, err
	}

	// Insert the new symbol.
	sym, err := st.assignSymbol(s, f, true)
	return LGEC(sym), err
}

// String converts an LGEC back to a string.  It panics if given an LGEC that was
// not created using NewLGEC.
func (s LGEC) String() string {
	return lgec.toString(uint64(s), "LGEC")
}
