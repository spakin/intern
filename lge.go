// This file provides the LGE type, which represents strings that can be
// compared for less than, greater than, or equal to another string.

package intern

import "fmt"

// An LGE is a string that has been interned to an integer.  An LGE supports
// less than, greater than, and equal to comparisons (<, <=, >, >=, ==, !=)
// with other LGEs.
type LGE symbol

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

// PreLGEMulti performs the same operation as PreLGE but accepts a slice of
// strings instead of an individual string.  This amortizes some costs when
// pre-allocating a large number of LGEs at once.
func PreLGEMulti(ss []string) {
	lge.Lock()
	lge.pending = append(lge.pending, ss...)
	lge.Unlock()
}

// NewLGE maps a string to an LGE symbol.  It guarantees that two equal strings
// will always map to the same LGE.  However, it is possible that the package
// cannot accommodate a particular string, in which case NewLGE returns a
// non-nil error.  Pre-allocate as many LGEs as possible using PreLGE to reduce
// the likelihood of that happening.
func NewLGE(s string) (LGE, error) {
	// Acquire a lock on LGE state.
	var err error
	lge.Lock()
	defer lge.Unlock()

	// Mark the new string as pending then flush all pending symbols.
	lge.pending = append(lge.pending, s)
	err = lge.flushPending()
	if err != nil {
		return 0, err
	}

	// Return the new symbol
	return LGE(lge.getSymbol(s)), nil
}

// NewLGEMulti performs the same operation as NewLGE but accepts a slice of
// strings instead of an individual string.  This amortizes some costs when
// allocating a large number of LGEs at once.
func NewLGEMulti(ss []string) ([]LGE, error) {
	// Acquire a lock on LGE state.
	var err error
	lge.Lock()
	defer lge.Unlock()

	// Mark all new strings as pending then flush all pending symbols.
	syms := make([]LGE, len(ss))
	if len(ss) == 0 {
		return syms, nil
	}
	lge.pending = append(lge.pending, ss...)
	err = lge.flushPending()
	if err != nil {
		return syms, err
	}

	// Return the new symbols.
	for i, s := range ss {
		syms[i] = LGE(lge.getSymbol(s))
	}
	return syms, nil
}

// String converts an LGE back to a string.  It panics if given an LGE that was
// not created using NewLGE.
func (s LGE) String() string {
	return lge.toString(symbol(s), "LGE")
}

// ForgetAllLGEs discards all existing mappings from strings to LGEs so the
// associated memory can be reclaimed.  Use this function only when you know
// for sure that no previously mapped LGEs will subsequently be used.
func ForgetAllLGEs() {
	lge.Lock()
	lge.forgetAll()
	lge.Unlock()
}

// RemapAllLGEs reassigns LGEs to strings to help clean up the mapping.  This
// provides a way to add strings that were previously rejected by NewLGE.
// RemapAllLGEs returns a mapping from old LGEs to new LGEs to assist programs
// with updating LGEs that are in use.
func RemapAllLGEs() (map[LGE]LGE, error) {
	// Store the existing LGE state then reinitialize it.
	lge.Lock()
	defer lge.Unlock()
	oldLge := state{
		pending:  lge.pending,
		strToSym: lge.strToSym,
	}
	lge.forgetAll()

	// Append the old list of strings to the pending list.
	lge.pending = oldLge.pending
	for s := range oldLge.strToSym {
		lge.pending = append(lge.pending, s)
	}

	// Map all pending strings to LGEs.
	err := lge.flushPending()
	if err != nil {
		return nil, err
	}

	// Construct a map from old to new LGEs and return it.
	m := make(map[LGE]LGE, len(lge.strToSym))
	for str, oldSym := range oldLge.strToSym {
		newSym, ok := lge.strToSym[str]
		if !ok {
			e := &PkgError{
				Code: ErrRemapFailed,
				Str:  str,
				msg:  fmt.Sprintf("Failed to remap string %q", str),
			}
			return nil, e
		}
		m[LGE(oldSym)] = LGE(newSym)
	}
	return m, nil
}

// MarshalText converts an LGE to a string and that string to a slice of bytes.
// With this method, LGE implements the encoding.TextMarshaler interface.
func (s *LGE) MarshalText() ([]byte, error) {
	return []byte(lge.toString(symbol(*s), "LGE")), nil
}

// UnmarshalText converts an slice of bytes to a string then interns that
// string to an LGE.  With this method, LGE implements the
// encoding.TextUnmarshaler interface.
func (s *LGE) UnmarshalText(text []byte) error {
	var err error
	*s, err = NewLGE(string(text))
	return err
}

// MarshalBinary converts an LGE to a string and that string to a slice of
// bytes.  With this method, LGE implements the encoding.BinaryMarshaler
// interface.
func (s *LGE) MarshalBinary() ([]byte, error) {
	return []byte(lge.toString(symbol(*s), "LGE")), nil
}

// UnmarshalBinary converts an slice of bytes to a string then interns that
// string to an LGE.  With this method, LGE implements the
// encoding.BinaryUnmarshaler interface.
func (s *LGE) UnmarshalBinary(data []byte) error {
	var err error
	*s, err = NewLGE(string(data))
	return err
}
