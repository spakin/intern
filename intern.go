/*
Package intern maps strings to integers for fast comparisons.

Description

Consider the string "This is string one" stored in variable x and the string
"This is string two" stored in variable y.  When a program tests if x == y, it
compares their first characters and finds that both are "T".  It then compares
their second characters and finds that both are "h".  This continues until the
16th character, where it finds that "o" is different from "t" and concludes
that the two strings are different.

The intern package provides two "symbol" data types to speed up such
comparisons.  The first, Eq, is an integer that is associated with particular
string contents.  For example, "This is string one" may be mapped to 52413
while "This is string two" may be mapped to 29488.  Comparing 52413 == 29488
fails in a single step, which can save a lot of time if such comparisons are
performed frequently.

The second data type, LGE, supports less than, greater than, and equal to
comparisons (and variations such as <=, >=, and !=).  For example, if x is "But
I have promises to keep", y is "And miles to go before I sleep", and z is also
"And miles to go before I sleep", then these might map to symbols 48712, 29519,
and 29519, respectively.  Let's call these sx, sy, and sz.  Then, x > y and sx
> sy; y == z and sy == sz; z < x and sz < sx; and any other such comparison of
x, y, and z, also holds for sx, sy, and sz.

A String method is defined for both Eq and LGE.  This maps a symbol back to its
original string.  Ergo, no information is lost when mapping strings to symbols.

Usage

NewEq maps a string to an Eq symbol, and NewLGE maps a string to an LGE symbol.
The former is faster and always succeeds (until the set of 64-bit integers is
exhausted).  NewLGE, in addition to being slower, can fail if earlier
assignments of integers to strings preclude a new string from being mapped to
an integer that respects comparisons with all existing symbols.  (In the
current implementation, a worst-case sequence of NewLGE calls will fail on the
65th call.)  As a workaround, the package provides a PreLGE function that
indicates an intention to invoke NewLGE on a particular string but without
actually assigning an integer.

Best practice is to pre-allocate as many LGE symbols as possible before calling
NewLGE.  When NewLGE is called, all strings previously passed to PreLGE are
assigned integers in an order that helps ensure that the desired relations are
preserved.  The process of pre-allocating LGE symbols with PreLGE and later
allocating them with NewLGE can be repeated as many times as necessary but with
increasing likelihood of failure with each repetition.  If NewLGE does fail,
the RemapAllLGEs function can be called to completely redo the mapping from
strings to LGE symbols.  The program will need to update any live LGE symbols
it has stored in data structures.

All functions in this package are thread-safe.

Performance

It's tricky to discuss the speed of Eq and LGE symbol comparisons relative to
string comparisons.  First, the time needed to compare two strings is a
function of the length of the common prefix.  For example, comparing "Jonas
Grumby" to "Roy Hinkley" is faster than comparing
"antidisestablishmentarianism" to "antidisestablishmentarianist".  In contrast,
symbol comparisons take the same amount of time regardless of the strings they
represent.  Second, symbols must be allocated, which takes time (more for an
LGE than an Eq).  Consequently, if few comparisons are performed per
allocation, it may be faster to use strings than symbols.  As a further twist,
the NewEqMulti and PreLGEMulti/NewLGEMulti functions help amortize some of the
allocation costs when allocating multiple symbols at once, but not all programs
have multiple strings they need to intern at once.  Third, there is a memory
cost associated with maintaining a bidirectional mapping between strings and
symbols.  If this extra memory causes a program's working set to expand beyond
the size of a cache, it may be faster to use strings than symbols.

The intern package includes a number of benchmarks to help programmers
determine if it may be beneficial to use the package.  Run them in
the usual manner:

    go test --bench=. --run=None

On my computer, the results show, very roughly, the following:

1) For random strings (i.e., very short common prefixes), symbols are faster
than strings if the program performs more than twice as many comparisons as Eq
allocations or more than 13 times as many comparisons as LGE allocations.

2) For very long, very similar strings (well over 100 initial characters in
common), both Eq and LGE symbols are faster than strings if the program
performs at least as many comparisons as symbol allocations.

*/
package intern

import (
	"fmt"
	"sync"
)

// These constants represent the various error codes the package can return.
const (
	ErrTableFull   = iota + 1 // Symbol table is full
	ErrRemapFailed            // Symbol remapping failed
)

// PkgError represents an error specific to the intern package, as opposed to
// an error generated by an underlying package.
type PkgError struct {
	Code int    // Type of error that occurred
	Str  string // String that triggered the error (if applicable)
	msg  string // Textual description of the error
}

// Error reports a PkgError as a string.
func (e *PkgError) Error() string {
	return e.msg
}

// symbol represents either package symbol type (Eq or LGE).
type symbol uint64

// symbolList is a list of symbols.  It implements sort.Interface.
type symbolList []symbol

// Len returns the length of a symbolList.
func (sl symbolList) Len() int { return len(sl) }

// Less says is one symbolList element is less than another.
func (sl symbolList) Less(i, j int) bool { return sl[i] < sl[j] }

// Swap swaps two elements of a symbolList.
func (sl symbolList) Swap(i, j int) { sl[i], sl[j] = sl[j], sl[i] }

// state includes all the state needed to manipulate all interned-string types.
type state struct {
	symToStr     map[symbol]string // Mapping from symbols to strings
	strToSym     map[string]symbol // Mapping from strings to symbols
	tree         *tree             // Tree for maintaining symbols assignments
	pending      []string          // Strings not yet mapped to symbols
	sync.RWMutex                   // Mutex protecting all of the above
}

// forgetAll discards all extant string/symbol mappings and resets the
// assignment tables to their initial state.
func (st *state) forgetAll() {
	st.symToStr = make(map[symbol]string)
	st.strToSym = make(map[string]symbol)
	st.tree = nil
	st.pending = make([]string, 0, 100)
}

// toString converts a symbol back to a string.  It panics if given a symbol
// that was not created using New*.
func (st *state) toString(s symbol, ty string) string {
	st.RLock()
	defer st.RUnlock()
	if str, ok := st.symToStr[s]; ok {
		return str
	}
	panic(fmt.Sprintf("%d is not a valid intern.%s", s, ty))
}

// flushPending flushes all pending symbols, converting strings to symbols.
// The function returns an error status.
func (st *state) flushPending() error {
	var err error
	if len(st.pending) > 0 {
		var sMap map[string]symbol
		st.tree, sMap, err = st.tree.insertMany(st.pending)
		if err != nil {
			return err
		}
		st.pending = st.pending[:0]
		for k, v := range sMap {
			st.strToSym[k] = v
			st.symToStr[v] = k
		}
	}
	return nil
}

// getSymbol looks up and returns the symbol associated with a string.  It
// aborts the program on failure.
func (st *state) getSymbol(s string) symbol {
	sym, ok := st.strToSym[s]
	if !ok {
		panic(fmt.Sprintf("Internal error: Expected to find an interned version of %q", s))
	}
	return sym
}
