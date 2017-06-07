// This file provides the LGE type, which represents strings that can be
// compared for less than, greater than, or equal to another string.

package intern

import (
	"fmt"
	"sync"
)

// An LGE is a string that has been interned to an integer.  An LGE supports
// less than, greater than, and equal to comparisons (<, <=, >, >=, ==, !=)
// with other LGEs.
type LGE uint64

// lgeState maintains all the state needed to manipulate LGEs.
var lgeState struct {
	symToStr     map[LGE]string // Mapping from LGEs to strings
	strToSym     map[string]LGE // Mapping from strings to LGEs
	tree         *tree          // Tree for maintaining LGE assignments
	pending      []string       // Strings not yet mapped to LGEs
	sync.RWMutex                // Mutex protecting all of the above
}

// forgetAllLGE discards all extant string/symbol mappings and resets the
// assignment tables to their initial state.
func forgetAllLGE() {
	lgeState.Lock()
	lgeState.symToStr = make(map[LGE]string)
	lgeState.strToSym = make(map[string]LGE)
	lgeState.tree = nil
	lgeState.pending = make([]string, 0, 100)
	lgeState.Unlock()
}

// init initializes our global state.
func init() {
	forgetAllLGE()
}

// PreLGE provides advance notice of a string that will be interned using
// NewSymbolLGE.  Batching up a large number of PreLGE calls before calling
// NewSymbolLGE helps avoid running out of symbols that are properly comparable
// with all other symbols.
func PreLGE(s string) {
	lgeState.Lock()
	lgeState.pending = append(lgeState.pending, s)
	lgeState.Unlock()
}

// NewLGE maps a string to an LGE symbol.  It guarantees that the same string
// contents will always return the same symbol.  However, it is possible that
// the package cannot accomodate a particular string, in which case NewLGE
// returns a non-nil error.  Pre-allocate as many LGEs as possible using PreLGE
// to reduce the likelihood of this happening.
func NewLGE(s string) (LGE, error) {
	var err error
	st := &lgeState
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
			st.symToStr[LGE(t.sym)] = t.str
			st.strToSym[t.str] = LGE(t.sym)
		})
	}

	// Insert the new symbol.
	st.tree, err = st.tree.insert(s)
	if err != nil {
		return 0, err
	}
	t := st.tree.find(s)
	if t == nil {
		panic("Internal error: Failed to find a string just inserted into a tree")
	}
	sym := LGE(t.sym)
	st.symToStr[sym] = s
	st.strToSym[s] = sym
	return sym, nil
}

// String converts an LGE back to a string.  It panics if given an LGE that was
// not created using NewLGE.
func (s LGE) String() string {
	lgeState.RLock()
	defer lgeState.RUnlock()
	if str, ok := lgeState.symToStr[s]; ok {
		return str
	}
	panic(fmt.Sprintf("%d is not a valid intern.LGE", s))
}
