package game

import (
	"time"
	"github.com/Toront0/poker/internal/utils"
	"fmt"
	"math/rand"
	"sync"
	"slices"
	"cmp"
)

type PokerTable struct {
	ID int `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Name string `json:"name"`
	Cards []string `json:"-"`
	Players []PokerPlayer `json:"players"`
	Winner []int `json:"winner"`
	Ante int `json:"ante"`
	AnteBetsAmount int `json:"-"`
	MinBet int `json:"min_bet"`
	Pot int `json:"pot"`
	Flop []string `json:"flop"`
	Turn string `json:"turn"`
	River string `json:"river"`
	IsAllInStage bool `json:"isAllInStage"`
	PotentialStreetCards []string `json:"potentialStreetCards"`
	PrizeDestribution string `json:"prizeDestribution"`
	Mutex sync.RWMutex `json:"-"`
	Prize int `json:"prize"`
	AmountOfPlayers int `json:"amountOfPlayers"`
	Mode string `json:"mode"`
	BuyIn int `json:"buyIn"`
}

func (t *PokerTable) NewDeck() {
	deck := []string{}


	for i := 1; i < 54; i++ {

		if (i <= 13) {
			card := fmt.Sprintf("%d%s", i, "spades")

			deck = append(deck, card)
		}

		if (i > 13 && i <= 26) {
			card := fmt.Sprintf("%d%s", i - 13, "hearts")

			deck = append(deck, card)
		}

		if (i > 26 && i <= 39) {
			card := fmt.Sprintf("%d%s", i - 26, "diamonds")

			deck = append(deck, card)
		}

		if (i > 40 && i <= 53) {
			card := fmt.Sprintf("%d%s", i - 40, "clubs")

			deck = append(deck, card)
		}
	}

	t.Cards = deck

	rand.Shuffle(len(t.Cards), func(i, j int) {
		t.Cards[i], t.Cards[j] = t.Cards[j], t.Cards[i]
	})
}

func (t *PokerTable) ShouldRevealFlop() bool {

	defer t.Mutex.RUnlock()

	t.Mutex.RLock()

	if len(t.Flop) > 2 {
		return false
	}

	for _, p := range t.Players {



		if p.Action == "" {
			return false
		}

	}

	return true
}

func (t *PokerTable) RevealNextTableCard() {

	t.Mutex.Lock()
	defer t.Mutex.Unlock()



	if len(t.Flop) == 0 {

		t.PotentialStreetCards = t.Cards[:3]

		return

	}

	if t.Turn == "" {

		t.PotentialStreetCards = append(t.PotentialStreetCards, t.Cards[0])

		return
	}

	if t.River == "" {

		t.PotentialStreetCards = append(t.PotentialStreetCards, t.Cards[0])

		return
	}

	

}

func (t *PokerTable) ShouldRevealNextTableCard() {
	res := true

	defer t.Mutex.Unlock()

	t.Mutex.Lock()
	
	for _, p := range t.Players {

		if p.Action == "" {
			res = false
			return 
		}

		if p.Action == "raise" {
			res = false
			return 
		}

	}


	t.PrepareToCardDeal()

	


	if len(t.Flop) == 0 && res {

		t.RevealFlop()
		return
	}

	if t.Turn == "" && res {
		t.RevealTurn()
		return
	}

	if t.River == "" && res {
		t.RevealRiver()
		return
	}


}

func (t *PokerTable) ShouldRevealTurn() bool {

	defer t.Mutex.RUnlock()
	t.Mutex.RLock()

	if len(t.Flop) < 2 {
		return false
	}

	if t.Turn != "" {
		return false
	}

	for _, p := range t.Players {

		if p.Action == "" {
			return false
		}

		if p.Action == "raise" {
			return false
		}

	}

	return true
}

func (t *PokerTable) ShouldRevealRiver() bool {

	defer t.Mutex.RUnlock()
	t.Mutex.RLock()

	if len(t.Flop) < 2 {
		return false
	}

	if t.Turn == "" {
		return false
	}

	fmt.Println("REVEALING RIVER")


	for _, p := range t.Players {

		if p.Action == "" || p.Action == "raise" {
			return false
		}

	}

	return true
}



func (t *PokerTable) RevealFlop() {
	
	t.Flop = t.Cards[:3]


	

	t.Cards = t.Cards[3:]

	

}

func (t *PokerTable) RevealTurn() {
	fmt.Println("cards before turn", t.Cards)

	t.Turn = t.Cards[0]

	fmt.Println("cards after turn", t.Cards)

	t.Cards = t.Cards[1:]

}

func (t *PokerTable) RevealRiver() {

	t.River = t.Cards[0]

	t.Cards = t.Cards[1:]
}

func (t *PokerTable) RequiredPaysAfterDealer() {


	dealerIdx := utils.SliceIndex(len(t.Players), func(i int) bool { return t.Players[i].IsDealer == true })

	if len(t.Players) == 2 {
		madeBetsAmount := 0

		minBet := utils.Min(t.Players[dealerIdx].Chips, t.MinBet / 2)

		// t.Players[dealerIdx].MakeBet(utils.Min(t.Players[dealerIdx].Chips ,t.MinBet / 2))
		t.Players[dealerIdx].MakeBet(minBet)

		madeBetsAmount += minBet

		if dealerIdx + 1 >= len(t.Players) {

			minBet := utils.Min(t.Players[0].Chips, t.MinBet)

			t.Players[0].MakeBet(minBet)

			madeBetsAmount += minBet

		} else {
			minBet := utils.Min(t.Players[dealerIdx + 1].Chips, t.MinBet)

			// t.Players[dealerIdx + 1].MakeBet(utils.Min(t.Players[dealerIdx + 1].Chips, t.MinBet))
			t.Players[dealerIdx + 1].MakeBet(minBet)

			madeBetsAmount += minBet
		}
		t.AddChipsToPot(madeBetsAmount)
		return
	}

	if dealerIdx + 1 >= len(t.Players) {

		smallBlindBet := utils.Min(t.Players[0].Chips, t.MinBet / 2)
		bigBlindBet := utils.Min(t.Players[1].Chips, t.MinBet)

		t.Players[0].MakeBet(smallBlindBet)
		t.Players[1].MakeBet(bigBlindBet)
		t.AddChipsToPot(smallBlindBet + bigBlindBet)
		return
	}
	

	minBet := utils.Min(t.Players[dealerIdx + 1].Chips, t.MinBet / 2)
	t.Players[dealerIdx + 1].MakeBet(minBet)
	t.AddChipsToPot(minBet)

	
	
	if dealerIdx + 2 >= len(t.Players) {
		minBet := utils.Min(t.Players[0].Chips, t.MinBet)

		t.Players[0].MakeBet(minBet)
		t.AddChipsToPot(minBet)
		
		return
	}

	minAvailableBet := utils.Min(t.Players[dealerIdx + 2].Chips, t.MinBet)

	t.Players[dealerIdx + 2].MakeBet(minAvailableBet)
	t.AddChipsToPot(minAvailableBet)
}

func (t *PokerTable) AddChipsToPot(a int) {

	t.Pot += a


}

func (t *PokerTable) DestributeCards() {

	for i := 0; i < len(t.Players); i++ {
		if i == 0 {
			t.Players[i].GetCard(t.Cards[0])
			t.Players[i].GetCard(t.Cards[1])
		
		} else {
			
			t.Players[i].GetCard(t.Cards[i + i])
			t.Players[i].GetCard(t.Cards[i + i + 1])
		}

	}

	t.Cards = t.Cards[len(t.Players) * 2:]

}

func (t *PokerTable) ChangeDealer()  {

	dealerIdx := utils.SliceIndex(len(t.Players), func(i int) bool { return t.Players[i].IsDealer == true })

	

	if dealerIdx == -1 {
		t.Players[0].MakeDealer(true)
		return
	}

	t.Players[dealerIdx].MakeDealer(false)

	if dealerIdx + 1 >= len(t.Players) {

		t.Players[0].MakeDealer(true)

	} else {
		t.Players[dealerIdx + 1].MakeDealer(true)
	}

	

}

func (t *PokerTable) FindActivePlayer() int {
	
	dealerIdx := utils.SliceIndex(len(t.Players), func(i int) bool { return t.Players[i].IsDealer == true })
	
	fmt.Println("dealerIdx", dealerIdx)


	if dealerIdx + 3 >= len(t.Players) {

		if len(t.Players) == 2 {
			t.Players[dealerIdx + 2 - len(t.Players)].MakeStatus("active")
			return dealerIdx + 2 - len(t.Players)
		}

		t.Players[dealerIdx + 3 - len(t.Players)].MakeStatus("active")
		return dealerIdx + 3 - len(t.Players)


	}

	t.Players[dealerIdx + 3].MakeStatus("active")


	return dealerIdx + 3
}

func (t *PokerTable) FindActivePlayerAfterCardDeal() int {


	dealerIdx := utils.SliceIndex(len(t.Players), func(i int) bool { return t.Players[i].IsDealer == true })

	nextP := -1

	j := dealerIdx


	for i := 0; i < len(t.Players); i++ {

		if j + 1 >= len(t.Players) {
			j = 0
		} else {
			j++
		}

		if t.Players[j].Action == "fold" {
			continue
		} else {
			nextP = j
			t.Players[j].MakeStatus("active")
			return j
		}

	}

	return nextP
}

func (t *PokerTable) AllFoldWinnerId() int {

	defer t.Mutex.RUnlock()
	t.Mutex.RLock()

	playersAtTable := utils.Filter(t.Players, func (p PokerPlayer) bool { return p.Action != "fold" })

	winnerIdx := -1

	fmt.Println("playersAtTable 123", len(playersAtTable))

	if len(playersAtTable) == 1 {
		

		return playersAtTable[0].ID
	}

	return winnerIdx
}

func (t *PokerTable) AllowSystemToCheck(playerIdx int) bool {


	for _, p := range t.Players {

		if p.Action == "raise" || p.Action == "all-in" {

			fmt.Println("p.Action", p.Action)

			return false

		} 


	}

	fmt.Println("after all-n")

	maxBet := 0

	for _, p := range t.Players {

		if p.Bet > maxBet {
			maxBet = p.Bet
		}

	}


	if t.Players[playerIdx].Action == "" && t.Players[playerIdx].Bet == maxBet {
		fmt.Println("maxBet", t.Players[playerIdx].Bet )
		return true

	}

	return false
}

func (t *PokerTable) DetermineWinner(w int) {

	t.Mutex.Lock()

	t.Winner = append(t.Winner, w)

	t.Mutex.Unlock()
}


func (t *PokerTable) FindWinner() {

	defer t.Mutex.Unlock()
	t.Mutex.Lock()

	cards := []string{t.Flop[0], t.Flop[1], t.Flop[2], t.Turn, t.River}


	pOnTable := utils.Filter(t.Players, func(p PokerPlayer) bool { return p.Action != "fold" })

	fmt.Println("pOnTable", pOnTable)

	ws := DetermineWhosHandHigher(cards, pOnTable)

	t.Winner = ws

	fmt.Println("winners", ws)


}

func (t *PokerTable) ShouldFindWinner() bool {
	t.Mutex.RLock()
	defer t.Mutex.RUnlock()


	if len(t.Flop) < 2 {
		return false
	} 

	if t.Turn == "" {
		return false
	}

	if t.River == "" {
		return false
	}

	for _, p := range t.Players {

		if p.Action == "" {

			return false

		}

	} 

	return true	

}

// func (t *PokerTable) GetWinnersAndLosers() []int {
// 	res := t.KickLoser()

// 	return t.Winner, res
// }

func (t *PokerTable) GetWinner() []int {


	return t.Winner
}



func (t *PokerTable) ReturnRestChipsToPlayers(rest int) {


	for _, p := range t.Players {
		if p.Action == "fold" {
			continue
		}

		if utils.Contains(t.Winner, p.ID) {

			continue

		}
		

		


		if p.Action == "fold" {
			continue
		}

		p.ChangeChipsTo(rest / len(t.Winner))

	}



}

func (t *PokerTable) DestributeChipsAfterRound() {

	defer t.Mutex.Unlock()
	t.Mutex.Lock()

	playersAtTable := 0


	for _, p := range t.Players {

		if p.Action != "fold" {
			playersAtTable++
		}

	}

	rest := t.Pot


	minBetPlayer := slices.MinFunc(t.Players, func (a, b PokerPlayer) int {
		return cmp.Compare(a.TotalBetsForRound, b.TotalBetsForRound)
	})

	// rest -= minBetPlayer.TotalBetsForRound

	fmt.Println("minBetPlayer.TotalBetsForRound", minBetPlayer.TotalBetsForRound)
	noBets := t.AnteBetsAmount + t.MinBet + t.MinBet / 2 == t.Pot

	for _, w := range t.Winner {

		winnerIdx := utils.SliceIndex(len(t.Players), func (i int) bool { return t.Players[i].ID == w })

		fmt.Println("t.Players[winnerIdx].TotalBetsForRound", t.Players[winnerIdx].TotalBetsForRound)

		playerTotalBets := t.Players[winnerIdx].TotalBetsForRound

		if len(t.Winner) == 1 && noBets  {
			fmt.Println("Inside")
			t.Players[winnerIdx].ChangeChipsTo(t.Pot)
			rest = 0
			break
		}

		prize := minBetPlayer.TotalBetsForRound * (playersAtTable - len(t.Winner) + 1)

		rest -= prize


		if prize < playerTotalBets && rest >= playerTotalBets - minBetPlayer.TotalBetsForRound {

			

			t.Players[winnerIdx].ChangeChipsTo(playerTotalBets - minBetPlayer.TotalBetsForRound)	
			rest -= t.Players[winnerIdx].TotalBetsForRound - minBetPlayer.TotalBetsForRound
		}

		t.Players[winnerIdx].ChangeChipsTo(prize)
	}

	if len(t.Winner) > 1 {
		for _, w :=  range t.Winner {

			winnerIdx := utils.SliceIndex(len(t.Players), func (i int) bool { return t.Players[i].ID == w })
	
			fmt.Println("t.Players[winnerIdx].TotalBetsForRound", t.Players[winnerIdx].TotalBetsForRound)
	
			if t.Players[winnerIdx].TotalBetsForRound == minBetPlayer.TotalBetsForRound {
				continue
			}
	
			maxPrize := t.Players[winnerIdx].TotalBetsForRound * (playersAtTable - len(t.Winner) + 1) - minBetPlayer.TotalBetsForRound * (playersAtTable - len(t.Winner) + 1)
	
	
	
			fmt.Println("maxPrize", maxPrize)
			fmt.Println("rest", rest)
	
			if rest - maxPrize < 0 {
				t.Players[winnerIdx].ChangeChipsTo(rest)
				rest = 0
			} else {
				t.Players[winnerIdx].ChangeChipsTo(maxPrize)
				rest -= maxPrize
			}

		} 
	}


	if rest > 0 && !noBets {
		
		for i, p := range t.Players {
			if p.Action == "fold" {
				continue
			}

			if utils.Contains(t.Winner, p.ID) {
				continue
			}
		
			t.Players[i].ChangeChipsTo(rest / (playersAtTable - len(t.Winner)))
			// rest = rest / (playersAtTable - len(t.Winner))
			
		} 
	}

	return 
}


func (t *PokerTable) TakeAntePays() {

	for i := 0; i < len(t.Players); i++ {
		takenAnte := utils.Min(t.Players[i].Chips, t.Ante)

		if takenAnte >= t.Players[i].Chips {
			t.Players[i].MakeBet(takenAnte)
			t.AddChipsToPot(takenAnte)
			continue
		}

		t.Players[i].ChangeChipsTo(-takenAnte)
		t.Players[i].ChangeTotalBetsTo(takenAnte)
		t.AnteBetsAmount += takenAnte
		t.AddChipsToPot(takenAnte)
	}

	
	
}

func (t *PokerTable) IncreaseAntePay() {


	t.Ante = int(float64(t.MinBet) * 0.2)

}

func (t *PokerTable) IncreaseMinBet() {


	t.MinBet *= 2

}

func (t *PokerTable) GetWinnerDetailForPrize() PlayerPrizeInfo {

	prize := 0

	if t.PrizeDestribution == "winner-takes-all" {
		prize = t.Prize
	}

	return PlayerPrizeInfo{
		ID: t.Players[0].ID,
		Prize: prize,
		Place: 1,
	}

}

func (t *PokerTable) PrepareToCardDestribution() bool {

	// res := t.KickLoser()

	if len(t.Players) == 1 {
		return false
	}

	t.Cards = nil
	t.Winner = []int{}
	t.Pot = 0
	t.Flop = []string{}
	t.Turn = ""
	t.River = ""
	t.PotentialStreetCards = []string{}
	t.IsAllInStage = false
	t.AnteBetsAmount = 0
	
	for i := 0; i < len(t.Players); i++ {
		t.Players[i].MakeAction("")

		t.Players[i].MakeStatus("waiting")

		t.Players[i].Hand = []string{}

		t.Players[i].TimeRemains = 20

		t.Players[i].Bet = 0

		t.Players[i].TotalBetsForRound = 0
	}

	return true

}
func (t *PokerTable) ResetPlayersActionAfterRaise(raiserID int) {
	


	for i, p := range t.Players {
		if p.NextAction == "call" || p.NextAction == "check" || p.NextAction == "raise"  {

			t.Players[i].NextAction = ""

		}

		if p.ID == raiserID || p.Chips == 0 {
			continue
		}

		if p.Action == "call" || p.Action == "raise" || p.Action == "check" {
			t.Players[i].Action = ""
		}


	}

}

func (t *PokerTable) PrepareToCardDeal() {
	// defer t.Mutex.RUnlock()
	// t.Mutex.RLock()

	a := t.IsAllInStage


	for i, p := range t.Players {
		

		t.Players[i].Bet = 0

		if a {
			continue
		}

		if p.Action == "all-in" || p.Chips == 0 {
			continue
		}

		if p.Action == "call" || p.Action == "check" {
		

			t.Players[i].MakeAction("")
			t.Players[i].TimeRemains = 20
			t.Players[i].SetNextAction("")

		}

	}

}

type PlayerPrizeInfo struct {
	ID int
	Place int 
	Prize int
}

func (t *PokerTable) KickLoser() []PlayerPrizeInfo {
	kicked := []PlayerPrizeInfo{}
	t.Mutex.Lock()
	defer t.Mutex.Unlock()

	for _, p := range t.Players {

		if p.Chips == 0 {
			prize := 0

			if t.PrizeDestribution == "winner-takes-all" && len(t.Players) == 1 {
				prize = t.Prize
			}

			p := PlayerPrizeInfo{
				ID: p.ID,
				Prize: prize,
				Place: len(t.Players),
			}

			kicked = append(kicked, p)

		}


	}

	t.Players = utils.Filter(t.Players, func(p PokerPlayer) bool { return p.Chips > 0  })


	return kicked
}

func (t *PokerTable) IsGameOnAllInStage() bool {
	activeQty := 0
	allInQty := 0

	
	defer t.Mutex.Unlock()
	t.Mutex.Lock()

	for _, p := range t.Players {

		if p.Action == "" {
			return false
		}

		if p.Action == "all-in" || p.Chips == 0 {
			allInQty++
			continue
		}



		if p.Action != "fold" {
			activeQty++
			continue
		}
	}

	fmt.Println("allInQty", allInQty)
	fmt.Println("activeQty", activeQty)

	if activeQty >= 2 {
		return false
	}

	// if activeQty < 1 || activeQty >= 2 {
	// 	return false
	// }

	if allInQty - activeQty >= 0  {
		t.IsAllInStage = true

		return true
		
	}


	return false

}

func (t *PokerTable) DetermineNextActivePlayer(currIdx int) int {

	nextP := -1


	t.Players[currIdx].MakeStatus("waiting")
	// t.Players[currIdx].MakeAction("fold")

	j := currIdx


	for i := 0; i < len(t.Players); i++ {

		if j + 1 >= len(t.Players) {
			j = 0
		} else {
			j++
		}

		
		if t.Players[j].Action == "check" {
			
			t.Players[j].MakeStatus("active")
			nextP = j
			break
		}

		if t.Players[j].Action == "" {
			
			t.Players[j].MakeStatus("active")
			nextP = j
			break
		}



		if t.Players[j].Action == "fold" || t.Players[j].Action == "all-in" {
			continue
		}

	}

	return nextP
}


func (t *PokerTable) ShuffleCards() {

	rand.Shuffle(len(t.Cards), func(i, j int) {
		t.Cards[i], t.Cards[j] = t.Cards[j], t.Cards[i]
	})

}


func (t *PokerTable) ResetRaisersActionAfterCall() {
	

	for i, p := range t.Players {
		if p.Action == "raise" {


			t.Players[i].Action = "call"

			// p.MakeAction("call")
		}

	}
}


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





func (t PokerTable) SendData(uID int) PokerTable {

	// if a := t.IsGameOnAllInStage(); a == falset. {
	// 	return HideCard(t, uID)
	// } else {
	// 	return t
	// }

	fmt.Println("before send ", t.IsAllInStage)	

	if !t.IsAllInStage {
		return HideCard(t, uID)
	} else {
		return t
	}
}

func (t *PokerTable) RevealPlayerCards() []PokerPlayer {


	playersAtTable := utils.Filter(t.Players, func (p PokerPlayer) bool { return p.Action != "fold" })


	return playersAtTable
}

func HideCard(t PokerTable, eID int) PokerTable {
	res := t

	ps := []PokerPlayer{}

	for _, p := range t.Players {

		// if p.Action == "fold" {
		// 			newP := PokerPlayer{}

		// 	newP.ID = p.ID
		// 	newP.Username = p.Username
		// 	newP.Chips = p.Chips
		// 	newP.Action = p.Action
		// 	newP.Status = p.Status
		// 	newP.Img = p.Img
		// 	newP.TimeRemains = p.TimeRemains
		// 	newP.TimeBank = p.TimeBank
		// 	newP.Hand = p.Hand
		// 	newP.Bet = p.Bet
		// 	newP.Combination = p.Combination
		// 	newP.IsDealer = p.IsDealer
		// 	newP.NextAction = ""
		// 	ps = append(ps, newP)

		// 	continue
		// }
		
		if eID == p.ID {
			ps = append(ps, p)
		} else {
			newP := PokerPlayer{}

			newP.ID = p.ID
			newP.Username = p.Username
			newP.Chips = p.Chips
			newP.Action = p.Action
			newP.Status = p.Status
			newP.Img = p.Img
			newP.TimeRemains = p.TimeRemains
			newP.TimeBank = p.TimeBank
			newP.Hand = []string{"", ""}
			newP.Bet = p.Bet
			newP.Combination = p.Combination
			newP.IsDealer = p.IsDealer
			newP.NextAction = ""
			ps = append(ps, newP)
			
		}

	}

	res.Players = ps

	return res
}