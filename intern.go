/*
Package intern maps symbols to integers for fast comparisons.
*/
package intern

import (
	"fmt"
)

// An EqSymbol is represented by an integer.  It supports only equality
// comparisons, not inequality comparisons.  (No checks are performed to
// enforce that property, however.)
type EqSymbol uint64

// eqSymbolTostring maintains the mapping from EqSymbols to strings.
var eqSymbolToString = make(map[EqSymbol]string)

// stringToEqSymbol maintains the mapping from strings to EqSymbols.
var stringToEqSymbol = make(map[string]EqSymbol)

// NewEqSymbol maps a string to an EqSymbol.  It guarantees that the same
// string contents will always return the same EqSymbol.
func NewEqSymbol(s string) EqSymbol {
	if sym, ok := stringToEqSymbol[s]; ok {
		return sym
	}
	sym := EqSymbol(len(eqSymbolToString) + 1) // Reserve 0 to help catch program errors.
	eqSymbolToString[sym] = s
	stringToEqSymbol[s] = sym
	return sym
}

// String converts an EqSymbol back to a string.  It panics if given an invalid
// input.
func (s EqSymbol) String() string {
	if str, ok := eqSymbolToString[s]; ok {
		return str
	} else {
		panic(fmt.Sprintf("%d is not a valid EqSymbol", s))
	}
}
