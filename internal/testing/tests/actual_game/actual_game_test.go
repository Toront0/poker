package tests

import (
	"testing"
	"github.com/Toront0/poker/internal/types/game"
)



func TestDestributeChipsAfterRound(t *testing.T) {

	g := game.PokerTable{
		Pot: 60000,
	}

	g.Players = []game.PokerPlayer{{8, "toronto", 0, "", "waiting", "", 10, 10, []string{}, true, 0, "", 10000, ""}, {6, "admin", 0, "", "waiting", "", 10, 10, []string{}, false, 0, "", 50000, ""}}

	g.Winner = []int{6}

	g.DestributeChipsAfterRound()

	if g.Players[1].Chips != 60000 {
		t.Errorf("TestDestributeChipsAfterRound Failed. Got %d, want %d", g.Players[1].Chips, 60000)
	}

}