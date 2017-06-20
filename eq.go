// This file provides the Eq type, which represents strings that can
// be compared only for equality.

package intern

// An Eq is a string that has been interned to an integer.  Eq supports only
// equality and inequality comparisons, not greater than/less than comparisons.
// (No checks are performed to enforce that usage model, unfortunately.)
type Eq symbol

// eq maintains all the state needed to manipulate Eqs.
var eq state

// init initializes our global state.
func init() {
	eq.forgetAll()
}

// assignEq assigns the next available Eq symbol to a string and returns the
// new symbol.  If the string already has an Eq associated with it, return the
// old Eq without allocating a new one.
func assignEq(s string) Eq {
	// Check if the string was already assigned a symbol.
	sym, ok := eq.strToSym[s]
	if ok {
		return Eq(sym)
	}

	// We haven't seen this string before.  Find a symbol for it.
	sym = symbol(len(eq.symToStr) + 1)
	eq.symToStr[sym] = s
	eq.strToSym[s] = sym
	return Eq(sym)
}

// NewEq maps a string to an Eq symbol.  It guarantees that two equal strings
// will always map to the same Eq.
func NewEq(s string) Eq {
	eq.Lock()
	defer eq.Unlock()
	return assignEq(s)
}

// NewEqMulti performs the same operation as NewEq but accepts a slice of
// strings instead of an individual string.  This amortizes some costs when
// allocating a large number of Eqs at once.
func NewEqMulti(ss []string) []Eq {
	eq.Lock()
	defer eq.Unlock()
	syms := make([]Eq, len(ss))
	for i, s := range ss {
		syms[i] = assignEq(s)
	}
	return syms
}

// String converts an Eq back to a string.  It panics if given an Eq that was
// not created using NewEq.
func (s Eq) String() string {
	return eq.toString(symbol(s), "Eq")
}

// ForgetAllEqs discards all existing mappings from strings to Eqs so the
// associated memory can be reclaimed.  Use this function only when you know
// for sure that no previously mapped Eqs will subsequently be used.
func ForgetAllEqs() {
	eq.Lock()
	eq.forgetAll()
	eq.Unlock()
}
