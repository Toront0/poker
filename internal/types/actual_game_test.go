package types

import (
	"testing"
	"fmt"
	"github.com/Toront0/poker/internal/types"

)

func TestDetermineWhosHandHigher(t *testing.T) {

	cards := []string{"10spades", "7spades", "12spades", "10hearts", "4spades"}

	players := []types.PokerPlayer{{8, "toronto", 10000, "", "waiting", "", 10, 10, []string{"2spades", "3spades"}, true, 0, "", 0, ""}, {6, "admin", 50000, "", "waiting", "", 10, 10, []string{"2spades", "3clubs"}, false, 0, "", 0, ""}}


	res := types.DetermineWhosHandHigher(cards, players)

	fmt.Println("res", res)

	if len(res) != 2 {
		t.Errorf("Result incorrect got: %v, want: %v", res, []int{5})
	}


}
func TestHighestCard(t *testing.T) {


	players := []types.PokerPlayer{{8, "toronto", 10000, "", "waiting", "", 10, 10, []string{"2spades", "3spades"}, true, 0, "", 0, ""}, {6, "admin", 50000, "", "waiting", "", 10, 10, []string{"1clubs", "4clubs"}, false, 0, "", 0, ""}}


	res := types.HighestCard(players, []int{8, 6})

	

	if len(res) != 0 {
		t.Errorf("Result incorrect got: %v, want: %v", res, []int{5})
	}


}