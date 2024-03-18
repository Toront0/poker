package game

import (
	"strconv"
	"sort"
	"slices"
	"fmt"

	"github.com/Toront0/poker/internal/utils"

)

func IsFlush(card1, card2 string, deck []string) bool {

	firstCardMatches := 0
	secondCardMatches := 0


	for _, p := range deck {


		if CardSuit(p) == CardSuit(card1) {
			firstCardMatches++
			continue
		}

		if CardSuit(p) == CardSuit(card2) {
			secondCardMatches++
			continue
		}


	} 

	if CardSuit(card1) == CardSuit(card2) {

		if firstCardMatches + secondCardMatches + 2 >= 5 {
			return true
		} else {
			return false
		}

	} else {
		if firstCardMatches >= 4 || secondCardMatches >= 4 {
			return true
		} else {
			return false
		}
	}
}

// 10clubs", "7clubs", "12clubs", "10hearts", "4spades

func FindHigherFlush(cards []string, players []PokerPlayer) []int {

	matches := 0
	flushBasedOnSuit := ""

	for i := 0; i < len(cards); i++ {

		for j := i + 1; j < len(cards); j++ {
			if CardSuit(cards[i]) == CardSuit(cards[j]) {
				matches++
			} 

			


			if matches >= 2 {
				flushBasedOnSuit = CardSuit(cards[i])
				break
			}

		}

		matches = 0
	}

	playerFlushCards := []int{}

	for _, p := range players {
		


		if CardSuit(p.Hand[0]) == flushBasedOnSuit {
			if CardValue(p.Hand[0]) == 1 {
				return []int{p.ID}
			}


			playerFlushCards = append(playerFlushCards, CardValue(p.Hand[0]))

		}
		if CardSuit(p.Hand[1]) == flushBasedOnSuit {
			if CardValue(p.Hand[1]) == 1 {
				return []int{p.ID}
			}

			playerFlushCards = append(playerFlushCards, CardValue(p.Hand[1]))

		}
	}

	highestPlayerFlushCard := slices.Max(playerFlushCards)

	playerWin := utils.Filter(players, func (p PokerPlayer) bool { return slices.Contains(p.Hand, strconv.Itoa(highestPlayerFlushCard) + flushBasedOnSuit) })


	return utils.Map(playerWin, func (i PokerPlayer) int { return i.ID })

}

func IsFlushRoyal(card1, card2 string, deck []string) string {

	res := []string{}


	IsFlushRoyal := "flush-royal"

	lowCardsAmount := 0

	if CardValue(card1) > 1 && CardValue(card1) < 10 && CardValue(card2) > 1 && CardValue(card2) < 10  {
		return ""
	}

	if CardValue(card1) == 1 || CardValue(card1) >= 10 {
		res = append(res, card1)
	}

	if CardValue(card2) == 1 || CardValue(card2) >= 10 {
		res = append(res, card2)
	}


	for _, p := range deck {

		if lowCardsAmount > 2 {
			return ""
		}

		if CardValue(p) == 1 || CardValue(p) > 10 {
			
			res = append(res, p)
		}

	} 

	if len(res) < 5 {
		return ""
	}

	flushSuit := CardSuit(res[0])

	for _, p := range res {

		if CardSuit(p) != flushSuit {
			return ""
		}

	}

	return IsFlushRoyal
}

func IsStraight(card1, card2 string, deck []string) bool {

	cards := []int{CardValue(card1), CardValue(card2), CardValue(deck[0]), CardValue(deck[1]), CardValue(deck[2]), CardValue(deck[3]), CardValue(deck[4])}

	sort.Slice(cards, func(i, j int) bool { return cards[i] < cards[j] })

	
	if cards[6] == 13 && cards[5] == 12 && cards[4] == 11 && cards[3] == 10 && cards[0] == 1 {

		return true

	}	


	matches := 0

	for i := 0; i < len(cards); i++ {

		if i + 1 >= len(cards) {
			break
		}
		
		if matches >= 4 {
			break
		}

		if cards[i] - cards[i + 1] == -1 {

			matches++

		} else {

			matches = 0

		}

	}


	if matches >= 4 {
		return true
	}  else {
		return false
	}
}

