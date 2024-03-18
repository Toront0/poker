package game

import (
	"testing"
	"github.com/Toront0/poker/internal/card_combos"

)

func TestIsFlushRoyal(t *testing.T) {
	
	card1 := "1spades"
	card2 := "10spades"

	tableCards := []string{"11spades", "12spades", "13spades", "5spades", "1spades"} 

	res := card_combos.IsFlushRoyal(card1, card2, tableCards)

	if res == "" {

		t.Errorf("Result was incorrect: got %s, want %s", res, "")

	}
}


func TestHighestCard(t *testing.T) {
	
	card1 := "13hearts"
	card2 := "1hearts"


	res := card_combos.HighestCard(card1, card2)

	if res != card1 {

		t.Errorf("Result was incorrect: got %s, want %s", res, card1)

	}
}





