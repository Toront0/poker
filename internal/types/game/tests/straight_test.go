package tests

import (
	"testing"
	"github.com/Toront0/poker/internal/types/game"

)

func TestIsStraight(t *testing.T) {


	card1 := "7hearts"
	card2 := "10spades"

	tableCards := []string{"8spades", "6spades", "13spades", "1spades", "9clubs"} 

	res := game.IsStraight(card1, card2, tableCards)

	if !res {

		t.Errorf("TestIsStraight Failed. Got %v, want %v", res, true)

	}

}

func TestIsStraightWithAce(t *testing.T) {


	card1 := "5hearts"
	card2 := "10spades"

	tableCards := []string{"11spades", "12spades", "13spades", "5spades", "1clubs"} 

	res := game.IsStraight(card1, card2, tableCards)

	if !res {

		t.Errorf("TestIsStraightWithAce Failed. Got %v, want %v", res, true)

	}

}

func TestFindHigherStraight(t *testing.T) {

	cards := []string{"3clubs", "10clubs", "11clubs", "12hearts", "9clubs"}

	players := []game.PokerPlayer{{8, "toronto", 10000, "", "waiting", "", 10, 10, []string{"13spades", "9clubs"}, true, 0, "", 0, ""}, {6, "admin", 50000, "", "waiting", "", 10, 10, []string{"9spades", "8clubs"}, false, 0, "", 0, ""}}

	res := game.FindHigherStraight(cards, players)

	if res[0] != 8  {
		t.Errorf("TestFindHigherStraight Test failed - got: %v, want: %d", res, 8)
	}
}

func TestFindHigherStraightWithAce(t *testing.T) {

	cards := []string{"3clubs", "11clubs", "12clubs", "13hearts", "9clubs"}

	players := []game.PokerPlayer{{8, "toronto", 10000, "", "waiting", "", 10, 10, []string{"10spades", "1clubs"}, true, 0, "", 0, ""}, {6, "admin", 50000, "", "waiting", "", 10, 10, []string{"9spades", "10clubs"}, false, 0, "", 0, ""}}

	res := game.FindHigherStraight(cards, players)

	if len(res) > 0  {
		t.Errorf("TestFindHigherStraightWithAce Test failed - got: %v, want: %d", res, 8)
	}
}