// This file provides the EqC type, which represents strings that can be
// compared only for equality.  Unlike Eq, EqC symbols are first canonicalized
// using a program-provided canonicalization function.

package intern

import (
	"fmt"
)

// An EqC is a string that has been interned to an integer after being
// canonicalized using a program-provided transformation function.  EqC
// supports only equality comparisons, not inequality comparisons.  (No checks
// are performed to enforce that property, unfortunately.)
//
// It is strongly recommended that programs alias EqC once for each
// transformation function.  This will help the compiler catch program errors
// if strings interned using different transformation functions are compared.
type EqC uint64

// eqc maintains all the state needed to manipulate EqCs.
var eqc state

// init initializes our global state.
func init() {
	eqc.forgetAll()
}

// NewEqC maps a string to an EqC symbol.  It guarantees that two strings that
// are equal after being passed through a function f will return the same EqC.
func NewEqC(s string, f func(string) string) EqC {
	var err error
	st := &eqc
	st.Lock()
	defer st.Unlock()
	sym, err := st.assignSymbol(s, f, false)
	if err != nil {
		panic(fmt.Sprintf("Internal error: Unexpected error (%s)", err))
	}
	return EqC(sym)
}

// String converts an EqC back to a string.  It panics if given an EqC that was
// not created using NewEqC.
func (s EqC) String() string {
	return eqc.toString(uint64(s), "EqC")
}
