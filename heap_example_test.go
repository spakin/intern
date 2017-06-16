// This file shows how to implement a priority queue of LGEs.

package intern_test

import (
	"container/heap"
	"fmt"

	"github.com/spakin/intern"
)

// Define a priority queue of LGEs.
type SymQ []intern.LGE

// Len returns the length of an LGE priority queue.
func (sq SymQ) Len() int { return len(sq) }

// Less says if one LGEs precedes another.
func (sq SymQ) Less(i, j int) bool { return sq[i] < sq[j] }

// Swap swaps two LGEs.
func (sq SymQ) Swap(i, j int) { sq[i], sq[j] = sq[j], sq[i] }

// Push interns a string and pushes the result on the priority queue.
func (sq *SymQ) Push(x interface{}) {
	xStr := x.(string)
	sym, err := intern.NewLGE(xStr)
	if err != nil {
		// We ran out of LGEs.  Forget all existing LGEs and
		// start over with whatever is currently in the queue.
		strs := make([]string, len(*sq))
		for i, s := range *sq {
			strs[i] = s.String()
		}
		intern.ForgetAllLGEs()
		for _, s := range strs {
			intern.PreLGE(s)
		}
		intern.PreLGE(xStr)
		for i, s := range strs {
			(*sq)[i], err = intern.NewLGE(s)
			if err != nil {
				panic(err)
			}
		}
		sym, err = intern.NewLGE(xStr)
		if err != nil {
			panic(err)
		}
	}
	*sq = append(*sq, sym)
}

// Pop pops an LGE from the priority queue and returns it as a string.
func (sq *SymQ) Pop() interface{} {
	ns := len(*sq)
	s := (*sq)[ns-1]
	*sq = (*sq)[:ns-1]
	return s.String()
}

// Sort a list of strings using a priority queue.
func ExampleForgetAllLGEs() {
	// Define some strings in reverse alphabetical order because
	// this is a worst case for LGEs.  It forces the use of the
	// control path that forgets all symbols then remaps all
	// currently queued strings to LGEs.
	colors := []string{
		"yellow green",
		"yellow",
		"white smoke",
		"white",
		"wheat",
		"violet red",
		"violet",
		"turquoise",
		"tomato",
		"thistle",
		"tan",
		"steel blue",
		"spring green",
		"snow",
		"slate grey",
		"slate gray",
		"slate blue",
		"sky blue",
		"sienna",
		"seashell",
		"sea green",
		"sandy brown",
		"salmon",
		"saddle brown",
		"royal blue",
		"rosy brown",
		"red",
		"purple",
		"powder blue",
		"plum",
		"pink",
		"peru",
		"peach puff",
		"papaya whip",
		"pale violet red",
		"pale turquoise",
		"pale green",
		"pale goldenrod",
		"orchid",
		"orange red",
		"orange",
		"olive drab",
		"old lace",
		"navy blue",
		"navy",
		"navajo white",
		"moccasin",
		"misty rose",
		"mint cream",
		"midnight blue",
		"medium violet red",
		"medium turquoise",
		"medium spring green",
		"medium slate blue",
		"medium sea green",
		"medium purple",
		"medium orchid",
		"medium blue",
		"medium aquamarine",
		"maroon",
		"magenta",
		"linen",
		"lime green",
		"light yellow",
		"light steel blue",
		"light slate grey",
		"light slate gray",
		"light slate blue",
		"light sky blue",
		"light sea green",
		"light salmon",
		"light pink",
		"light grey",
		"light green",
		"light gray",
		"light goldenrod yellow",
		"light goldenrod",
		"light cyan",
		"light coral",
		"light blue",
		"lemon chiffon",
		"lawn green",
		"lavender blush",
		"lavender",
		"khaki",
		"ivory",
		"indian red",
		"hot pink",
		"honeydew",
		"grey",
		"green yellow",
		"green",
		"gray",
		"goldenrod",
		"gold",
		"ghost white",
		"gainsboro",
		"forest green",
		"floral white",
		"firebrick",
		"dodger blue",
		"dim grey",
		"dim gray",
		"deep sky blue",
		"deep pink",
		"dark violet",
		"dark turquoise",
		"dark slate grey",
		"dark slate gray",
		"dark slate blue",
		"dark sea green",
		"dark salmon",
		"dark red",
		"dark orchid",
		"dark orange",
		"dark olive green",
		"dark magenta",
		"dark khaki",
		"dark grey",
		"dark green",
		"dark gray",
		"dark goldenrod",
		"dark cyan",
		"dark blue",
		"cyan",
		"cornsilk",
		"cornflower blue",
		"coral",
		"chocolate",
		"chartreuse",
		"cadet blue",
		"burlywood",
		"brown",
		"blue violet",
		"blue",
		"blanched almond",
		"black",
		"bisque",
		"beige",
		"azure",
		"aquamarine",
		"antique white",
		"alice blue",
	}

	// Insert all of the strings into a priority queue.
	sq := &SymQ{}
	heap.Init(sq)
	for _, c := range colors {
		heap.Push(sq, c)
	}

	// Read out all the strings in alphabetical order.
	for sq.Len() > 0 {
		c := heap.Pop(sq).(string)
		fmt.Println(c)
	}
}
