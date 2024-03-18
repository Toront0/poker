package services

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"context"
	"github.com/Toront0/poker/internal/types"
	"time"
)

type UserProfileStorer interface {
	GetAllProfileImages() ([]*types.ProfileImg, error)
	ChangeProfileImage(userID int, url string) error
	GetUserByID(userID int) (*types.UserDetail, error)
	ChangeUsername(userID int, newV string) error
	FindUsers(search string) ([]*types.FoundUser, error)
	GetUserGames(userID , limit, page int) (*types.UserGames, error)
	GetLastMoneyTransactionStatus(userID int) (time.Time, error)
	GetFreeMoney(userID int) error
}

type userProfileStore struct {
	store *pgxpool.Pool
} 

func NewUserProfileStore(store *pgxpool.Pool) *userProfileStore {
	return &userProfileStore{
		store: store,
	}
}

func (s *userProfileStore) GetAllProfileImages() ([]*types.ProfileImg, error) {
	imgs := []*types.ProfileImg{}

	rows, err := s.store.Query(context.Background(), `select is_free, url from avatars`)

	if err != nil {
		return imgs, err
	}

	for rows.Next() {
		i := &types.ProfileImg{}

		rows.Scan(&i.IsFree, &i.URL)

		imgs = append(imgs, i)
	}

	return imgs, nil
}

func (s *userProfileStore) ChangeProfileImage(userID int, url string) error {

	_, err := s.store.Exec(context.Background(), `update users set profile_img = $1 where id = $2`, url, userID)

	return err
}

func (s *userProfileStore) GetUserByID(userID int) (*types.UserDetail, error) {
	u := &types.UserDetail{}

	err := s.store.QueryRow(context.Background(), `select id, created_at, username, profile_img, banner_img from users where id = $1`, userID).Scan(&u.ID, &u.CreatedAt, &u.Username, &u.ProfileImg, &u.BannerImg)

	if err != nil {
		return u, err
	}

	return u, nil
}

func (s *userProfileStore) ChangeUsername(userID int, newV string) error {
	
	_, err := s.store.Exec(context.Background(), `update users set username = $1 where id = $2`, newV, userID)

	return err
}

func (s *userProfileStore) FindUsers(search string) ([]*types.FoundUser, error) {
	res := []*types.FoundUser{}

	arg := "%" + search + "%"


	rows, err := s.store.Query(context.Background(),  `select id, username, profile_img from users where username ilike $1`, arg)

	if err != nil {
		return res, err
	}

	for rows.Next() {
		r := &types.FoundUser{}

		rows.Scan(&r.ID, &r.Username, &r.ProfileImg)

		res = append(res, r)
	}

	return res, nil
}

func (s *userProfileStore) GetUserGames(userID , limit, page int) (*types.UserGames, error) {
	res := &types.UserGames{}

	rows, err := s.store.Query(context.Background(), `select t1.id, t1.name, t1.prize, t1.state, t2.place from games t1 join game_players t2 on t1.id = t2.game_id where t2.player_id = $1 limit $2 offset $3`, userID, limit, limit * page)

	if err != nil {
		return res, err
	}

	err = s.store.QueryRow(context.Background(), `select count(*) from game_players where player_id = $1`, userID).Scan(&res.TotalGames)

	if err != nil {
		return res, err
	}

	for rows.Next() {
		r := types.UserGamePreview{}

		rows.Scan(&r.ID, &r.Name, &r.Prize, &r.State, &r.Place)

		res.Games = append(res.Games, r)
	}


	return res, nil
}

func (s *userProfileStore) GetLastMoneyTransactionStatus(userID int) (time.Time, error) {
	var res time.Time

	err := s.store.QueryRow(context.Background(), `select money_transaction from users where id = $1`, userID).Scan(&res)

	return res, err
}

func (s *userProfileStore) GetFreeMoney(userID int) error {

	_, err := s.store.Exec(context.Background(), `update users set money = users.money + 10000, money_transaction = $2 where id = $1`, userID, time.Now())

	return err

}