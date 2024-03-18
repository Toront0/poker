package tests

import (
	"testing"
	"github.com/Toront0/poker/internal/types/game"

)

func TestFindHighestCard(t *testing.T) {

	cards := []string{"10clubs", "3clubs", "12clubs", "10hearts", "7clubs"}

	players := []card_combos.PokerPlayer{{8, "toronto", 10000, "", "waiting", "", 10, 10, []string{"1spades", "7clubs"}, true, 0, "", 0, ""}, {6, "admin", 50000, "", "waiting", "", 10, 10, []string{"2spades", "4clubs"}, false, 0, "", 0, ""}}

	res := game.FindHighestCard(cards, players)


	if res[0] != 8 {
		t.Errorf("TestFindHighestCard Failed. Got %v, want %v", res, []int{8})
	}

}