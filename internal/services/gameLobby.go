package services

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"context"
	"github.com/Toront0/poker/internal/types"

)

type GameLobbyStorer interface {
	GetAllGames() ([]*types.Game, error)
	GetGame(gameID int) (types.GameDetail, error)
	CreateGame(req *types.CreateGameReq) error
	GetGamePassword(gameId int) (string, error)
	JoinGame(id, gameID int) error
	FindGames(search string) ([]*types.FoundGame, error)
	GetEmojies() ([]*types.Emoji, error)
	GetGameResults(gameID int) (*types.GameSummary, error)
}

type gameLobbyStore struct {
	store *pgxpool.Pool
}

func NewGameLobbyStore(store *pgxpool.Pool) *gameLobbyStore {
	return &gameLobbyStore{
		store: store,
	}
}

func (s *gameLobbyStore) GetAllGames() ([]*types.Game, error) {
	games := []*types.Game{}

	rows, err := s.store.Query(context.Background(), `select games.*, (select array_agg(player_id) from game_players where game_players.game_id = games.id) as tot from games where state != 'Завершен'`)

	if err != nil {
		return games, err
	}

	

	for rows.Next() {
		g := &types.Game{}

		rows.Scan(&g.ID, &g.CreatedAt, &g.Name, &g.BuyIn, &g.AmountOfPlayers, &g.Prize, &g.IsPrivate, &g.RoomPassword, &g.Mode, &g.PrizeDestribution, &g.State, &g.PlayersInRoom)

		games = append(games, g)
	}

	return games, nil
}

func (s *gameLobbyStore) GetGame(gameID int) (types.GameDetail, error) {
	res := types.GameDetail{}

	err := s.store.QueryRow(context.Background(), `select id, created_at, name, buy_in, amount_of_players, prize, mode, prize_destribution from games where id = $1`, gameID).Scan(&res.ID, &res.CreatedAt, &res.Name, &res.BuyIn, &res.AmountOfPlayers, &res.Prize, &res.Mode, &res.PrizeDestribution)

	if err != nil {
		return res, err
	}

	rows, err := s.store.Query(context.Background(), `select t1.id, t1.username, t1.profile_img from users t1 join game_players t2 on t1.id = t2.player_id  where t2.game_id = $1`, gameID)

	for rows.Next() {
		p := types.PlayerDataFromDB{}

		rows.Scan(&p.ID, &p.Username, &p.Img)

		res.Players = append(res.Players, p)
	}

	return res, err
}

func (s *gameLobbyStore) CreateGame(req *types.CreateGameReq) error {
	var gameID int 


	err := s.store.QueryRow(context.Background(), `insert into games (name, buy_in, amount_of_players, prize, is_private, room_password, mode, prize_destribution) values ($1, $2, $3, $4, $5, $6, $7, $8) returning id`, req.Name, req.BuyIn, req.AmountOfPlayers, req.Prize, req.IsPrivate, req.RoomPassword, req.Mode, req.PrizeDestribution).Scan(&gameID)

	if err != nil {
		return err
	}

	_, err = s.store.Exec(context.Background(), `insert into game_players (player_id, game_id) values($1, $2)`, req.CreatorID, gameID)

	if err != nil {
		return err
	}


	_, err = s.store.Exec(context.Background(), `update users set money = users.money - $1 where id = $2`, req.BuyIn, req.CreatorID)


	return err
}

func (s *gameLobbyStore) JoinGame(id, gameID int) error {



	_, err := s.store.Exec(context.Background(), `insert into game_players (player_id, game_id) values($1, $2)`, id, gameID)

	if err != nil {
		return err
	}

	_, err = s.store.Exec(context.Background(), `update users set money = users.money - (select buy_in from games where id = $2) where id = $1`, id, gameID)

	return err
}

func (s *gameLobbyStore) GetGamePassword(gameId int) (string, error) {
	var res string

	err := s.store.QueryRow(context.Background(), `select room_password from games where id = $1`, gameId).Scan(&res)

	return res, err
}

func (s *gameLobbyStore) FindGames(search string) ([]*types.FoundGame, error) {
	res := []*types.FoundGame{}

	q := "%" + search + "%"

	rows, err := s.store.Query(context.Background(), `select id, name, prize, state from games where name ilike $1`, q)

	if err != nil {
		return res, err
	}

	for rows.Next() {
		r := &types.FoundGame{}

		rows.Scan(&r.ID, &r.Name, &r.Prize, &r.State)

		res = append(res, r)
	}

	return res, nil
}

func (s *gameLobbyStore) GetEmojies() ([]*types.Emoji, error) {
	res := []*types.Emoji{}

	rows, err := s.store.Query(context.Background(), `select * from emojies`)

	if err != nil {
		return res, err
	}

	for rows.Next() {
		e := &types.Emoji{}

		rows.Scan(&e.ID, &e.Emoji, &e.Title)

		res = append(res, e)

	}

	return res, nil
}

func (s *gameLobbyStore) GetGameResults(gameID int) (*types.GameSummary, error) {
	res := &types.GameSummary{}

	err := s.store.QueryRow(context.Background(), `select id, created_at, name, buy_in, amount_of_players, prize, is_private, mode, prize_destribution, state from games where id = $1`, gameID).Scan(&res.ID, &res.CreatedAt, &res.Name, &res.BuyIn, &res.AmountOfPlayers, &res.Prize, &res.IsPrivate, &res.Mode, &res.PrizeDestribution, &res.State)

	if err != nil {
		return res, err
	}


	rows, err := s.store.Query(context.Background(), `select t1.id, t1.username, t1.profile_img, t2.prize, t2.place from users t1 join game_players t2 on t2.player_id = t1.id where t2.game_id = $1`, gameID)

	if err != nil {
		return res, err
	}

	for rows.Next() {
		p := types.PlayerFinalResult{}

		rows.Scan(&p.ID, &p.Username, &p.ProfileImg, &p.Prize, &p.Place)

		res.Players = append(res.Players, p)
	}


	return res, nil

}