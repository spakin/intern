// This file provides a simple binary-tree abstraction for use by LGE symbols.

package intern

import (
	"fmt"
	"sort"
)

// A tree represents a binary tree of strings.
type tree struct {
	str   string // Contents of this node
	sym   symbol // Symbol to assign to this string (LTE or LTEC)
	left  *tree  // Left child or nil
	right *tree  // Right child or nil
}

// insert inserts a string into a tree, returning the new tree, the inserted
// symbol, and an error value.
func (t *tree) insert(s string) (*tree, symbol, error) {
	return t.insertHelper(s, 1<<63, 1<<62)
}

// insertHelper inserts a string into a tree, returning the new tree, the
// inserted symbol, and an error value.  It performs almost all of the work for
// the top-level insert method.
func (t *tree) insertHelper(s string, val, incr symbol) (*tree, symbol, error) {
	if t == nil {
		return &tree{str: s, sym: val}, val, nil
	}
	if incr == 0 {
		e := &PkgError{
			Code: ErrTableFull,
			Str:  s,
			msg:  fmt.Sprintf("Unable to insert %q; symbol table is full", s),
		}
		return nil, 0, e
	}
	var sym symbol
	var err error
	switch {
	case s == t.str:
		sym = val
	case s < t.str:
		t.left, sym, err = t.left.insertHelper(s, val-incr, incr/2)
	case s > t.str:
		t.right, sym, err = t.right.insertHelper(s, val+incr, incr/2)
	}
	return t, sym, err
}

// insertMany inserts a list of strings into a tree, attempting to maintain
// balance as it does so.  A new tree, a map from strings to symbols, and an
// error value are returned.  It is assumed that the given list of strings is
// non-empty.
func (t *tree) insertMany(ss []string) (*tree, map[string]symbol, error) {
	// Create a sorted version of the list of strings.
	sss := make([]string, len(ss))
	for i, s := range ss {
		sss[i] = s
	}
	sort.Strings(sss)

	// Call our helper function then construct a map based on the list of
	// symbols it returns.
	tNew, syms, err := t.insertManySorted(sss)
	if err != nil {
		return nil, nil, err
	}
	sort.Sort(syms)
	m := make(map[string]symbol, len(sss))
	for i, s := range sss {
		m[s] = syms[i]
	}
	return tNew, m, nil
}

// insertManySorted inserts a sorted list of strings into a tree, attempting to
// maintain balance as it does so.  It performs most of the work for
// insertMany.  It is assumed that the given list of strings is non-empty.
func (t *tree) insertManySorted(ss []string) (*tree, symbolList, error) {
	// Handle the base case (a single string) first.
	n := len(ss)
	if n == 1 {
		tNew, s, err := t.insert(ss[0])
		return tNew, symbolList{s}, err
	}

	// Insert the middle element, then recursively insert the left and
	// right sub-slices.
	mid := n / 2
	tNew, sym, err := t.insert(ss[mid])
	if err != nil {
		return nil, nil, err
	}
	var lSyms, rSyms symbolList
	if mid > 0 {
		tNew, lSyms, err = tNew.insertManySorted(ss[:mid])
		if err != nil {
			return nil, nil, err
		}
	}
	if mid+1 < n {
		tNew, rSyms, err = tNew.insertManySorted(ss[mid+1:])
		if err != nil {
			return nil, nil, err
		}
	}
	sList := append(lSyms, sym)
	sList = append(sList, rSyms...)
	return tNew, sList, nil
}
