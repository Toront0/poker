package actualGame

import (
	"github.com/go-chi/chi/v5"

	"github.com/Toront0/poker/internal/types/game"
	"github.com/Toront0/poker/internal/types"
	"github.com/Toront0/poker/internal/utils"
	"github.com/Toront0/poker/internal/services"
	"time"
	"fmt"
	"sync"
	"encoding/json"

)

type PokerGame struct {
	GameData *game.PokerTable
	hub *Hub
	mux *chi.Mux
	store services.ActualGameStorer
}

func NewPokerGame(mux *chi.Mux, store services.ActualGameStorer, data types.GameDetail) *PokerGame {


	// players := []game.PokerPlayer{{8, "toronto", 10000, "", "waiting", "", 10, 10, []string{}, true, 0, "", 0, ""}, {6, "admin", 50000, "", "waiting", "", 10, 10, []string{}, false, 0, "", 0, ""}}
	// players := []game.PokerPlayer{{7, "anton", 50000, "", "waiting", "", 10, 10, []string{}, true, 0, "", 0, ""}, {6, "admin", 50000, "", "waiting", "", 10, 10, []string{}, false, 0, "", 0, ""}, {8, "toronto", 10000, "", "waiting", "", 10, 10, []string{}, false, 0, "", 0, ""},{10, "ivan", 10000, "", "waiting", "", 10, 10, []string{}, false, 0, "", 0, ""}}
	players := []game.PokerPlayer{}

	for i, p := range data.Players {


		player := game.PokerPlayer{}

		player.ID = p.ID
		player.Username = p.Username
		player.Img = p.Img
		player.Chips = 50000
		player.Status = "waiting"
		player.TimeBank = 50
		player.Hand = []string{"", ""}
		player.TimeRemains = 20
		
		if i == 0 {
			player.Status = "active"
		}

		players = append(players, player)



	}


	var mutex sync.RWMutex

	d := &game.PokerTable{
		ID: data.ID,
		CreatedAt: data.CreatedAt,
		Name: data.Name,
		BuyIn: data.BuyIn,
		AmountOfPlayers: data.AmountOfPlayers,
		Prize: data.Prize,
		Mode: data.Mode,
		PrizeDestribution: data.PrizeDestribution,
		MinBet: 100,
		Ante: 20,
		Players: players,
		Mutex: mutex,
		Flop: []string{},
		Turn: "",
		River: "",
	}

	hub := NewHub()
	go hub.Run()

	return &PokerGame{
		GameData: d,
		hub: hub,
		mux: mux,
		store: store,
	}
}


func (g *PokerGame) SetupRoutes() {

	g.mux.Get(fmt.Sprintf("/poker/table/%d/{id}", g.GameData.ID), g.GetGameData)
	g.mux.Get(fmt.Sprintf("/poker/fold/%d/{id}", g.GameData.ID), g.Fold)
	g.mux.Get(fmt.Sprintf("/poker/check/%d/{id}", g.GameData.ID), g.Check)
	g.mux.Post(fmt.Sprintf("/poker/call/%d", g.GameData.ID), g.Call)
	g.mux.Post(fmt.Sprintf("/poker/raise/%d", g.GameData.ID), g.Raise)
	g.mux.Post(fmt.Sprintf("/poker/emoji/%d", g.GameData.ID), g.Raise)
	g.mux.HandleFunc(fmt.Sprintf("/ws/poker/table/%d/{id}", g.GameData.ID), g.ServeWs)

}

