package game

import (
	"time"
	"fmt"
	"sync"
)






type PokerPlayer struct {
	ID int `json:"id"`
	Username string `json:"username"`
	Chips int `json:"chips"`
	Action string `json:"action"`
	Status string `json:"status"`
	Img string `json:"img"`
	TimeRemains int `json:"time_remains"`
	TimeBank int `json:"time_bank"`
	Hand []string `json:"hand"`
	IsDealer bool `json:"is_dealer"`
	Bet int `json:"bet"`
	Combination string `json:"combination"`
	TotalBetsForRound int
	NextAction string `json:"nextAction"`
}

func (p *PokerPlayer) SetNextAction(a string) {

	p.NextAction = a

}

func (p *PokerPlayer) Fold() {

	p.MakeAction("fold")

}

func (p *PokerPlayer) Check() {

	p.MakeAction("check")

}

func (p *PokerPlayer) Call(bet int) {

	p.MakeBet(bet)
	
	p.MakeAction("call")

}

func (p *PokerPlayer) Raise(bet int) {

	if bet - p.Chips == 0 {
		p.MakeAction("all-in")
	} else {
		p.MakeAction("raise")
	}

	p.MakeBet(bet)

}

func (p *PokerPlayer) InitTimer(RWMutex *sync.RWMutex) int {

	v := 20

	for i := v; i > 0; i-- {

		RWMutex.RLock()

		if p.Action != "" || p.Action == "all-in" || p.NextAction != "" {
			
			fmt.Println("ALL-IN")

			RWMutex.RUnlock()


			
			return v

		}

		RWMutex.RUnlock()
		
		RWMutex.Lock()

		p.TimeRemains--
		v--

		RWMutex.Unlock()

		fmt.Println("p.TimeRemains", p.TimeRemains, p.Action)

		time.Sleep(time.Second)
	}


	return v
}

func (p *PokerPlayer) InitTimeBank(RWMutex *sync.RWMutex) int {

	RWMutex.RLock()

	v := p.TimeBank

	RWMutex.RUnlock()

	for i := v; i > 0; i-- {
		
		RWMutex.RLock()
		if p.Action != "" || p.Action == "all-in" {
			
			RWMutex.RUnlock()
			return v

		}
		RWMutex.RUnlock()

		RWMutex.Lock()

		p.TimeBank--
		v--

		RWMutex.Unlock()

		fmt.Println("p.TimeBank", p.Bet)

		time.Sleep(time.Second)
	} 

	return v
}

func (p *PokerPlayer) MakeStatus(status string) {

	p.Status = status

}

func (p *PokerPlayer) MakeAction(action string) {

	p.Action = action

}

func (p *PokerPlayer) GetCard(card string) {

	p.Hand = append(p.Hand, card)

}


func (p *PokerPlayer) MakeDealer(v bool) {

	p.IsDealer = v

}

func (p *PokerPlayer) ChangeChipsTo(v int) {

	p.Chips += v

}

func (p *PokerPlayer) ChangeTotalBetsTo(v int) {
	p.TotalBetsForRound += v
}

func (p *PokerPlayer) MakeBet(bet int) {

	p.Bet += bet

	p.TotalBetsForRound += bet

	fmt.Println("BET IS ", bet)

	p.Chips -= bet
}