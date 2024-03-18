package roomHandler

import (
	"github.com/Toront0/poker/internal/services"
	"github.com/Toront0/poker/internal/types"
	"github.com/Toront0/poker/internal/handlers/actualGame"

	"golang.org/x/crypto/bcrypt"

	"net/http"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
)

type gameLobbyHandler struct {
	store services.GameLobbyStorer
	hub *Hub
}

func NewGameLobbyHandler(store services.GameLobbyStorer) *gameLobbyHandler {
	hub := NewHub()

	go hub.Run()

	return &gameLobbyHandler{
		store: store,
		hub: hub,
	}
}

func (h *gameLobbyHandler) ServeWs(w http.ResponseWriter, r *http.Request) {
	conn, err := actualGame.Upgrader.Upgrade(w, r, nil)
	
	if err != nil {
		log.Println(err)
		return
	}

	client := &Client{hub: h.hub, conn: conn, send: make(chan []byte, 256)}
	h.hub.register <- client


	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
}

func (h *gameLobbyHandler) HandleGetAllGames(w http.ResponseWriter, r *http.Request) {

	games, err :=h.store.GetAllGames()

	if err != nil {
		fmt.Printf("could not get games %s", err)
		return
	}


	json.NewEncoder(w).Encode(games)
}

type Res struct {
	Action string `json:"action"`
	PlayerID int `json:"playerId"`
	GameID int `json:"gameId"`
}

func (h *gameLobbyHandler) HandleCreateGame(w http.ResponseWriter, r *http.Request) {
	req := &types.CreateGameReq{}

	json.NewDecoder(r.Body).Decode(req)

	epw, err := bcrypt.GenerateFromPassword([]byte(req.RoomPassword), bcrypt.DefaultCost)

	if err != nil {
		fmt.Printf("could not hash password %s", err)
		return
	}

	if req.IsPrivate {
		req.RoomPassword = string(epw)
	}

	err = h.store.CreateGame(req)

	if err != nil {
		fmt.Printf("could not create the game %s", err)
		return
	}
}

func (h *gameLobbyHandler) HandleJoinGame(w http.ResponseWriter, r *http.Request, mux *chi.Mux, store services.ActualGameStorer) {
	req := &struct{
		PlayerID int `json:"playerId"`
		GameID int `json:"gameId"`
		Password string `json:"password"`
		PlayersInRoom []int `json:"playersInRoom"`
		AmountOfPlayers int `json:"amountOfPlayers"`
	} {
		PlayerID: 0,
		GameID: 0,
		Password: "",
		PlayersInRoom: []int{},
		AmountOfPlayers: 0,
	}

	json.NewDecoder(r.Body).Decode(req)

	pwrd, err := h.store.GetGamePassword(req.GameID)



	if err != nil {
		fmt.Printf("could not check password of the game %s", err)
		w.WriteHeader(400)
		return
	}

	if pwrd != "" {
		err =  bcrypt.CompareHashAndPassword([]byte(pwrd), []byte(req.Password))

		if err != nil {
			fmt.Printf("password of the room is wrong %s", err)
			w.WriteHeader(409)
			return 
		}
	} 


	err = h.store.JoinGame(req.PlayerID, req.GameID)

	if err != nil {
		fmt.Printf("could not join the game %s", err)
		w.WriteHeader(400)
		return
	}


	s := &Res{
		Action: "player-in",
		PlayerID: req.PlayerID,
		GameID: req.GameID,
	}

	en, _ := json.Marshal(s)


	h.hub.broadcast <- []byte(en)


	// len(req.PlayersInRoom) + 1 Because we count requester as a potential future player
	if len(req.PlayersInRoom) + 1 == req.AmountOfPlayers {

		res, err := h.store.GetGame(req.GameID)

		if err != nil {
			fmt.Printf("could not start the game %s", err)
			w.WriteHeader(500)
			return
		}

		game := actualGame.NewPokerGame(mux, store, res)

		game.SetupRoutes()
	
		go game.Run()

		s := &GameStartRes{
			Action: "game-start",
			PlayerID: req.PlayerID,
			GameID: req.GameID,
		}
		fmt.Println("two or more")
	
	
		en, _ := json.Marshal(s)
	
		h.hub.broadcast <- []byte(en)


		fmt.Println("two or more")

	}

}


func (h *gameLobbyHandler) HandleStartActualGame(w http.ResponseWriter, r *http.Request, mux *chi.Mux, store services.ActualGameStorer) {

	ps := []types.PlayerDataFromDB{{6, "admin",""}, { 8, "toronto", "" }}

	data := types.GameDetail{
		ID: 1,
		CreatedAt: time.Now(),
		Name: "High Rollers",
		BuyIn: 12500,
		AmountOfPlayers: 2,
		Prize: 25000,
		Mode: "hyper-turbo",
		PrizeDestribution: "winner takes all",
		Players: ps,
	}

	game := actualGame.NewPokerGame(mux, store, data)

	game.SetupRoutes()

	go game.Run()
}

func (h *gameLobbyHandler) HandleGetGameResults(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	gID, err := strconv.Atoi(id)

	if err != nil {
		fmt.Printf("could not parse user ID %s", err)
		w.WriteHeader(400)
		return
	}


	res, err := h.store.GetGameResults(gID)

	if err != nil {
		fmt.Printf("could not get game results %s", err)
		w.WriteHeader(404)
		return
	}

	json.NewEncoder(w).Encode(res)

}

func (h *gameLobbyHandler) HandleFindGames(w http.ResponseWriter, r *http.Request) {
	req := r.URL.Query().Get("search")


	res, err := h.store.FindGames(req)

	if err != nil {
		fmt.Printf("could not find the game %s", err)
		w.WriteHeader(400)
		return
	}

	json.NewEncoder(w).Encode(res)
}

func (h *gameLobbyHandler) HandleGetEmojies(w http.ResponseWriter, r *http.Request) {

	res, err := h.store.GetEmojies()

	if err != nil {
		fmt.Printf("could not get emojies %s", err)
		return
	}

	json.NewEncoder(w).Encode(res)

}


// func (h *gameLobbyHandler ) HandleTest(w http.ResponseWriter, r *http.Request) {
// 	req := &struct{
// 		Value string `json:"value"`
// 		Value2 string `json:"value2"`
// 	} {
// 		Value: ""
// 	}


// }