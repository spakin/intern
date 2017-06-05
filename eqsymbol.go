/*
Package intern maps symbols to integers for fast comparisons.
*/
package intern

import (
	"fmt"
	"sync"
)

// An Eq is represented by an integer.  It supports only equality
// comparisons, not inequality comparisons.  (No checks are performed to
// enforce that property, however.)
type Eq uint64

// eqSymbolState maintains all the state needed to manipulate Eqs.
var eqSymbolState struct {
	symToStr     map[Eq]string // Mapping from Eqs to strings
	strToSym     map[string]Eq // Mapping from strings to Eqs
	sync.RWMutex               // Mutex protecting both of the above
}

// init initializes our global state.
func init() {
	eqSymbolState.symToStr = make(map[Eq]string)
	eqSymbolState.strToSym = make(map[string]Eq)
}

// NewEq maps a string to an Eq.  It guarantees that the same
// string contents will always return the same Eq.
func NewEq(s string) Eq {
	eqSymbolState.Lock()
	defer eqSymbolState.Unlock()
	if sym, ok := eqSymbolState.strToSym[s]; ok {
		return sym
	}
	sym := Eq(len(eqSymbolState.symToStr) + 1) // Reserve 0 to help catch program errors.
	eqSymbolState.symToStr[sym] = s
	eqSymbolState.strToSym[s] = sym
	return sym
}

// String converts an Eq back to a string.  It panics if given an invalid
// input.
func (s Eq) String() string {
	eqSymbolState.RLock()
	defer eqSymbolState.RUnlock()
	if str, ok := eqSymbolState.symToStr[s]; ok {
		return str
	}
	panic(fmt.Sprintf("%d is not a valid intern.Eq", s))
}
