package types

import "time"

type User struct {
	ID int `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	Username string `json:"username"`
	Email string `json:"email"`
	Password string `json:"-"`
	IsVerified bool `json:"isVerified"`
	ProfileImg string `json:"profileImg"`
	BannerImg string `json:"bannerImg"`
	MoneyTransaction time.Time `json:"moneyTransaction"`
	Money int `json:"money"`
}

type ProfileImg struct {
	IsFree bool `json:"isFree"`
	URL string `json:"url"`
}

type UserDetail struct {
	ID int `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	Username string `json:"username"`
	ProfileImg string `json:"profileImg"`
	BannerImg string `json:"bannerImg"`
}

type FoundUser struct {
	ID int `json:"id"`
	Username string `json:"username"`
	ProfileImg string `json:"profileImg"`
}