// This file provides the Eq type, which represents strings that can
// be compared only for equality.

package intern

import (
	"fmt"
)

// An Eq is a string that has been interned to an integer.  Eq
// supports only equality comparisons, not inequality comparisons.
// (No checks are performed to enforce that property, unfortunately.)
type Eq uint64

// eq maintains all the state needed to manipulate Eqs.
var eq state

// init initializes our global state.
func init() {
	eq.forgetAll()
}

// NewEq maps a string to an Eq symbol.  It guarantees that two equal strings
// will always map to the same Eq.
func NewEq(s string) Eq {
	var err error
	st := &eq
	st.Lock()
	defer st.Unlock()
	sym, err := st.assignSymbol(s, false)
	if err != nil {
		panic(fmt.Sprintf("Internal error: Unexpected error (%s)", err))
	}
	return Eq(sym)
}

// String converts an Eq back to a string.  It panics if given an Eq that was
// not created using NewEq.
func (s Eq) String() string {
	return eq.toString(uint64(s), "Eq")
}
