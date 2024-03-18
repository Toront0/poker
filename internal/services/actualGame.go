package services

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"context"
	"fmt"
)

type ActualGameStorer interface {
	GivePlayersPrize(playerID, prize, place, gameID int) error
	EndTheGame(gameID int) error
}

type actualGameStore struct {
	store *pgxpool.Pool
}

func NewActualGameStore(store *pgxpool.Pool) *actualGameStore {
	return &actualGameStore{
		store: store,
	}
}

func (s *actualGameStore) GivePlayersPrize(playerID, prize, place, gameID int) error {

	fmt.Println("prize playerID", prize, playerID)

	_, err := s.store.Exec(context.Background(), `update game_players set place = $1, prize = $2 where game_id = $3 and player_id = $4`, place, prize, gameID, playerID)

	if err != nil {
		
		return err
	}


	return nil
}

func (s *actualGameStore) EndTheGame(gameID int) error {

	_, err := s.store.Exec(context.Background(), `update games set state = $1 where id = $2`, "Завершен", gameID)

	return err
}