func FindHigherStraight(tableCards []string, players []PokerPlayer) []int {

	matches := 0

	cards := []int{CardValue(tableCards[0]), CardValue(tableCards[1]), CardValue(tableCards[2]), CardValue(tableCards[3]), CardValue(tableCards[4])}

	for _, p := range players {

		cards = append(cards, CardValue(p.Hand[0]))
		cards = append(cards, CardValue(p.Hand[1]))

	}

	



	sort.Slice(cards, func(i, j int) bool { return cards[i] < cards[j] })

	straight := []int{}

	


	for i := len(cards) - 1; i > 0; i-- {
		if i - 1 < 0 {
			break
		}

		if matches == 4 {
			straight = append(straight, cards[i])
			
			break

		}

		if cards[i] - cards[i - 1] == 0 {

			continue
		}


		if cards[i] - cards[i - 1] != 1 {
			straight = []int{}
			matches = 0
			continue
		} else {
		

			straight = append(straight, cards[i])
		
			matches++

		} 
	
	}

	fmt.Println("cards", cards)

	// Exclusive case when value 1 as a highest value after 13 
	if straight[0] == 13 && straight[1] == 12 && straight[2] == 11 && straight[3] == 10 && cards[0] == 1 {

		straight = []int{14,13,12,11,10}

	} 


	res := []int{}

	fmt.Println("straight", straight)
	fmt.Println("res", res)

	for _, p := range players {
		card1 := CardValue(p.Hand[0])
		card2 := CardValue(p.Hand[1])

		if card1 == 1 {
			card1 = 14
		}

		if card2 == 1 {
			card2 = 14
		}

		if card1 == straight[0] || card2 == straight[0] {

			res = append(res, p.ID)

		} 

	}

	if len(res) == 0 {

		for _, p := range players {

			res = append(res, p.ID)

		}

	}


	return res

}



type FullHouseHelper struct {
	Triple int
	Double int
}

// To Find higher full-hosue, first we check highest card in there trips combination and then two last cards will be checked.
func FindHigherFullHouse(cards []string, players []PokerPlayer) []int {

	tripleMaxValue := 0
	pairMaxValue := 0

	mapRes := make(map[int]FullHouseHelper)

	for _, p := range players {
		fMatches := 0
		sMatches := 0


		for _, c := range cards {

			if CardValue(p.Hand[0]) == CardValue(c) {
				fMatches++
			}

			if CardValue(p.Hand[1]) == CardValue(c) {
				sMatches++
			}

		}

		if fMatches >= 2 && CardValue(p.Hand[0]) > tripleMaxValue {

			tripleMaxValue = CardValue(p.Hand[0])
		}

		if sMatches >= 2 && CardValue(p.Hand[1]) > tripleMaxValue {


			tripleMaxValue = CardValue(p.Hand[1])
		}

		if fMatches == 1 && CardValue(p.Hand[0]) > pairMaxValue {


			pairMaxValue = CardValue(p.Hand[0])
		}

		if sMatches == 1 && CardValue(p.Hand[1]) > pairMaxValue {


			pairMaxValue = CardValue(p.Hand[1])
		}

		

	}


	res := []int{}

	for _, p := range players {

		maxV := CardValue(p.Hand[0])

		if CardValue(p.Hand[1]) > CardValue(p.Hand[0]) {
			maxV = CardValue(p.Hand[1])
		}


		if maxV >= tripleMaxValue {

			fmt.Println("mapRes[p.ID].Triple", mapRes[p.ID].Triple)

			res = append(res, p.ID)

		}


	}


	if len(res) > 1 {
		res := []int{}

		for _, p := range players {

			minV := CardValue(p.Hand[0])

			if CardValue(p.Hand[1]) < CardValue(p.Hand[0]) {
				minV = CardValue(p.Hand[1])
			}

			if minV >= pairMaxValue {
				res = append(res, p.ID)
			}

		}
		
		return res

	}

	return res

}

func IsFullHouse(card1, card2 string, deck []string) bool {

	fMatches := 0 
	sMatches := 0 

	for _, c := range deck {

		if CardValue(card1) == CardValue(c) {
			fMatches++
		}

		if CardValue(card2) == CardValue(c) {
			sMatches++
		}

	}

	if fMatches > 1 && sMatches == 1 || fMatches == 1 && sMatches > 1 {

		return true

	}

	return false

}


func IsStraightFlush(card1, card2 string, deck []string) bool {

	f := IsFlush(card1, card2, deck)

	if !f {
		return false
	}


	s := IsStraight(card1, card2, deck)

	if !s {
		return false
	}

	return true
}

func IsFourOfAKind(card1, card2 string, deck []string) bool {

	fMatches := 0
	sMatches := 0

	for _, p := range deck {

		if CardValue(card1) == CardValue(p) {
			fMatches++
		}

		if CardValue(card2) == CardValue(p) {
			sMatches++
		}
	}


	if CardValue(card1) == CardValue(card2) {
		if fMatches >= 2 || sMatches >= 2 {
			return true
		} else {
			return false
		}
	} else {
		if fMatches >= 3 || sMatches >= 3 {
			return true
		} else {
			return false
		}
	}
}

