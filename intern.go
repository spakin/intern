/*
Package intern maps symbols to integers for fast comparisons.
*/
package intern

import (
	"fmt"
	"sync"
)

// state includes all the state needed to manipulate all interned-string types.
type state struct {
	symToStr     map[uint64]string // Mapping from symbols to strings
	strToSym     map[string]uint64 // Mapping from strings to symbols
	tree         *tree             // Tree for maintaining symbols assignments
	pending      []string          // Strings not yet mapped to symbols
	sync.RWMutex                   // Mutex protecting all of the above
}

// forgetAll discards all extant string/symbol mappings and resets the
// assignment tables to their initial state.
func (st *state) forgetAll() {
	st.Lock()
	st.symToStr = make(map[uint64]string)
	st.strToSym = make(map[string]uint64)
	st.tree = nil
	st.pending = make([]string, 0, 100)
	st.Unlock()
}

// id is the identity function (string → string).  It's used when no string
// canonicalization is required.
func id(s string) string { return s }

// toString converts a symbol back to a string.  It panics if given a symbol
// that was not created using New*.
func (st *state) toString(s uint64, ty string) string {
	st.RLock()
	defer st.RUnlock()
	if str, ok := st.symToStr[s]; ok {
		return str
	}
	panic(fmt.Sprintf("%d is not a valid intern.%s", s, ty))
}

// flushPending flushes all pending symbols.  It returns an error value.
func (st *state) flushPending() error {
	var err error
	if len(st.pending) > 0 {
		st.tree, err = st.tree.insertMany(st.pending)
		if err != nil {
			return err
		}
		st.pending = st.pending[:0]
		st.tree.walk(func(t *tree) {
			st.symToStr[t.sym] = t.str
			st.strToSym[t.str] = t.sym
		})
	}
	return nil
}

// assignSymbol assigns a new symbol to a string.  It takes as arguments the
// string to assign, a canonicalization function, and a Boolean that indicates
// whether to use a tree to preserve order.  This method returns the assigned
// symbol and an error value.
func (st *state) assignSymbol(s string, f func(string) string, useTree bool) (uint64, error) {
	fs := f(s)
	var sym uint64 // Symbol to assign to string s
	if useTree {
		// Assign the symbol using a tree to maintain order.
		var err error
		st.tree, err = st.tree.insert(fs)
		if err != nil {
			return 0, err
		}
		t := st.tree.find(fs)
		if t == nil {
			panic("Internal error: Failed to find a string just inserted into a tree")
		}
		sym = t.sym
	} else {
		// Assign the next available number, starting at 1 to ensure
		// that an uninitialized symbol is treated as invalid.
		var ok bool
		sym, ok = st.strToSym[fs]
		if ok {
			// The string was already assigned a symbol.
			return sym, nil
		}
		sym = uint64(len(st.symToStr)) + 1
	}
	st.symToStr[sym] = s  // Use the original string when mapping a symbol to a string.
	st.strToSym[fs] = sym // Use the transformed string when mapping a string to a symbol.
	return sym, nil
}
