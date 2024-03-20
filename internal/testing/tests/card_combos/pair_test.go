package tests

import (
	"testing"
	"github.com/Toront0/poker/internal/types/game"

)

func TestIsPair(t *testing.T) {

	card1 := "2hearts"
	card2 := "10spades"

	tableCards := []string{"10diamonds", "12spades", "13spades", "5spades", "1spades"} 

	res := game.IsPair2(card1, card2, tableCards)


	if !res {
		t.Errorf("TestIsPair Failed. Got %t, want %t", res, true)
	}

}

func TestFindHigherPair(t *testing.T) {

	cards := []string{"10clubs", "3clubs", "13clubs", "10hearts", "1clubs"}

	players := []card_combos.PokerPlayer{{8, "toronto", 10000, "", "waiting", "", 10, 10, []string{"1spades", "7clubs"}, true, 0, "", 0, ""}, {6, "admin", 50000, "", "waiting", "", 10, 10, []string{"2spades", "1clubs"}, false, 0, "", 0, ""}}

	res := card_combos.FindHigherPair(cards, players)

	if len(res) > 0 {
		t.Errorf("TestFindHigherPair Failed. Got %d, want %d", res, 1)
	}

}