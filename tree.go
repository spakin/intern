// This file provides a simple binary-tree abstraction for use by LGE symbols.

package intern

import (
	"fmt"
	"sort"
)

// A tree represents a binary tree of strings.
type tree struct {
	str   string // Contents of this node
	sym   uint64 // Symbol to assign to this string (LTE or LTEC)
	left  *tree  // Left child or nil
	right *tree  // Right child or nil
}

// insert inserts a string into a tree, returning the new tree.
func (t *tree) insert(s string) (*tree, error) {
	return t.insertHelper(s, 1<<63, 1<<62)
}

// insertHelper inserts a string into a tree, returning the new tree.  It
// performs almost all of the work for the top-level insert method.
func (t *tree) insertHelper(s string, val, incr uint64) (*tree, error) {
	if t == nil {
		return &tree{str: s, sym: val}, nil
	}
	if incr == 0 {
		return nil, fmt.Errorf("Unable to insert %q; sub-tree is full", s)
	}
	var err error
	switch {
	case s == t.str:
	case s < t.str:
		t.left, err = t.left.insertHelper(s, val-incr, incr/2)
	case s > t.str:
		t.right, err = t.right.insertHelper(s, val+incr, incr/2)
	}
	return t, err
}

// insertMany inserts a list of strings into a tree, attempting to maintain
// balance as it does so.  A new tree is returned.
func (t *tree) insertMany(ss []string) (*tree, error) {
	// Create a sorted version of the list of strings then call our helper
	// function.
	sss := make([]string, len(ss))
	for i, s := range ss {
		sss[i] = s
	}
	sort.Strings(sss)
	return t.insertManySorted(sss)
}

// insertManySorted inserts a sorted list of strings into a tree, attempting to
// maintain balance as it does so.  It performs most of the work for
// insertMany.  A new tree is returned.
func (t *tree) insertManySorted(ss []string) (*tree, error) {
	// Handle the easy cases first.
	n := len(ss)
	switch n {
	case 0:
		return nil, nil
	case 1:
		return t.insert(ss[0])
	}

	// Insert the middle element, then recursively insert the left and
	// right sub-slices.
	mid := n / 2
	tNew, err := t.insert(ss[mid])
	if err != nil {
		return nil, err
	}
	if mid > 0 {
		tNew, err = tNew.insertManySorted(ss[:mid])
		if err != nil {
			return nil, err
		}
	}
	if mid+1 < n {
		tNew, err = tNew.insertManySorted(ss[mid+1:])
		if err != nil {
			return nil, err
		}
	}
	return tNew, nil
}

// walk invokes a given function on each node of a tree.
func (t *tree) walk(f func(*tree)) {
	if t == nil {
		return
	}
	t.left.walk(f)
	f(t)
	t.right.walk(f)
}

// findString returns the node containing a given string or nil if not found.
func (t *tree) find(s string) *tree {
	for t != nil {
		switch {
		case s == t.str:
			return t
		case s < t.str:
			t = t.left
		case s > t.str:
			t = t.right
		}
	}
	return nil
}
