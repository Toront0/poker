package tests


import (
	"testing"
	"github.com/Toront0/poker/internal/card_combos"

)

func TestIsFourOfKind(t *testing.T) {
	card1 := "2hearts"
	card2 := "10spades"

	tableCards := []string{"10diamonds", "10spades", "10clubs", "5spades", "1spades"} 

	res := card_combos.IsFourOfAKind2(card1, card2, tableCards)


	if !res {
		t.Errorf("TestIsFourOfKind Failed. Got %t, want %t", res, true)
	}
}

func TestFindHigherFourOfAKind(t *testing.T) {

	cards := []string{"1clubs", "1clubs", "1clubs", "13hearts", "13clubs"}

	players := []card_combos.PokerPlayer{{8, "toronto", 10000, "", "waiting", "", 10, 10, []string{"1spades", "10clubs"}, true, 0, "", 0, ""}, {6, "admin", 50000, "", "waiting", "", 10, 10, []string{"13spades", "13clubs"}, false, 0, "", 0, ""}}

	res := card_combos.FindHigherFourOfAKind(cards, players)

	if len(res) > 0 {
		t.Errorf("TestFindHigherFourOfAKind Failed. Got %d, want %d", res, 1)
	}


	

}