func FindHigherFourOfAKind(cards []string, players []PokerPlayer) []int {

	max := 0

	for _, p := range players {

		fMatches := 0
		sMatches := 0

		alreadyPair := CardValue(p.Hand[0]) == CardValue(p.Hand[1])

		for _, c := range cards {

			if CardValue(p.Hand[0]) == CardValue(c) {
				fMatches++
			}

			if CardValue(p.Hand[1]) == CardValue(c) {
				sMatches++
			}

		
		}

		if alreadyPair && (fMatches >= 2 || sMatches >= 2) {
			if CardValue(p.Hand[0]) == 1 || CardValue(p.Hand[1]) == 1 {

				max = 14

			}

			if CardValue(p.Hand[0]) >= max {
				max = CardValue(p.Hand[0])
			}

			if CardValue(p.Hand[1]) >= max {
				max = CardValue(p.Hand[1])
			}

			
		}

		if fMatches == 3 {

			if CardValue(p.Hand[0]) == 1 {

				max = 14

			} 

			if CardValue(p.Hand[0]) >= max {
				max = CardValue(p.Hand[0])
			}

		}

		if sMatches == 3 {

			if CardValue(p.Hand[1]) == 1 {

				max = 14

			} 

			if CardValue(p.Hand[1]) >= max {
				max = CardValue(p.Hand[1])
			}

		}

		fMatches = 0
		sMatches = 0
	}

	fmt.Println("max", max)

	res := []int{}

	for _, p := range players {
		card1 := CardValue(p.Hand[0]) 
		card2 := CardValue(p.Hand[1]) 

		if CardValue(p.Hand[0]) == 1 {
			card1 = 14
		}

		if CardValue(p.Hand[1]) == 1 {
			card2 = 14
		}

		if card1 == max || card2 == max {
			res = append(res, p.ID)
		}

	}

	return res

}



func IsSet(card1, card2 string, deck []string) bool {

	fMatches := 0
	sMatches := 0

	for _, p := range deck {

		if CardValue(card1) == CardValue(p) {
			fMatches++
		}

		if CardValue(card2) == CardValue(p) {
			sMatches++
		}

	}

	if CardValue(card1) == CardValue(card2) {
		if fMatches >= 1 || sMatches >= 1 {
			return true
		} else {
			return false
		}
	} else {
		if fMatches >= 2 || sMatches >= 2 {
			return true
		} else {
			return false
		}
	}

}

func FindHigherSet(cards []string, players []PokerPlayer) []int {
	max := 0
	topKicker := 0


	res := []int{}

	fMatches := 0
	sMatches := 0
		
	for _, p := range players {

		alreadyPair := CardValue(p.Hand[0]) == CardValue(p.Hand[1]) 



		for _, c := range cards {



			if CardValue(p.Hand[0]) == CardValue(c) {
				fMatches++
			}  

			if CardValue(p.Hand[1]) == CardValue(c) {
				sMatches++
			}

		}

		if alreadyPair && (fMatches > 0 || sMatches > 0) {
			max = CardValue(p.Hand[0])	
			fMatches = 0
			sMatches = 0
			continue
		}

		if fMatches >= 2 {
			if CardValue(p.Hand[0]) == 1 {
				max = 14
				
			}

			if CardValue(p.Hand[0]) >= max {
				max = CardValue(p.Hand[0])
			}

			if CardValue(p.Hand[1]) == 1 {
				topKicker = 14
			} 

			if CardValue(p.Hand[1]) >= topKicker {
				topKicker = CardValue(p.Hand[1])
			}
		}

		if sMatches >= 2 {
			if CardValue(p.Hand[1]) == 1 {
				max = 14
			}

			if CardValue(p.Hand[1]) >= max {
				max = CardValue(p.Hand[1])
			}

			if CardValue(p.Hand[0]) == 1 {
				topKicker = 14
			} 

			if CardValue(p.Hand[0]) >= topKicker {
				topKicker = CardValue(p.Hand[0])
			}
		}

		fMatches = 0
		sMatches = 0
	}


	for _, p := range players {
		if CardValue(p.Hand[0]) == 1 && fMatches >= 2 || CardValue(p.Hand[1]) == 1 && sMatches >= 2 {
			res = append(res, p.ID)
			continue
		} 

		if CardValue(p.Hand[0]) == max || CardValue(p.Hand[1]) == max {
			res = append(res, p.ID)
		} 

	}



	if len(res) > 1 {


		newRes := []int{}
		for _, p := range players {

			card1 := CardValue(p.Hand[0])
			card2 := CardValue(p.Hand[1])

			if card1 == 1 {
				card1 = 14
			}

			if card2 == 1 {
				card2 = 14
			}

			if card1 + card2 >=  max + topKicker {

				// r := utils.RemoveIndex(res, i)
				newRes = append(newRes, p.ID)
			} 

		}

		res = newRes

	} 

	return res

}

