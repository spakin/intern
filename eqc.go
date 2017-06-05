// This file provides the EqC type, which represents strings that can be
// compared only for equality.  Unlike Eq, EqC symbols are first canonicalized
// using a program-provided canonicalization function.

package intern

import (
	"fmt"
	"sync"
)

// An EqC is a string that has been interned to an integer.  EqC supports only
// equality comparisons, not inequality comparisons.  (No checks are performed
// to enforce that property, unfortunately.)  Unlike Eq, equality comparisons
// are insensitive to a transformation function provided to NewEqC.
//
// It is strongly recommended that programs alias EqC once for each
// transformation function.  This will help the compiler catch program errors
// if strings interned using different transformation functions are compared.
type EqC Eq

// eqcState maintains all the state needed to manipulate EqCs.
var eqcState struct {
	symToStr     map[EqC]string // Mapping from EqCs to strings
	strToSym     map[string]EqC // Mapping from strings to EqCs
	sync.RWMutex                // Mutex protecting both of the above
}

// init initializes our global state.
func init() {
	eqcState.symToStr = make(map[EqC]string)
	eqcState.strToSym = make(map[string]EqC)
}

// NewEqC maps a string to an EqC.  It guarantees that two strings that are
// equal after being passed through a function f will return the same EqC.
// That is, if f(s1) == f(s2) then NewEqC(s1) == NewEqC(s2).
func NewEqC(s string, f func(string) string) EqC {
	eqcState.Lock()
	defer eqcState.Unlock()
	fs := f(s)
	if sym, ok := eqcState.strToSym[fs]; ok {
		return sym
	}
	sym := EqC(len(eqcState.symToStr) + 1) // Reserve 0 to help catch program errors.
	eqcState.symToStr[sym] = s             // Use the original string when mapping a symbol to a string.
	eqcState.strToSym[fs] = sym            // Use the transformed string when mapping a string to a symbol.
	return sym
}

// String converts an EqC back to a string.  It panics if given an EqC that was
// not created using NewEqC.
func (s EqC) String() string {
	eqcState.RLock()
	defer eqcState.RUnlock()
	if str, ok := eqcState.symToStr[s]; ok {
		return str
	}
	panic(fmt.Sprintf("%d is not a valid intern.EqC", s))
}
