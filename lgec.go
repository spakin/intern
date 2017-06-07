// This file provides the LGEC type, which represents strings that can be
// compared for less than, greater than, or equal to another string.  Unlike
// LGE, LGEC symbols are first canonicalized using a program-provided
// canonicalization function.

package intern

import (
	"fmt"
	"sync"
)

// An LGEC is a string that has been interned to an integer after being
// canonicalized using a program-provided transformation function.  An LGEC
// supports less than, greater than, and equal to comparisons (<, <=, >, >=,
// ==, !=) with other LGECs.
type LGEC uint64

// lgecState maintains all the state needed to manipulate LGECs.
var lgecState struct {
	symToStr     map[LGEC]string // Mapping from LGECs to strings
	strToSym     map[string]LGEC // Mapping from strings to LGECs
	tree         *tree           // Tree for maintaining LGEC assignments
	pending      []string        // Strings not yet mapped to LGECs
	sync.RWMutex                 // Mutex protecting all of the above
}

// forgetAllLGEC discards all extant string/symbol mappings and resets the
// assignment tables to their initial state.
func forgetAllLGEC() {
	lgecState.Lock()
	lgecState.symToStr = make(map[LGEC]string)
	lgecState.strToSym = make(map[string]LGEC)
	lgecState.tree = nil
	lgecState.pending = make([]string, 0, 100)
	lgecState.Unlock()
}

// init initializes our global state.
func init() {
	forgetAllLGEC()
}

// PreLGEC provides advance notice of a string that will be interned using
// NewSymbolLGEC.  A provided function canonicalizes the string.  Batching up a
// large number of PreLGEC calls before calling NewSymbolLGEC helps avoid
// running out of symbols that are properly comparable with all other symbols.
func PreLGEC(s string, f func(string) string) {
	lgecState.Lock()
	lgecState.pending = append(lgecState.pending, f(s))
	lgecState.Unlock()
}

// NewLGEC maps a string to an LGEC symbol.  It guarantees that two strings
// that are equal after being passed through a function f will return the same
// LGEC.  However, it is possible that the package cannot accomodate a
// particular string, in which case NewLGEC returns a non-nil error.
// Pre-allocate as many LGECs as possible using PreLGEC to reduce the
// likelihood of this happening.
func NewLGEC(s string, f func(string) string) (LGEC, error) {
	var err error
	st := &lgecState
	st.Lock()
	defer st.Unlock()

	// Flush all pending symbols.
	if len(st.pending) > 0 {
		st.tree, err = st.tree.insertMany(st.pending)
		if err != nil {
			return 0, err
		}
		st.pending = st.pending[:0]
		st.tree.walk(func(t *tree) {
			st.symToStr[LGEC(t.sym)] = t.str
			st.strToSym[t.str] = LGEC(t.sym)
		})
	}

	// Insert the new symbol.
	fs := f(s)
	st.tree, err = st.tree.insert(fs)
	if err != nil {
		return 0, err
	}
	t := st.tree.find(fs)
	if t == nil {
		panic("Internal error: Failed to find a string just inserted into a tree")
	}
	sym := LGEC(t.sym)
	st.symToStr[sym] = s  // Use the original string when mapping a symbol to a string.
	st.strToSym[fs] = sym // Use the transformed string when mapping a string to a symbol.
	return sym, nil
}

// String converts an LGEC back to a string.  It panics if given an LGEC that was
// not created using NewLGEC.
func (s LGEC) String() string {
	lgecState.RLock()
	defer lgecState.RUnlock()
	if str, ok := lgecState.symToStr[s]; ok {
		return str
	}
	panic(fmt.Sprintf("%d is not a valid intern.LGEC", s))
}
