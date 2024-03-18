package types

import (
	"time"
)

type Game struct {
	ID int `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	Name string `json:"name"`
	BuyIn int `json:"buyIn"`
	AmountOfPlayers int `json:"amountOfPlayers"`
	Prize int `json:"prize"`
	IsPrivate bool `json:"isPrivate"`
	RoomPassword string `json:"-"`
	PlayersInRoom []int `json:"playersInRoom"`
	Mode string `json:"mode"`
	PrizeDestribution string `json:"prizeDestribution"`
	State string `json:"state"`
}

type GameDetail struct {
	ID int `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	Name string `json:"name"`
	BuyIn int `json:"buyIn"`
	AmountOfPlayers int `json:"amountOfPlayers"`
	Prize int `json:"prize"`
	Mode string `json:"mode"`
	PrizeDestribution string `json:"prizeDestribution"`
	Players []PlayerDataFromDB `json:"players"`
}

type PlayerDataFromDB struct {
	ID int `json:"id"`
	Username string `json:"username"`
	Img string `json:"img"`
}

type CreateGameReq struct {
	Name string `json:"name"`
	BuyIn int `json:"buyIn"`
	AmountOfPlayers int `json:"amountOfPlayers"`
	Prize int `json:"prize"`
	IsPrivate bool `json:"isPrivate"`
	RoomPassword string `json:"roomPassword"`
	AutoStart bool `json:"autoStart"`
	CreatorID int `json:"creatorID"`
	Mode string `json:"mode"`
	PrizeDestribution string `json:"prizeDestribution"`
}

type FoundGame struct {
	ID int `json:"id"`
	Name string `json:"name"`
	Prize int `json:"prize"`
	State string `json:"state"`
}

type UserGames struct {
	Games []UserGamePreview `json:"games"`
	TotalGames int `json:"totalGames"`
}

type UserGamePreview struct {
	ID int `json:"id"`
	Name string `json:"name"`
	Prize int `json:"prize"`
	State string `json:"state"`
	Place int `json:"place"`
}

type Emoji struct {
	ID int `json:"id"`
	Emoji string `json:"emoji"`
	Title string `json:"title"`
}

type GameSummary struct {
	ID int `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	Name string `json:"name"`
	BuyIn int `json:"buyIn"`
	AmountOfPlayers int `json:"amountOfPlayers"`
	Prize int `json:"prize"`
	IsPrivate bool `json:"isPrivate"`
	Mode string `json:"mode"`
	PrizeDestribution string `json:"prizeDestribution"`
	State string `json:"state"`
	Players []PlayerFinalResult `json:"players"`
}

type PlayerFinalResult struct {
	ID int `json:"id"`
	Username string `json:"username"`
	ProfileImg string `json:"profileImg"`
	Prize int `json:"prize"`
	Place int `json:"place"`
}