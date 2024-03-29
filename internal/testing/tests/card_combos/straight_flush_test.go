package tests

import (
	"testing"
	"github.com/Toront0/poker/internal/types/game"

)

func TestIsStraightFlush(t *testing.T) {

	card1 := "2hearts"
	card2 := "7clubs"

	tableCards := []string{"3spades", "4spades", "5spades", "6spades", "10hearts"} 

	res := game.IsStraightFlush(card1, card2, tableCards)

	if !res {

		t.Errorf("TestIsStraightFlush Failed. Got %t, want %t", res, true)

	}

}