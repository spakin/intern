// This file provides examples of intern usage.

package intern_test

import (
	"container/heap"
	"fmt"
	"math/rand"

	"github.com/spakin/intern"
)

// StringQueue represents a priority queue of strings.
type StringQueue struct {
	s  intern.LGESlice // Slice of LGE symbols
	ps []string        // Newly pushed strings
}

// Len returns the length of a string priority queue.
func (q *StringQueue) Len() int { return q.s.Len() }

// Less says if one element of a string priority queue is less than another.
func (q *StringQueue) Less(i, j int) bool { return q.s.Less(i, j) }

// Swap swaps two elements of a string priority queue.
func (q *StringQueue) Swap(i, j int) { q.s.Swap(i, j) }

// Push adds a new string to a string priority queue.
func (q *StringQueue) Push(x interface{}) {
	// Don't push anything yet.  Merely inform the intern package
	// that we intend to convert x to an LGE symbol.
	str := x.(string)
	intern.PreLGE(str)
	q.ps = append(q.ps, str)
}

// Pop returns the last string in a string priority queue.
func (q *StringQueue) Pop() interface{} {
	// Convert our newly added strings to symbols and enqueue them.
	for _, p := range q.ps {
		sym, err := intern.NewLGE(p)
		if err != nil {
			panic(err)
		}
		q.s = append(q.s, sym)
	}
	q.ps = q.ps[:0]
	heap.Init(q)

	// Return the final element of the slice as a string.
	n := len(q.s)
	x := q.s[n-1]
	q.s = q.s[:n-1]
	return x.String()
}

// generateString produces a semi-random string.  It mocks up whatever
// process would normally drive inserting strings into a priority
// queue.
func generateString() string {
	return fmt.Sprintf("My favorite number is %03d.", rand.Intn(100))
}

func Example_heap() {
	// Initialize a priority queue.
	h := &StringQueue{}
	heap.Init(h)

	// Insert five strings, each of which will be interned to a symbol.
	for _, s := range []string{
		"Isengrim",
		"Isumbras",
		"Ferumras",
		"Bandobras",
		"Gerontius",
	} {
		heap.Push(h, s)
	}

	// Alternate pushing five more strings and popping and
	// printing the alphabetically first five strings.
	for _, s := range []string{
		"Hildigrim",
		"Adalgrim",
		"Bungo",
		"Mirabella",
		"Gorbadoc",
	} {
		fmt.Println(heap.Pop(h).(string))
		heap.Push(h, s)
	}

	// Pop and print the remaining five strings in alphabetical order.
	for i := 0; i < 5; i++ {
		fmt.Println(heap.Pop(h).(string))
	}
}
// Output:
// Isengrim
// Hildigrim
// Gerontius
// Bungo
// Mirabella
// Gorbadoc
// Isumbras
// Ferumras
// Bandobras
// Adalgrim
