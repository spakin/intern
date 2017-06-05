/*
Package intern maps symbols to integers for fast comparisons.
*/
package intern

import (
	"fmt"
	"sync"
)

// An EqSymbol is represented by an integer.  It supports only equality
// comparisons, not inequality comparisons.  (No checks are performed to
// enforce that property, however.)
type EqSymbol uint64

// eqSymbolState maintains all the state needed to manipulate EqSymbols.
var eqSymbolState struct {
	symToStr     map[EqSymbol]string // Mapping from EqSymbols to strings
	strToSym     map[string]EqSymbol // Mapping from strings to EqSymbols
	sync.RWMutex                     // Mutex protecting both of the above
}

// init initializes our global state.
func init() {
	eqSymbolState.symToStr = make(map[EqSymbol]string)
	eqSymbolState.strToSym = make(map[string]EqSymbol)
}

// NewEqSymbol maps a string to an EqSymbol.  It guarantees that the same
// string contents will always return the same EqSymbol.
func NewEqSymbol(s string) EqSymbol {
	eqSymbolState.Lock()
	defer eqSymbolState.Unlock()
	if sym, ok := eqSymbolState.strToSym[s]; ok {
		return sym
	}
	sym := EqSymbol(len(eqSymbolState.symToStr) + 1) // Reserve 0 to help catch program errors.
	eqSymbolState.symToStr[sym] = s
	eqSymbolState.strToSym[s] = sym
	return sym
}

// String converts an EqSymbol back to a string.  It panics if given an invalid
// input.
func (s EqSymbol) String() string {
	eqSymbolState.RLock()
	defer eqSymbolState.RUnlock()
	if str, ok := eqSymbolState.symToStr[s]; ok {
		return str
	}
	panic(fmt.Sprintf("%d is not a valid EqSymbol", s))
}
