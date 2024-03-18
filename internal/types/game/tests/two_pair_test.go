package tests

import (
	"testing"
	"github.com/Toront0/poker/internal/card_combos"

)

func TestIsTwoPair(t *testing.T) {

	card1 := "2hearts"
	card2 := "10spades"

	tableCards := []string{"10diamonds", "12spades", "13spades", "5spades", "2spades"} 

	res := card_combos.IsTwoPair2(card1, card2, tableCards)

	if !res {

		t.Errorf("TestIsTwoPair Failed. Got %t, want %t", res, true)

	}

}

func TestFindHigherTwoPair(t *testing.T) {
	cards := []string{"10clubs", "3clubs", "13clubs", "10hearts", "1clubs"}

	players := []card_combos.PokerPlayer{{8, "toronto", 10000, "", "waiting", "", 10, 10, []string{"12spades", "7clubs"}, true, 0, "", 0, ""}, {6, "admin", 50000, "", "waiting", "", 10, 10, []string{"2spades", "13clubs"}, false, 0, "", 0, ""}}

	res := card_combos.FindHigherTwoPair(cards, players)

	if res[0] != 6 {
		t.Errorf("TestFindHigherTwoPair Failed. Got %v, want %v", res, []int{6})
	}


}