func IsTwoPair(card1, card2 string, deck []string) bool {

	fMatches := 0
	sMatches := 0

	for _, p := range deck {

		if CardValue(card1) == CardValue(p) {

			fMatches++

		}

		if CardValue(card2) == CardValue(p) {

			sMatches++

		}

	}


	if fMatches >= 1 && sMatches >= 1 {
		return true
	} else {
		return false
	}

}

func FindHigherTwoPair(cards []string, players []PokerPlayer) []int {

	max := 0

	for _, p := range players {
		if CardValue(p.Hand[0]) == 1 || CardValue(p.Hand[1]) == 1 {
			max = 14
		}


		if CardValue(p.Hand[0]) >= max {

			max = CardValue(p.Hand[0])

		}

		if CardValue(p.Hand[1]) >= max {

			max = CardValue(p.Hand[0])

		}

	}


	res := []int{}

	for _, p := range players {
		if CardValue(p.Hand[0]) == 1 || CardValue(p.Hand[1]) == 1 {
			res = append(res, p.ID)
			continue
		}

		if CardValue(p.Hand[0]) == max || CardValue(p.Hand[1]) == max {

			res = append(res, p.ID)

		} 

	} 

	return res

}

func IsPair(card1, card2 string, deck []string) bool {


	for _, p := range deck {

		if CardValue(card1) == CardValue(p) || CardValue(card2) == CardValue(p) {

			return true

		}

	}



	if CardValue(card1) == CardValue(card2) {
		return true
	}
	
	return false
}

func FindHigherPair(cards []string, players []PokerPlayer) []int {
	pairs := make(map[int]int)

	for _, p := range players {

		if CardValue(p.Hand[0]) == CardValue(p.Hand[1]) {
			pairs[p.ID] = CardValue(p.Hand[0])
			continue
		}

		for _, c := range cards {

			if CardValue(p.Hand[0]) == CardValue(c) {

				pairs[p.ID] = CardValue(c)

			}

			if CardValue(p.Hand[1]) == CardValue(c) {

				pairs[p.ID] = CardValue(c)

			}
		}
	}

	res := []int{}

	max := 0

	for _, c := range pairs {

		if c >= max {
			max = c
		}

	}

	for id, c := range pairs {
		if c == 1 {
			res = append(res, id)
			// Usually in system we do not represent ace as 14 value card, but in this case i thought it would not be that bad.
			max = 14

			continue
		}

		if c == max {
			res = append(res, id)
		}

	}
	return res
}


// func HighestCard(card1, card2 string) string {
	
// 	if CardValue(card1) == 1 {
// 		return card1
// 	}

// 	if CardValue(card2) == 1 {
// 		return card2
// 	}

// 	if CardValue(card1) >= CardValue(card2) {
// 		return card1
// 	} else {
// 		return card2
// 	}

// }


func CardSuit(card string) string {

	suitLetter := 2

	_, err := strconv.Atoi(string(card[1]))

	if err != nil {
		suitLetter = 1
	}



	return card[suitLetter:]


}


func CardValue(card string) int {

	suitLetter := 2

	_, err := strconv.Atoi(string(card[1]))

	if err != nil {
		suitLetter = 1
	}


	val, _ := strconv.Atoi(card[:suitLetter])

	return val
}

