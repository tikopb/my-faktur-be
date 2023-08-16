package model

import (
	"time"

	"bemyfaktur/internal/model/constant"
)

type Payment struct {
	ID          int       `json:"id" gorm:"primaryKey;autoIncrement"`
	CreatedAt   time.Time `gorm:"column:created_at"`
	CreatedBy   string    `gorm:"column:created_by" json:"created_by"`
	User        User      `gorm:"foreignKey:created_by"`
	PartnerID   int       `json:"partner_id" gorm:"column:partner_id"`
	Partner     Partner   `gorm:"foreignKey:partner_id"`
	GrandTotal  float64   `gorm:"column:grand_total"`
	Discount    float64   `gorm:"column:discount"`
	BatchNo     string    `json:"batchno" gorm:"column:batch_no"`
	PaymentLine []PaymentLine
	Status      constant.PaymentStatus
}

type PaymentLine struct {
	ID        int     `json:"id" gorm:"primaryKey;autoIncrement"`
	PaymentID int     `gorm:"column:Payment_id"`
	Price     float64 `gorm:"column:price"`
	Amount    float64 `gorm:"column:amount"`
	CreatedBy string  `gorm:"column:created_by" json:"created_by"`
	User      User    `gorm:"foreignKey:created_by"`
}