func (g *PokerGame) Run() {
	cardDest := make(chan bool, 1)
	winnerCh := make(chan bool, 1)
	flopCh := make(chan bool, 1)
	// turnCh := make(chan bool, 1)
	// riverCh := make(chan bool, 1)
	gameEnds := make(chan bool, 1)


	cardDest <- true

	go g.InitAnteMinBetHandler(gameEnds)

	for {

		select {
		case <- cardDest:
			
			g.GameData.Mutex.Lock()

			a := g.GameData.PrepareToCardDestribution()


			if !a {
				winnerDetail := g.GameData.GetWinnerDetailForPrize()

				err := g.store.GivePlayersPrize(winnerDetail.ID, winnerDetail.Prize, winnerDetail.Place, g.GameData.ID)

				if err != nil {
					fmt.Printf("could not give players's prize %s", err)
				}


				g.store.EndTheGame(g.GameData.ID)

				gameEnds <- true

				g.hub.Broadcast <- []byte("game-end")
				
				return
			}
			

			g.GameData.ChangeDealer()

			g.GameData.TakeAntePays()

			g.GameData.RequiredPaysAfterDealer()

			g.GameData.NewDeck()

			g.GameData.DestributeCards()
			
			

			idx := g.GameData.FindActivePlayer()

			g.GameData.Mutex.Unlock()

			g.hub.Broadcast <- []byte("global-changes")

			t := g.GameData.Players[idx].InitTimer(&g.GameData.Mutex)

			fmt.Println("time remains", t)

			if t == 0 {
				g.GameData.Mutex.RLock()
				if (g.GameData.Players[idx].TimeBank > 0) {

					g.GameData.Mutex.RUnlock()

					t = g.GameData.Players[idx].InitTimeBank(&g.GameData.Mutex)

					if t == 0 {
						g.GameData.Mutex.Lock()

						g.GameData.Players[idx].MakeAction("fold")

						g.GameData.Mutex.Unlock()
					}

				} else {
					g.GameData.Mutex.RUnlock()

					g.GameData.Mutex.Lock()

					g.GameData.Players[idx].MakeAction("fold")

					g.GameData.Mutex.Unlock()
				}

				
			}

			
		case <- winnerCh:
			g.GameData.DestributeChipsAfterRound()

			losers := g.GameData.KickLoser()

			g.GameData.RevealNextTableCard()

			res := &WSWinnerResponse{
				Action: "winner",
				Winner: g.GameData.GetWinner(),
				PotentialStreetCards: g.GameData.PotentialStreetCards,
				KickedPlayers: []int{},
			}
			
			if len(losers) > 0 {

				for _, p := range losers {
					res.KickedPlayers = append(res.KickedPlayers, p.ID)

					err := g.store.GivePlayersPrize(p.ID, p.Prize, p.Place, g.GameData.ID)

					if err != nil {
						fmt.Printf("could not give players's prize %s", err)
					}
				}

			}

			

			bytes, err := json.Marshal(res)

			if err != nil {
				fmt.Printf("could not marshal data for winner response %s", err)
			}

			g.hub.Broadcast <- []byte(bytes)

			time.Sleep(time.Second * 6)

			

			cardDest <- true

		// case <- flopCh:
		// 	g.GameData.PrepareToCardDeal()
			
		// 	g.GameData.RevealFlop()
			

			

		// case <- turnCh:
		// 	g.GameData.PrepareToCardDeal()
			
		// 	g.GameData.RevealTurn()

			
		// case <- riverCh:
		// 	g.GameData.PrepareToCardDeal()
			
		// 	g.GameData.RevealRiver()

		default:

			

			s := g.GameData.ShouldFindWinner()

			if s {
				g.GameData.FindWinner()

				g.hub.Broadcast <- []byte("reveal-cards")

				time.Sleep(time.Second * 1)

				winnerCh <- true

				continue
			}

			a := g.GameData.IsGameOnAllInStage()

			if a {
				fmt.Println("IN")
				flop := g.GameData.ShouldRevealFlop()

				if flop == true {
					g.GameData.PrepareToCardDeal()
				
					g.GameData.RevealFlop()

					g.hub.Broadcast <- []byte("global-changes")

					time.Sleep(time.Second * 3)

					continue

				}

				turn := g.GameData.ShouldRevealTurn()


				if turn == true {
					g.GameData.PrepareToCardDeal()
				
					g.GameData.RevealTurn()

					g.hub.Broadcast <- []byte("global-changes")

					time.Sleep(time.Second * 3)
					continue

				}

				river := g.GameData.ShouldRevealRiver()

			

				if river == true {
					g.GameData.PrepareToCardDeal()
				
					g.GameData.RevealRiver()

					g.hub.Broadcast <- []byte("global-changes")

					time.Sleep(time.Second * 3)
					continue
				}

				continue

			}

			

			v := g.GameData.AllFoldWinnerId()

			if v != -1 {
				g.GameData.DetermineWinner(v)

				

				winnerCh <- true

				continue
			}

			fmt.Println("333")

			g.GameData.ShouldRevealNextTableCard()

			fmt.Println("333")

	


				g.GameData.Mutex.Lock()

				actIdx := utils.SliceIndex(len(g.GameData.Players), func (i int) bool { return g.GameData.Players[i].Status == "active" })


				idx := 0

				if actIdx == -1 {
					idx = g.GameData.FindActivePlayerAfterCardDeal()
				} else {
					idx = g.GameData.DetermineNextActivePlayer(actIdx)
				}

				g.GameData.Mutex.Unlock()
	

				if idx == -1 {
					g.GameData.Mutex.Lock()

					g.GameData.Players[actIdx].MakeAction("fold")

					g.GameData.Mutex.Lock()
					flopCh <- true
					continue
				}

				g.hub.Broadcast <- []byte("global-changes")

					

					t := g.GameData.Players[idx].InitTimer(&g.GameData.Mutex)

					g.GameData.Mutex.RLock()

					if t != 0 && g.GameData.Players[idx].Action == "" {
						g.GameData.Mutex.RUnlock()

						g.GameData.Mutex.Lock()

						g.GameData.Players[idx].MakeAction(g.GameData.Players[idx].NextAction)

						g.GameData.Mutex.Unlock()

						g.hub.Broadcast <- []byte("global-changes")

					} else {
						g.GameData.Mutex.RUnlock()
					}

					



					if t == 0 {
						
						g.GameData.Mutex.RLock()
						if (g.GameData.Players[idx].TimeBank > 0) {
							g.GameData.Mutex.RUnlock()
							
		
							t = g.GameData.Players[idx].InitTimeBank(&g.GameData.Mutex)
	
							
							if t == 0 {
								g.GameData.Mutex.Lock()
								val := g.GameData.AllowSystemToCheck(idx)
	
								fmt.Println("val", val)
	
								if val {
									g.GameData.Players[idx].MakeAction("check")
								} else {
	
									g.GameData.Players[idx].MakeAction("fold")
								}
								g.GameData.Mutex.Unlock()
	
							}
		
						} else {
							g.GameData.Mutex.RUnlock()

							g.GameData.Mutex.Lock()

							val := g.GameData.AllowSystemToCheck(idx)
	
							if val {
								g.GameData.Players[idx].MakeAction("check")
							} else {
	
								g.GameData.Players[idx].MakeAction("fold")
							}

							g.GameData.Mutex.Unlock()
						}
						
						
					}
			
		}

	}
}

func (g *PokerGame) InitAnteMinBetHandler(gameEnds chan bool) {
	v := 0

	g.GameData.Mutex.RLock()

	interval := 600

	if g.GameData.Mode == "turbo" {
		interval = 300 
	}

	if g.GameData.Mode == "hyper-turbo" {
		interval = 120
	} 

	g.GameData.Mutex.RUnlock()

	fmt.Println("interval", interval)


	for {
		select {
		case <- gameEnds:
			fmt.Println("InitAnteMinBetHandler ENDS--")
			return

		default:
			
	
		
	
			if v != 0 && v % interval == 0 {
				g.GameData.Mutex.Lock()
	
				fmt.Println("old min bet", g.GameData.MinBet)
	
				g.GameData.IncreaseMinBet()
				g.GameData.IncreaseAntePay()
	
				fmt.Println("new min bet", g.GameData.MinBet)
	
				res := &struct{
					Action string `json:"action"`
				} {
					Action: "limit-changes",
				}
	
				bytes, _ := json.Marshal(res)
	
				g.hub.Broadcast <- []byte(bytes)
	
	
				
				g.GameData.Mutex.Unlock()
	
			}
	
	
			v++
	
	
			time.Sleep(time.Second)
		}
	}

}