// This file provides helper routines needed by multiple tests.

package intern_test

import (
	"fmt"
	"math/rand"
	"unicode/utf8"
)

// charSet is the list of characters from which to draw for randomly generated
// strings.
var charSet []rune

// init initializes our global state.
func init() {
	const cs = "ÂBÇDÈFGHÍJKLMÑÖPQRSTÛVWXÝZ0123456789âbçdèfghíjklmñöpqrstûvwxýz @#$*-+<>一二三"
	charSet = make([]rune, 0, utf8.RuneCountInString(cs))
	for _, r := range cs {
		charSet = append(charSet, r)
	}
}

// randomString returns a random string of a given length.
func randomString(r *rand.Rand, n int) string {
	rs := make([]rune, n)
	nc := len(charSet)
	for i := range rs {
		rs[i] = charSet[r.Intn(nc)]
	}
	return string(rs)
}

// reverseString returns a string with its characters reversed.
func reverseString(s string) string {
	rs := []rune(s)
	nr := len(rs)
	for i := nr / 2; i >= 0; i-- {
		rs[i], rs[nr-i-1] = rs[nr-i-1], rs[i]
	}
	return string(rs)
}

// Dummy is used to prevent benchmarks from being treated as dead code.
var Dummy uint64

// nComp is the number of strings to compare all the others to when benchmarking.
const nComp = 1000

// generateSimilarStrings generates a list of strings that have a substantial
// prefix in common.
func generateSimilarStrings(n int) []string {
	strs := make([]string, n)
	for i := range strs {
		strs[i] = fmt.Sprintf("String comparisons can be slow when the strings to compare have a long prefix in common.  My favorite number is now %015d.", i+1)
	}
	return strs
}

// generateRandomStrings generates a list of long strings with random content.
func generateRandomStrings(n int) []string {
	prng := rand.New(rand.NewSource(1920)) // Constant for reproducibility
	strs := make([]string, n)
	for i := range strs {
		nc := prng.Intn(50) + 10 // Number of characters
		strs[i] = randomString(prng, nc)
	}
	return strs
}

// ozChars is a list of characters in the Oz series of books.
var ozChars = []string{
	"A-B-Sea Serpent",
	"Abatha",
	"Agnes",
	"Army of Oogaboo",
	"Aunt Em",
	"Belfaygor of Bourne",
	"Bell-snickle",
	"Betsy Bobbin",
	"Billina",
	"Blinkem",
	"Boq",
	"Bristle",
	"Button-Bright",
	"Cap'n Bill",
	"Cayke",
	"China Princess",
	"Chiss",
	"Chopfyt",
	"Cinnamon Bunn",
	"Cowardly Lion",
	"Dorothy Gale",
	"Dr. Pipt",
	"Ervic",
	"Eureka",
	"Evoldo",
	"Foolish Owl",
	"Frogman",
	"Fyter the Tin Soldier",
	"Gayelette",
	"Glass Cat",
	"Glinda",
	"Good Witch of the North",
	"Good Witch of the South",
	"Graham Gems",
	"Great Royal Marshmallow",
	"Guardian of the Gates",
	"Herby",
	"Hungry Tiger",
	"Jack Pumpkinhead",
	"Jellia Jamb",
	"Jenny Jump",
	"Jester",
	"Jim the Cab-Horse",
	"Jinjur",
	"Jinnicky the Red Jinn",
	"John Dough",
	"Johnny Cake",
	"Johnny Dooit",
	"Kabumpo",
	"Kalidah",
	"Kaliko",
	"King Kinda Jolly",
	"King Kleaver",
	"King Krewl",
	"King Kynd",
	"King Pastoria",
	"King of Bunnybury",
	"King of the Fairy Beavers",
	"Ku-Klip",
	"Lavender Bear",
	"Lonesome Duck",
	"Mombi",
	"Mr. Muffin",
	"Mr. Yoop",
	"Mrs. Yoop",
	"Munchkins",
	"Nimmie Amee",
	"Nome King",
	"Ojo the Lucky",
	"Patchwork Girl",
	"Phonograph",
	"Polychrome",
	"Pop Over",
	"Prince Karver",
	"Princess Langwidere",
	"Princess Ozma",
	"Professor Woggle-Bug",
	"Queen Ann Soforth",
	"Queen Coo-ee-oh",
	"Queen Lurline",
	"Rak",
	"Robin Brown",
	"Sally Lunn",
	"Sawhorse",
	"Scarecrow",
	"Shaggy Man",
	"Sir Hokus of Pokes",
	"Smith",
	"Soldier with the Green Whiskers",
	"The Gump",
	"Tik-Tok",
	"Tin Woodman",
	"Tinker",
	"Tititi-Hoochoo",
	"Toto",
	"Trot",
	"Tugg",
	"Ugu the Shoemaker",
	"Unc Nunkie",
	"Uncle Henry",
	"Wicked Witch of the East",
	"Wicked Witch of the North",
	"Wicked Witch of the South",
	"Wicked Witch of the West",
	"Wise Donkey",
	"Wiser the Owl",
	"Wizard of Oz",
	"Woozy",
	"Zeb Hugson",
}