func DetermineWhosHandHigher(cards []string, players []PokerPlayer) []int {
	res := []int{}

	r := &WinnerCombinations{}

	for i, p := range players {

		fmt.Println("checking player", i)

		if v := IsFlushRoyal(p.Hand[0], p.Hand[1], cards); v == "flush-royal" {
			r.FlushRoyal = append(r.FlushRoyal, p.ID)
			return r.FlushRoyal
		}

		if v := IsStraightFlush(p.Hand[0], p.Hand[1], cards); v {
			r.StraighFlush = append(r.StraighFlush, p.ID)
			continue
		}

		if v := IsFourOfAKind(p.Hand[0], p.Hand[1], cards); v {
			r.FourOfAKind = append(r.FourOfAKind, p.ID)
			continue
		}

		if v := IsFullHouse(p.Hand[0], p.Hand[1], cards); v {
			r.FullHouse = append(r.FullHouse, p.ID)
			continue
		}

		if v := IsFlush(p.Hand[0], p.Hand[1], cards); v {
			r.Flush = append(r.Flush, p.ID)
			continue
		}

		if v := IsStraight(p.Hand[0], p.Hand[1], cards); v {
			r.Straight = append(r.Straight, p.ID)
			continue
		}

		
		if v := IsSet(p.Hand[0], p.Hand[1], cards); v {
			r.Set = append(r.Set, p.ID)
			continue
		}

		if v := IsTwoPair(p.Hand[0], p.Hand[1], cards); v {
			r.TwoPair = append(r.TwoPair, p.ID)
			continue
		}

		if v := IsPair(p.Hand[0], p.Hand[1], cards); v {
			r.Pair = append(r.Pair, p.ID)
			continue
		}

		r.HighestCard = append(r.HighestCard, p.ID)

	}


	if len(r.FlushRoyal) != 0 {
		return r.FlushRoyal
	} else if len(r.StraighFlush) != 0 {
		if len(r.StraighFlush) > 1 {
			return FindHigherFlush(cards, players)
		} else {
			return r.StraighFlush
		}
	} else if len(r.FourOfAKind) != 0 {
		if len(r.FourOfAKind) > 1 {
			return FindHigherFourOfAKind(cards, players)
		} else {
			return r.FourOfAKind
		}
	} else if len(r.FullHouse) != 0 {
		if len(r.FullHouse) > 1 {
			return FindHigherFullHouse(cards, players)
		} else {
			return r.FullHouse
		}
	} else if len(r.Flush) != 0 {
		if len(r.Flush) > 1 {

			return FindHigherFlush(cards, players)

		} else {
			return r.Flush
		}
	} else if len(r.Straight) != 0 {
		if len(r.Straight) > 1 {
		
			return FindHigherStraight(cards, players) 
		} else {
			return r.Straight
		}
	} else if len(r.Set) != 0 {
		if len(r.Set) > 1 {
			return FindHigherSet(cards, players)  
		} else {
			return r.Set
		}
	} else if len(r.TwoPair) != 0 {
		if len(r.TwoPair) > 1 {
			return FindHigherTwoPair(cards, players)  
		} else {
			return r.TwoPair
		}
	} else if len(r.Pair) != 0 {
		if len(r.Pair) > 1 {
			return FindHigherPair(cards, players)  
		} else {
			return r.Pair
		}
	} else {
		return FindHighestCard(cards, players) 
	}

	return res
}

func FindHighestCard(cards []string, players []PokerPlayer) []int {
	max := 0


	for _, p := range players {
		if CardValue(p.Hand[0]) == 1 || CardValue(p.Hand[1]) == 1 {
			max = 14
			continue
		}

		if CardValue(p.Hand[0]) >= max {
			max = CardValue(p.Hand[0])
		}

		if CardValue(p.Hand[1]) >= max {
			max = CardValue(p.Hand[1])
		}

	}


	res := []int{}

	for _, p := range players {
		card1 := CardValue(p.Hand[0]) 
		card2 := CardValue(p.Hand[1]) 

		if CardValue(p.Hand[0]) == 1 {
			card1 = 14
		}

		if CardValue(p.Hand[1]) == 1 {
			card2 = 14
		}

		if card1 == max || card2 == max {

			res = append(res, p.ID)

		}

	}

	return res
}

type WinnerCombinations struct {
	FlushRoyal []int 
	StraighFlush []int 
	FourOfAKind []int 
	FullHouse []int 
	Flush []int 
	Straight []int 
	Set []int 
	TwoPair []int 
	Pair []int 
	HighestCard []int
}

// type PokerPlayer struct {
// 	ID int `json:"id"`
// 	Username string `json:"username"`
// 	Chips int `json:"chips"`
// 	Action string `json:"action"`
// 	Status string `json:"status"`
// 	Img string `json:"img"`
// 	TimeRemains int `json:"time_remains"`
// 	TimeBank int `json:"time_bank"`
// 	Hand []string `json:"hand"`
// 	IsDealer bool `json:"is_dealer"`
// 	Bet int `json:"bet"`
// 	Combination string `json:"combination"`
// 	TotalBetsForRound int
// 	NextAction string `json:"nextAction"`
// }
