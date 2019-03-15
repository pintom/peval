package eval

import (
	"strings"
)

// BetterHands takes in cards dealt cards (max 7)
// and returns the map of hand combinations that beat the owners cards.
func BetterHands(dealtCards Cards) (bestHand Cards, betterHands map[string][]Cards,
	outsOuts map[string][]Cards) {
	// Check if there are enough or not too many cards to determine the value.
	if (len(dealtCards) < 5) || (len(dealtCards) > 7) {
		return "", nil, nil
	}

	// Populate the deck minus the cards that already been dealt.
	de := make(map[string]cb)
	for x, y := range cardsTable {
		de[x] = y
	}
	for _, y := range dealtCards {
		delete(de, string(y))
	}

	//dealt cards value in bits.
	// !!! FIRST TWO CARDS ARE REPLACED LATER IN THE FUNCTION !!!
	dealtBits := dealtCards.CardsBits()
	_, bestHandValue, bestHandInt := eval(dealtBits)
	bestHand = BitsToChars(bestHandInt)

	// Make a map of cards that opponents might hold that beats me.
	betterHands = make(map[string][]Cards)
	// holds Outs that would make nuts
	outsOuts = make(map[string][]Cards)

	var outs []uint
	var opOuts []uint
	// Add two cards as 2 more will be dealt
	if len(dealtCards) == 5 {
		outs = make([]uint, len(dealtCards)+2)
		opOuts = make([]uint, len(dealtCards)+2)
	} else { // add one more as only the river is left.
		outs = make([]uint, len(dealtCards)+1)
		opOuts = make([]uint, len(dealtCards)+1)

	}
	copy(outs, dealtBits)
	copy(opOuts, dealtBits)

	// Range over the left cards in a deck and check if they beat the bestHandValue
	for i, v := range de {
		var f int
		var d int
		var str string
		for j, w := range de {
			// if cards are the same, continue
			if i == j {
				continue
			}

			// Calculates possible opponents' cards
			dealtBits[0] = v.bits
			dealtBits[1] = w.bits

			// Check if this "new" hand beats the dealt one.
			str, f, _ = eval(dealtBits)
			if f < bestHandValue {
				betterHands[str] = append(betterHands[str], Cards(i+j))
			}

			// Calculate Outs for two more cards to go.
			if len(dealtCards) == 5 {

				// append cards to the end (5 and six are last positions)
				outs[5] = v.bits
				opOuts[5] = v.bits
				outs[6] = w.bits
				opOuts[6] = w.bits

				str, d, _ = eval(outs)
				if d < bestHandValue && d < f {
					outsOuts[str] = append(outsOuts[str], Cards(i+j))
				}

			}

		}
		// Calculate Outs for one more card to go.
		if len(dealtCards) == 6 {
			outs[6] = v.bits

			str, d, _ := eval(outs)
			if d < bestHandValue {
				outsOuts[str] = append(outsOuts[str], Cards(i))
			}

		}
		delete(de, i)

	}

	return

}

// HandValue returns a (value, unique int value, Best cards) of given 5 to 7 cards.
func HandValue(c Cards) (string, int, Cards) {
	unt := c.CardsBits()

	st, in, sl := eval(unt)

	return st, in, BitsToChars(sl)
}

// BitsToChars takes in a slice of cards bits and returns Cards value.
func BitsToChars(uns []uint) Cards {
	sl := make([]string, 2)
	for _, v := range uns {
		for x, y := range cardsTable {
			if y.bits == v {
				sl = append(sl, x)
			}
		}
	}

	s := strings.Join(sl, "")
	return Cards(s)
}

func eval(hand []uint) (string, int, []uint) {
	if len(hand) < 5 {
		return "", 0, nil
	} else if len(hand) > 7 {
		return "", 0, nil
	}
	// Checks if cards do not repeat.
	for ind, val := range hand {
		for _, val2 := range hand[ind+1:] {
			if val2 == val {
				return "", 0, nil
			}

		}

	}

	var best int
	bestCards := make([]uint, 5)

	copy(bestCards, hand[:5])
	c := make([]uint, len(hand))

	if len(hand) == 5 {
		best = ev(hand)

		return handValue(best), best, bestCards
	}
	if len(hand) == 6 {
		// if there are more than 3 community cards
		best = ev(hand[:5])
		for i := 0; i < 5; i++ {
			copy(c, hand)
			c[i] = hand[5]
			b := ev(c[:5])
			if b < best { // less is better
				best = b
				copy(bestCards, c[:5])
			}
		}
		return handValue(best), best, bestCards

	} else if len(hand) == 7 {
		best = ev(hand[:5])
		for i := 0; i < 5; i++ {
			copy(c, hand)
			c[i] = hand[5]
			b := ev(c[:5])
			if b < best { // less is better
				best = b
				copy(bestCards, c[:5])

			}

		}
		for i := 0; i < 5; i++ {
			copy(c, hand)
			c[i] = hand[6]
			b := ev(c[:5])
			if b < best { // less is better
				best = b
				copy(bestCards, c[:5])

			}

		}
		for i := 0; i < 5; i++ {
			for j := i + 1; j < 5; j++ {
				copy(c, hand)
				c[i] = hand[5]
				c[j] = hand[6]
				b := ev(c[:5])
				if b < best { // less is better
					best = b
					copy(bestCards, c[:5])

				}
			}
		}

		return handValue(best), best, bestCards
	}

	return "", 0, bestCards
}

func ev(cards []uint) int {
	// Shift by 16 bits to the right
	q := (cards[0] | cards[1] | cards[2] | cards[3] | cards[4]) >> 16

	// Check if all suits are the same
	fl := cards[0] & cards[1] & cards[2] & cards[3] & cards[4] & 0xF000

	var v int

	if fl != 0 {
		// it it is not 0, then it is a Flush
		v = flushes[q]
	} else {
		// if it is 0, then it is either Straight or High cards
		v = unique5[q]

		// if unique5 returns 0, then it holds actual hand's distinct value.
		if v == 0 {
			d := (cards[0] & 0xFF) * (cards[1] & 0xFF) * (cards[2] & 0xFF) * (cards[3] & 0xFF) * (cards[4] & 0xFF)
			v = findIt(int(d))
			return values[v]
		}
	}

	return v
}

func findIt(key int) int {
	low := 0
	high := 4887
	mid := 0

	for low <= high {
		mid = (high + low) >> 1 // divide by two
		if key < products[mid] {
			high = mid - 1
		} else if key > products[mid] {
			low = mid + 1
		} else {
			return mid
		}
	}
	return -1
}

func handValue(val int) string {
	if val > 6185 {
		return "High cards"
	}
	if val > 3325 {
		return "One Pair"
	}
	if val > 2467 {
		return "Two Pairs"
	}
	if val > 1609 {
		return "Three of a Kind"
	}
	if val > 1599 {
		return "Straight"
	}
	if val > 322 {
		return "Flush"
	}
	if val > 166 {
		return "Full House"
	}
	if val > 10 {
		return "Four of a Kind"
	}
	if val > 0 {
		return "Straight Flush"
	}
	return ""
}
