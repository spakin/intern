// This file provides the Eq type, which represents strings that can
// be compared only for equality.

package intern

import (
	"fmt"
	"sync"
)

// An Eq is a string that has been interned to an integer.  Eq
// supports only equality comparisons, not inequality comparisons.
// (No checks are performed to enforce that property, unfortunately.)
type Eq uint64

// eqState maintains all the state needed to manipulate Eqs.
var eqState struct {
	symToStr     map[Eq]string // Mapping from Eqs to strings
	strToSym     map[string]Eq // Mapping from strings to Eqs
	sync.RWMutex               // Mutex protecting both of the above
}

// init initializes our global state.
func init() {
	eqState.symToStr = make(map[Eq]string)
	eqState.strToSym = make(map[string]Eq)
}

// NewEq maps a string to an Eq symbol.  It guarantees that the same string
// contents will always return the same symbol.
func NewEq(s string) Eq {
	eqState.Lock()
	defer eqState.Unlock()
	if sym, ok := eqState.strToSym[s]; ok {
		return sym
	}
	sym := Eq(len(eqState.symToStr) + 1) // Reserve 0 to help catch program errors.
	eqState.symToStr[sym] = s
	eqState.strToSym[s] = sym
	return sym
}

// String converts an Eq back to a string.  It panics if given an Eq that was
// not created using NewEq.
func (s Eq) String() string {
	eqState.RLock()
	defer eqState.RUnlock()
	if str, ok := eqState.symToStr[s]; ok {
		return str
	}
	panic(fmt.Sprintf("%d is not a valid intern.Eq", s))
}
