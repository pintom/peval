package eval

import "strings"

const (
	spades = iota
	hearts
	diamonds
	clubs
)

type cb struct {
	cards string
	suit  int
	bits  uint
}

// Card bits is an integer, made up of four bytes.  The high-order
// bytes are used to hold the rank bit pattern, whereas
// the low-order bytes hold the suit/rank/prime value
// of the card.
// +--------+--------+--------+--------+
// |xxxbbbbb|bbbbbbbb|cdhsrrrr|xxpppppp|
// +--------+--------+--------+--------+
//
// p = prime number of rank (deuce=2,trey=3,four=5,five=7,...,ace=41)
// r = rank of card (deuce=0,trey=1,four=2,five=3,...,ace=12)
// cdhs = suit of card
// b = bit turned on depending on rank of card
// ©Kevin Suffecool http://suffe.cool/poker/evaluator.html
var cardsTable = map[string]cb{
	"a": {cards: "2c", suit: clubs, bits: 98306},
	"A": {cards: "2h", suit: hearts, bits: 73730},
	"b": {cards: "3c", suit: clubs, bits: 164099},
	"B": {cards: "3h", suit: hearts, bits: 139523},
	"c": {cards: "4c", suit: clubs, bits: 295429},
	"C": {cards: "4h", suit: hearts, bits: 270853},
	"d": {cards: "5c", suit: clubs, bits: 557831},
	"D": {cards: "5h", suit: hearts, bits: 533255},
	"e": {cards: "6c", suit: clubs, bits: 1082379},
	"E": {cards: "6h", suit: hearts, bits: 1057803},
	"f": {cards: "7c", suit: clubs, bits: 2131213},
	"F": {cards: "7h", suit: hearts, bits: 2106637},
	"g": {cards: "8c", suit: clubs, bits: 4228625},
	"G": {cards: "8h", suit: hearts, bits: 4204049},
	"h": {cards: "9c", suit: clubs, bits: 8423187},
	"H": {cards: "9h", suit: hearts, bits: 8398611},
	"i": {cards: "Tc", suit: clubs, bits: 16812055},
	"I": {cards: "Th", suit: hearts, bits: 16787479},
	"j": {cards: "Jc", suit: clubs, bits: 33589533},
	"J": {cards: "Jh", suit: hearts, bits: 33564957},
	"k": {cards: "Qc", suit: clubs, bits: 67144223},
	"K": {cards: "Qh", suit: hearts, bits: 67119647},
	"l": {cards: "Kc", suit: clubs, bits: 134253349},
	"L": {cards: "Kh", suit: hearts, bits: 134228773},
	"m": {cards: "Ac", suit: clubs, bits: 268471337},
	"M": {cards: "Ah", suit: hearts, bits: 268446761},
	"n": {cards: "2d", suit: diamonds, bits: 81922},
	"N": {cards: "2s", suit: spades, bits: 69634},
	"o": {cards: "3d", suit: diamonds, bits: 147715},
	"O": {cards: "3s", suit: spades, bits: 135427},
	"p": {cards: "4d", suit: diamonds, bits: 279045},
	"P": {cards: "4s", suit: spades, bits: 266757},
	"q": {cards: "5d", suit: diamonds, bits: 541447},
	"Q": {cards: "5s", suit: spades, bits: 529159},
	"r": {cards: "6d", suit: diamonds, bits: 1065995},
	"R": {cards: "6s", suit: spades, bits: 1053707},
	"s": {cards: "7d", suit: diamonds, bits: 2114829},
	"S": {cards: "7s", suit: spades, bits: 2102541},
	"t": {cards: "8d", suit: diamonds, bits: 4212241},
	"T": {cards: "8s", suit: spades, bits: 4199953},
	"u": {cards: "9d", suit: diamonds, bits: 8406803},
	"U": {cards: "9s", suit: spades, bits: 8394515},
	"v": {cards: "Td", suit: diamonds, bits: 16795671},
	"V": {cards: "Ts", suit: spades, bits: 16783383},
	"w": {cards: "Jd", suit: diamonds, bits: 33573149},
	"W": {cards: "Js", suit: spades, bits: 33560861},
	"x": {cards: "Qd", suit: diamonds, bits: 67127839},
	"X": {cards: "Qs", suit: spades, bits: 67115551},
	"y": {cards: "Kd", suit: diamonds, bits: 134236965},
	"Y": {cards: "Ks", suit: spades, bits: 134224677},
	"z": {cards: "Ad", suit: diamonds, bits: 268454953},
	"Z": {cards: "As", suit: spades, bits: 268442665},
}

// Cards represent cards as follows:
// a - 2c
// A - 2h
// b - 3c
// B - 3h
// c - 4c
// C - 4h
// d - 5c
// D - 5h
// e - 6c
// E - 6h
// f - 7c
// F - 7h
// g - 8c
// G - 8h
// h - 9c
// H - 9h
// i - Tc
// I - Th
// j - Jc
// J - Jh
// k - Qc
// K - Qh
// l - Kc
// L - Kh
// m - Ac
// M - Ah
// n - 2d
// N - 2s
// o - 3d
// O - 3s
// p - 4d
// P - 4s
// q - 5d
// Q - 5s
// r - 6d
// R - 6s
// s - 7d
// S - 7s
// t - 8d
// T - 8s
// u - 9d
// U - 9s
// v - Td
// V - Ts
// w - Jd
// W - Js
// x - Qd
// X - Qs
// y - Kd
// Y - Ks
// z - Ad
// Z - As
type Cards string

// String method return card values in human readable form, e.g."wn" would return Jd2d.
// It implements Stringer interface from fmt package.
func (c Cards) String() string {
	var s []string
	for _, v := range c {
		v, ok := cardsTable[string(v)]
		if !ok {
			return "error"
		}
		s = append(s, v.cards)
	}
	st := strings.Join(s, "")

	return st
}

// Slice method separates given cards string ard returns a slice of single cards.
func (c Cards) Slice() []Cards {
	st := strings.Split(string(c), "")

	cst := make([]Cards, len(st))
	for i, v := range st {
		cst[i] = Cards(v)
	}

	return cst
}

// CardsBits method returns a slice of cards bits
func (c Cards) CardsBits() []uint {
	s := make([]uint, len(c))
	for i, v := range c {
		s[i] = cardsTable[string(v)].bits
	}
	return s
}

// only for two cards!!
func (c Cards) Suited() bool {
	if cardsTable[string(c[0])].suit == cardsTable[string(c[1])].suit {
		return true
	}

	return false
}

// Sorts only 2 cards!!!!
func (c Cards) Sort() Cards {
	if cardsTable[string(c[0])].bits > cardsTable[string(c[1])].bits {
		return c
	}

	return Cards(c[1]) + Cards(c[0])
}

func TwoCardEqv(cards []Cards) map[string]int {
	cardsMap := make(map[string]int)
	for _, v := range cards {
		v = v.Sort()
		w := v.String()
		var s []byte
		s = append(s, w[0], w[2])
		if v.Suited() {
			s = append(s, byte("s"[0]))
		}

		cardsMap[string(s)]++
	}

	return cardsMap
}

func CardsEqv(cards Cards) string {
	cards = cards.Sort()
	w := cards.String()
	var s []byte
	s = append(s, w[0], w[2])
	if cards.Suited() {
		s = append(s, byte("s"[0]))
	}

	return string(s)
}
