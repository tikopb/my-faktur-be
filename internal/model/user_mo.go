package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        string    `gorm:"primaryKey;autoIncrement" json:"-"`
	Username  string    `gorm:"unique" json:"username"`
	Hash      string    `json:"-"`
	FullName  string    `gorm:"column:full_name" json:"full_name"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	IsActive  bool      `gorm:"column:isactive" json:"isactive"`
	// OrganizationId int           `gorm:"column:org_id;index:idx_user_org_id"`
	// Organization   *Organization `gorm:"foreignKey:org_id"`
}

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	FullName string `json:"fullname"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserSession struct {
	AccessToken      string    `json:"access_token"`
	RefreshToken     string    `json:"refresh_token"`
	OrganizationID   int       `json:"org_id"`
	OrganizationUUID uuid.UUID `json:"org_uuid"`
}

type UserSessionRespond struct {
	UserInformation User        `json:"user_information"`
	UserSession     UserSession `json:"user_session"`
}

type UserPartial struct {
	UserId           string `json:"user_uuid"`
	Username         string `json:"username"`
	OrganizationUUID string `json:"org_uuid"`
	OrganizationID   int    `json:"org_id"`
}

type RefreshSession struct {
	RefreshToken string `json:"refresh_token"`
}

type AccessToken struct {
	AccessToken string `json:"access_token"`
}
