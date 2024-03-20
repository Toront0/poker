package tests

import (
	"testing"
	"github.com/Toront0/poker/internal/types/game"

)

func TestIsFlush(t *testing.T) {
	
	card1 := "1hearts"
	card2 := "10spades"

	tableCards := []string{"11diamonds", "12spades", "13spades", "5spades", "1spades"} 

	res := game.IsFlush(card1, card2, tableCards)

	if !res {

		t.Errorf("Result was incorrect: got %t, want %t", res, true)

	}
}

func TestFindHigherFlush(t *testing.T) {
	cards := []string{"10clubs", "3clubs", "12clubs", "10hearts", "7clubs"}

	players := []game.PokerPlayer{{8, "toronto", 10000, "", "waiting", "", 10, 10, []string{"2spades", "7clubs"}, true, 0, "", 0, ""}, {6, "admin", 50000, "", "waiting", "", 10, 10, []string{"2spades", "4clubs"}, false, 0, "", 0, ""}}


	res := game.FindHigherFlush(cards, players)

	if res[0] != 8  {
		t.Errorf("TestFindHigherFlush Test failed - got: %v, want: %d", res, 8)
	}
}


func TestFindHigherFlushWithAce(t *testing.T) {
	cards := []string{"10clubs", "3clubs", "12clubs", "10hearts", "7clubs"}

	players := []game.PokerPlayer{{8, "toronto", 10000, "", "waiting", "", 10, 10, []string{"2spades", "7clubs"}, true, 0, "", 0, ""}, {6, "admin", 50000, "", "waiting", "", 10, 10, []string{"2spades", "1clubs"}, false, 0, "", 0, ""}}


	res := game.FindHigherFlush(cards, players)

	if res[0] != 6  {
		t.Errorf("TestFindHigherFlush Test failed - got: %v, want: %d", res, 1)
	}
}