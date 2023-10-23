package model

import (
	"strconv"
	"time"

	"bemyfaktur/internal/model/constant"
)

// -- invoice
type Invoice struct {
	ID                int       `json:"id" gorm:"primaryKey;autoIncrement"`
	CreatedAt         time.Time `gorm:"column:created_at"`
	CreatedBy         string    `gorm:"column:created_by" json:"created_by"`
	User              User      `gorm:"foreignKey:created_by"`
	PartnerID         int       `json:"partner_id" gorm:"column:partner_id"`
	Partner           Partner   `gorm:"foreignKey:partner_id"`
	GrandTotal        float64   `gorm:"column:grand_total"`
	Discount          float64   `json:"discount" gorm:"column:discount"`
	BatchNo           string    `json:"batchno" gorm:"column:batch_no"`
	InvoiceLine       []InvoiceLine
	Status            constant.InvoiceStatus    `gorm:"column:status;default:DR"`
	DocAction         constant.InvoiceDocAction `json:"docaction" gorm:"column:docaction;default:DR"`
	OustandingPayment float64                   `json:"oustanding" gorm:"column:oustanding_payment"`
	DocumentNo        string                    `json:"documentno" gorm:"column:documentno;not null;unique"`
	IsPrecentage      bool                      `gorm:"column:isprecentage;default:false" json:"isprecentage"`
}

type InvoiceRequest struct {
	CreatedAt    time.Time                 `json:"created_at"`
	CreatedBy    string                    `json:"-"`
	UserId       string                    `json:"user_id"` //##@ until security module fixed
	PartnerID    int                       `json:"partner_id"`
	GrandTotal   float64                   `json:"grand_total"`
	Discount     float64                   `json:"discount"`
	BatchNo      string                    `json:"batchno"`
	Status       constant.InvoiceStatus    `json:"status"`
	DocAction    constant.InvoiceDocAction `json:"docaction"`
	DocumentNo   string                    `json:"documentno"`
	IsPrecentage bool                      `json:"isprecentage"`
}

type InvoiceRespont struct {
	ID                int       `json:"id"`
	CreatedAt         time.Time `json:"created_at"`
	GrandTotal        float64   `json:"grand_total"`
	Discount          float64   `json:"discount"`
	BatchNo           string    `json:"batchno"`
	Status            constant.InvoiceStatus
	DocAction         constant.InvoiceDocAction
	CreatedBy         User    `json:"createdby"`
	Partner           Partner `json:"partner"`
	OustandingPayment float64 `json:"oustanding"`
	DocumentNo        string  `json:"documentno"`
	IsPrecentage      bool    `json:"isprecentage"`
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
	Invoice_id      int
	Invoice_line_id int
	Created_at      time.Time
	Product_name    string
	Product_id      int
	Qty             float64
	Price           float64
	Amount          float64
	Discount        float64
	IsPrecentage    bool
}

func GetSeatchParamInvoice(q string) string {
	var value string
	q = "'%" + q + "%'"
	if IsIntegerVariable(q) {
		value = " lower(batch_no)  LIKE " + q + " OR lower(documentno) LIKE " + q + " OR grand_total::TEXT LIKE " + q
	} else {
		value = " lower(batch_no)  LIKE " + q + " OR lower(documentno) LIKE " + q
	}

	return value
}

func GetSeatchParamInvoiceLine(q string, invoiceId int) string {
	id := strconv.Itoa(invoiceId)
	value := " invoice_id = " + id
	if IsIntegerVariable(q) {
		q = "'%" + q + "%'"
		value = value + " lower(batch_no)  LIKE " + q + " OR lower(documentno) LIKE " + q + " OR grand_total::TEXT LIKE " + q
	} else {
		value = value + ""
	}

	return value
}
