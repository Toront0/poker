package actualGame

import (
	"net/http"
	"strconv"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"fmt"
	"log"

	"github.com/Toront0/poker/internal/utils"
)

func (g *PokerGame) GetGameData(w http.ResponseWriter, r *http.Request) {
	req := chi.URLParam(r, "id")

	uID, err := strconv.Atoi(req)

	if err != nil {
		fmt.Printf("could not parse user ID in GetGameData %s", err)
		w.WriteHeader(400)
		return
	}


	json.NewEncoder(w).Encode(g.GameData.SendData(uID))

}


func (g *PokerGame) Fold(w http.ResponseWriter, r *http.Request) {
	req := chi.URLParam(r, "id")
	
	id, err := strconv.Atoi(req)

	fmt.Println("id", id)

	idx := utils.SliceIndex(len(g.GameData.Players), func (i int) bool { return g.GameData.Players[i].ID == id})

	if err != nil {
		fmt.Printf("could not convert id %s", err)
		w.WriteHeader(400)
		return
	}

	g.GameData.Players[idx].MakeAction("fold")
}


func (g *PokerGame) Check(w http.ResponseWriter, r *http.Request) {
	req := chi.URLParam(r, "id")
	
	id, err := strconv.Atoi(req)


	idx := utils.SliceIndex(len(g.GameData.Players), func (i int) bool { return g.GameData.Players[i].ID == id})

	if err != nil {
		fmt.Printf("could not convert id %s", err)
		w.WriteHeader(400)
		return
	}

	g.GameData.Players[idx].MakeAction("check")
}

func (g *PokerGame) Call(w http.ResponseWriter, r *http.Request) {
	req := &struct{
		UserID int `json:"user_id"`
		Bet int `json:"bet"`
	} {
		UserID: 0,
		Bet: 0,
	}

	json.NewDecoder(r.Body).Decode(req)

	idx := utils.SliceIndex(len(g.GameData.Players), func (i int) bool { return g.GameData.Players[i].ID == req.UserID})

	g.GameData.AddChipsToPot(req.Bet)
	
	g.GameData.Players[idx].MakeAction("call")


	for i, p := range g.GameData.Players {
		if p.Action == "raise" {

			fmt.Println("HTTP RAISE")

			g.GameData.Players[i].Action = "call"

			// p.MakeAction("call")
		}


	}

	g.GameData.Players[idx].MakeBet(req.Bet)
}

// func (g *PokerGame) Emoji(w http.ResponseWriter, r *http.Request) {

// 	req := &struct {
// 		SenderID int `json:"senderId"`
// 		EmojiId
// 	}

// }

func (g *PokerGame) Raise(w http.ResponseWriter, r *http.Request) {
	req := &struct{
		UserID int `json:"user_id"`
		Bet int `json:"bet"`
	} {
		UserID: 0,
		Bet: 0,
	}

	

	json.NewDecoder(r.Body).Decode(req)

	fmt.Println("RAISING", req.Bet, req.UserID)

	idx := utils.SliceIndex(len(g.GameData.Players), func (i int) bool { return g.GameData.Players[i].ID == req.UserID})

	if req.Bet - g.GameData.Players[idx].Chips == 0 {
		g.GameData.Players[idx].MakeAction("all-in")
	} else {
		g.GameData.Players[idx].MakeAction("raise")
		
	}

	g.GameData.AddChipsToPot(req.Bet)
	g.GameData.ResetPlayersActionAfterRaise(req.UserID)
	
	
	g.GameData.Players[idx].MakeBet(req.Bet)
}

func (g *PokerGame) ServeWs(w http.ResponseWriter, r *http.Request) {


	req := chi.URLParam(r, "id")

	uID, err := strconv.Atoi(req)

	fmt.Println("uID", uID)

	if err != nil {
		fmt.Printf("could not parse playerID %s", err)
		log.Println(err)
		return
	}

	conn, err := Upgrader.Upgrade(w, r, nil)
	
	if err != nil {
		log.Println(err)
		return
	}

	client := &Client{Hub: g.hub, Conn: conn, Send: make(chan []byte, 256), PlayerID: uID}
	client.Hub.Register <- client

	fmt.Println("connection")

	

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump(g.GameData)
	go client.readPump(g.GameData)
}