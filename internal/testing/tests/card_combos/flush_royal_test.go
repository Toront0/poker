package tests

import (
	"testing"
	"github.com/Toront0/poker/internal/types/game"
)

func TestIsFlushRoyal(t *testing.T) {
	
	card1 := "1spades"
	card2 := "9spades"

	tableCards := []string{"11spades", "12spades", "13spades", "5spades", "1clubs"} 

	res := game.IsFlushRoyal(card1, card2, tableCards)

	if !res {

		t.Errorf("Result was incorrect: got %t, want %t", res, true)

	}
}