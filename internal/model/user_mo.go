package model

import "time"

type User struct {
	ID        int       `gorm:"primaryKey;autoIncrement" json:"id"`
	Username  string    `gorm:"unique" json:"username"`
	Hash      string    `json:"-"`
	FullName  string    `gorm:"column:full_name" json:"full_name"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	IsActive  bool      `gorm:"column:isactive" json:"isactive"`
}
