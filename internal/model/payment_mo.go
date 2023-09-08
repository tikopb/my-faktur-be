package model

import (
	"time"

	"bemyfaktur/internal/model/constant"
)

type Payment struct {
	ID          int       `json:"id" gorm:"primaryKey;autoIncrement"`
	CreatedAt   time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP"`
	CreatedBy   string    `gorm:"column:created_by" json:"created_by"`
	User        User      `gorm:"foreignKey:created_by"`
	PartnerID   int       `json:"partner_id" gorm:"column:partner_id"`
	Partner     Partner   `gorm:"foreignKey:partner_id"`
	GrandTotal  float64   `gorm:"column:grand_total"`
	Discount    float64   `gorm:"column:discount"`
	BatchNo     string    `json:"batchno" gorm:"column:batch_no"`
	PaymentLine []PaymentLine
	Status      constant.PaymentStatus    `gorm:"column:status;default:DR"`
	DocAction   constant.InvoiceDocAction `gorm:"column:docaction;default:DR"`
	DocumentNo  string                    `json:"documentno" gorm:"column:documentno;not null;unique"`
}

type PaymentRequest struct {
	ID         int                       `json:"id"`
	CreatedBy  string                    `json:"created_by"`
	PartnerID  int                       `json:"partner_id"`
	GrandTotal float64                   `json:"grand_total"`
	Discount   float64                   `json:"discount"`
	BatchNo    string                    `json:"batchno"`
	Status     constant.PaymentStatus    `json:"status"`
	DoAction   constant.PaymentDocAction `json:"docaction"`
	DocumentNo string                    `json:"documentno"`
}

type PaymentRespont struct {
	ID           int                       `json:"id"`
	CreatedBy    string                    `json:"created_by"`
	PartnerID    int                       `json:"partner_id"`
	Partner_name string                    `json:"partner_name"`
	GrandTotal   float64                   `json:"grand_total"`
	Discount     float64                   `json:"discount"`
	BatchNo      string                    `json:"batchno"`
	Status       constant.PaymentStatus    `json:"status"`
	DoAction     constant.PaymentDocAction `json:"docaction"`
	Partner      Partner
}

type PaymentLine struct {
	ID        int       `json:"id" gorm:"primaryKey;autoIncrement"`
	PaymentID int       `gorm:"column:payment_id"`
	Payment   Payment   `gorm:"foreignKey:payment_id"`
	Price     float64   `gorm:"column:price"`
	Amount    float64   `gorm:"column:amount"`
	CreatedAt time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP"`
	CreatedBy string    `gorm:"column:created_by" json:"created_by"`
	User      User      `gorm:"foreignKey:created_by"`
	InvoiceID int       `gorm:"column:invoice_id;not null" json:"invoice_id"`
	Invoice   Invoice   `gorm:"foreignKey:invoice_id"`
}

type PaymentLineRequest struct {
	PaymentID  int     `json:"payment_id"`
	Invoice_id int     `json:"invoice_id"`
	Price      float64 `json:"price"`
}

type PaymentLineRespont struct {
	ID         int     `json:"id"`
	Price      float64 `json:"price"`
	Amount     float64 `json:"amount"`
	BatchNo    string  `json:"batchno"`
	Payment    Payment
	Invoice_id int `json:"invoice_id"`
}
