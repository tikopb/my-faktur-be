package model

import (
	"time"

	"bemyfaktur/internal/model/constant"
)

// -- invoice
type Invoice struct {
	ID          int       `json:"id" gorm:"primaryKey;autoIncrement"`
	CreatedAt   time.Time `gorm:"column:created_at"`
	CreatedBy   string    `gorm:"column:created_by" json:"created_by"`
	User        User      `gorm:"foreignKey:created_by"`
	PartnerID   int       `json:"partner_id" gorm:"column:partner_id"`
	Partner     Partner   `gorm:"foreignKey:partner_id"`
	GrandTotal  float64   `gorm:"column:grand_total"`
	Discount    float64   `json:"discount" gorm:"column:discount"`
	BatchNo     string    `json:"batchno" gorm:"column:batch_no"`
	InvoiceLine []InvoiceLine
	Status      constant.InvoiceStatus
}

type InvoiceRequest struct {
	CreatedAt  time.Time `json:"created_at"`
	UserId     string    `json:"user_id"` //##@ until security module fixed
	PartnerID  int       `json:"partner_id"`
	GrandTotal float64   `json:"grand_total"`
	Discount   float64   `json:"discount"`
	BatchNo    string    `json:"batchno"`
	Status     constant.InvoiceStatus
}

type InvoiceRespont struct {
	ID         int       `json:"id"`
	CreatedAt  time.Time `json:"created_at"`
	GrandTotal float64   `json:"grand_total"`
	Discount   float64   `json:"discount"`
	BatchNo    string    `json:"batchno"`
	Status     constant.InvoiceStatus
	CreatedBy  User    `json:"createdby"`
	Partner    Partner `json:"partner"`
}

// -- invoice line
type InvoiceLine struct {
	ID           int       `json:"id" gorm:"primaryKey;autoIncrement"`
	CreatedAt    time.Time `json:"created_at"`
	Price        float64   `gorm:"column:price" json:"price"`
	Discount     float64   `gorm:"column:discount" json:"discount"`
	Qty          float64   `gorm:"column:qty" json:"qty"`
	Amount       float64   `gorm:"column:amount"`
	IsPrecentage bool      `gorm:"column:isprecentage;default:false" json:"isprecentage"`
	ProductID    int       `gorm:"column:product_id" json:"product_id"` // Fixed column name
	Product      Product   `gorm:"foreignKey:ProductID"`
	InvoiceID    int       `gorm:"column:invoice_id;not null" json:"invoice_id"`
	Invoice      Invoice   `gorm:"foreignKey:invoice_id"`
	CreatedBy    string    `gorm:"column:created_by" json:"created_by"`
	User         User      `gorm:"foreignKey:created_by"`
}

type InvoiceLineRespont struct {
	Invoice         string
	Invoice_line_id int
	Created_at      time.Time
	Product_name    string
	Product_id      int
	Qty             float64
	Price           float64
	Amount          float64
	Discount        float64
}
