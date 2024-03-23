package services

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"context"
	"github.com/Toront0/poker/internal/types"

	"fmt"
)

type AuthStorer interface {
	CreateUser(username, email, password string) (*types.AuthUser, error)
	GetUserBy(columnName string, value interface{}) (*types.AuthUser, error)
	InsertEmailCode(email string, code int) error
	DeleteCodeIfExist(email string) error
	VerifyCode(email string, code string) (bool, error)
	ChangePassword(password string) error
}

type authStore struct {
	conn *pgxpool.Pool
}

func NewAuthStore(conn *pgxpool.Pool) *authStore {
	return &authStore{
		conn: conn,
	}
}

func (s *authStore) CreateUser(username, email, password string) (*types.AuthUser, error) {
	acc := &types.AuthUser{}

	err := s.conn.QueryRow(context.Background(), `insert into users (username, email, password) values($1, $2, $3) returning id, username, profile_img, money`, username, email, password).Scan(&acc.ID, &acc.Username, &acc.ProfileImg, &acc.Money)

	// defaultTime := time.Date(1970, time.January, 1, 23, 0, 0, 0, time.UTC)

	// acc.VipFinishedAt = &defaultTime

	if err != nil {
		fmt.Printf("could not create the user %s", err)
		return &types.AuthUser{}, err
	}


	return acc, nil
}

func (s *authStore) GetUserBy(columnName string, value interface{}) (*types.AuthUser, error) {
	acc := &types.AuthUser{}

	query := fmt.Sprintf("select id, username, password, profile_img, money, (select finished_at from vip_subscriptions where user_id = t1.id) from users t1 where %s = $1", columnName)

	err := s.conn.QueryRow(context.Background(), query, value).Scan(&acc.ID, &acc.Username, &acc.Password, &acc.ProfileImg, &acc.Money, &acc.VipFinishedAt)

	if err != nil {
		fmt.Printf("could not get an user %s", err)
		return &types.AuthUser{}, err
	}

	return acc, nil
}

func (s *authStore) InsertEmailCode(email string, code int) error {

	_, err := s.conn.Exec(context.Background(), `insert into email_codes (email, code) values ($1, $2)`, email, code)

	return err
}

func (s *authStore) DeleteCodeIfExist(email string) error {

	_, err := s.conn.Exec(context.Background(), `delete from email_codes where email = $1`, email)

	return err
}

func (s *authStore) VerifyCode(email string, code string) (bool, error) {
	var res int
	
	err := s.conn.QueryRow(context.Background(), `select id from email_codes where email = $1 and code = $2`, email, code).Scan(&res)
	
	if err != nil {
		return false, err
	}

	if res == 0 {
		return false, nil
	} else {
		return true, nil
	}
}

func (s *authStore) ChangePassword(password string) error {

	_, err := s.conn.Exec(context.Background(), `update users set password = $1`, password)

	return err
}