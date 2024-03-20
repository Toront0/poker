package tests

import (
	"testing"
	"github.com/Toront0/poker/internal/types/game"

)

func TestIsFullHouse(t *testing.T) {
	card1 := "2hearts"
	card2 := "10spades"

	tableCards := []string{"10diamonds", "10spades", "4clubs", "5spades", "2spades"} 

	res := game.IsFullHouse2(card1, card2, tableCards)


	if !res {
		t.Errorf("TestIsFullHouse Failed. Got %t, want %t", res, true)
	}
}

func TestFindHigherFullHouse(t *testing.T) {

	cards := []string{"7clubs", "6clubs", "8clubs", "5hearts", "8clubs"}

	players := []game.PokerPlayer{{8, "toronto", 10000, "", "waiting", "", 10, 10, []string{"6spades", "8clubs"}, true, 0, "", 0, ""}, {6, "admin", 50000, "", "waiting", "", 10, 10, []string{"8spades", "5clubs"}, false, 0, "", 0, ""}}

	res := game.FindHigherFullHouse(cards, players)

	if res[0] != 8 {
		t.Errorf("TestFindHigherFullHouse Failed. Got %d, want %d", res[0], 8)
	}

}