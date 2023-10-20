package model

import "time"

type User struct {
	ID        string    `gorm:"primaryKey;autoIncrement" json:"id"`
	Username  string    `gorm:"unique" json:"username"`
	Hash      string    `json:"-"`
	FullName  string    `gorm:"column:full_name" json:"full_name"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	IsActive  bool      `gorm:"column:isactive" json:"isactive"`
}

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserSession struct {
	//JWTToken     string `json:"jwt_token"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type RefreshSession struct {
	RefreshToken string `json:"refresh_token"`
}

type AccessToken struct {
	AccessToken string `json:"access_token"`
}
