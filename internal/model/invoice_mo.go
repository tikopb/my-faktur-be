package model

import (
	"time"

	"bemyfaktur/internal/model/constant"
)

type Invoice struct {
	ID          int       `json:"id" gorm:"primaryKey;autoIncrement"`
	CreatedAt   time.Time `gorm:"column:created_at"`
	CreatedBy   string    `gorm:"column:created_by" json:"created_by"`
	User        User      `gorm:"foreignKey:created_by"`
	PartnerID   int       `json:"partner_id" gorm:"column:partner_id"`
	Partner     Partner   `gorm:"foreignKey:partner_id"`
	GrandTotal  float64   `gorm:"column:grand_total"`
	Discount    float64   `gorm:"column:discount"`
	BatchNo     string    `json:"batchno" gorm:"column:batch_no"`
	InvoiceLine []InvoiceLine
	Status      constant.InvoiceStatus
}

type InvoiceLine struct {
	ID        int     `json:"id" gorm:"primaryKey;autoIncrement"`
	InvoiceID int     `gorm:"column:invoice_id"`
	ProductID int     `gorm:"column:product_id"` // Fixed column name
	Product   Product `gorm:"foreignKey:ProductID"`
	Price     float64 `gorm:"column:price"`
	Discount  float64 `gorm:"column:discount"`
	Qty       float64 `gorm:"column:qty"`
	Amount    float64 `gorm:"column:amount"`
	CreatedBy string  `gorm:"column:created_by" json:"created_by"`
	User      User    `gorm:"foreignKey:created_by"`
}