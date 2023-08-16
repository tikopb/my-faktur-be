package model

import "time"

type Partner struct {
	ID        int       `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	CreatedBy string    `gorm:"column:created_by" json:"created_by"`
	User      User      `gorm:"foreignKey:created_by"`
	DNAmount  float64   `gorm:"column:dn_amount" json:"dn_amount"`
	CNAmount  float64   `gorm:"column:cn_amount" json:"cn_amount"`
	Isactive  bool      `gorm:"column:cn_amount" json:"isactive"`
}
