package tests

import (
	"testing"
	"github.com/Toront0/poker/internal/types/game"

)

func TestIsSet(t *testing.T) {

	card1 := "2hearts"
	card2 := "10spades"

	tableCards := []string{"10diamonds", "12spades", "13spades", "5spades", "10hearts"} 

	res := game.IsSet2(card1, card2, tableCards)

	if !res {

		t.Errorf("TestIsSet Failed. Got %t, want %t", res, true)

	}

}

func TestFindHigherSet(t *testing.T) {


	cards := []string{"2diamonds", "5clubs", "13clubs", "13hearts", "5diamonds"}

	players := []game.PokerPlayer{{8, "toronto", 10000, "", "waiting", "", 10, 10, []string{"1spades", "13clubs"}, true, 0, "", 0, ""}, {6, "admin", 50000, "", "waiting", "", 10, 10, []string{"12spades", "13clubs"}, false, 0, "", 0, ""}}

	res := game.FindHigherSet(cards, players)


	if len(res) > 0 {

		t.Errorf("TestFindHigherSet Failed. Got %v, want %v", res, []int{6})

	}

}