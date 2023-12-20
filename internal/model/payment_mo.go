package model

import (
	"strconv"
	"time"

	"bemyfaktur/internal/model/constant"

	"github.com/google/uuid"
)

type Payment struct {
	ID           int                       `json:"-" gorm:"primaryKey;autoIncrement"`
	CreatedAt    time.Time                 `gorm:"column:created_at;default:CURRENT_TIMESTAMP"`
	CreatedBy    string                    `gorm:"column:created_by" json:"created_by"`
	User         User                      `gorm:"foreignKey:created_by"`
	UpdatedBy    string                    `gorm:"column:updated_by" json:"updated_by"`
	UserUpdated  User                      `gorm:"foreignKey:updated_by"`
	PartnerID    int                       `json:"partner_id" gorm:"column:partner_id;index:idx_payment_partner_id"`
	Partner      Partner                   `gorm:"foreignKey:partner_id"`
	GrandTotal   float64                   `gorm:"column:grand_total;default:0"`
	Discount     float64                   `gorm:"column:discount;default:0"`
	BatchNo      string                    `json:"batchno" gorm:"column:batch_no"`
	Status       constant.PaymentStatus    `gorm:"column:status;default:DR;index:idx_payment_docstatus"`
	DocAction    constant.PaymentDocAction `gorm:"column:docaction;default:DR"`
	DocumentNo   string                    `json:"documentno" gorm:"column:documentno;not null;unique;index:idx_payment_documentno"`
	IsPrecentage bool                      `gorm:"column:isprecentage;default:false" json:"isprecentage"`
	UUID         uuid.UUID                 `json:"id" gorm:"type:uuid;default:uuid_generate_v4();index:idx_payment_uuid"`
}

type PaymentRequest struct {
	CreatedBy    string                    `json:"created_by"`
	PartnerID    int                       `json:"partner_id"`
	GrandTotal   float64                   `json:"grand_total"`
	Discount     float64                   `json:"discount"`
	BatchNo      string                    `json:"batchno"`
	Status       constant.PaymentStatus    `json:"status"`
	DocAction    constant.PaymentDocAction `json:"docaction"`
	DocumentNo   string                    `json:"documentno"`
	IsPrecentage bool                      `json:"isprecentage"`
}

type PaymentRespont struct {
	ID           uuid.UUID                 `json:"id"`
	DocumentNo   string                    `json:"documentno"`
	BatchNo      string                    `json:"batchno"`
	IsPrecentage bool                      `json:"isprecentage"`
	GrandTotal   float64                   `json:"grand_total"`
	Discount     float64                   `json:"discount"`
	Status       constant.PaymentStatus    `json:"status"`
	DoAction     constant.PaymentDocAction `json:"docaction"`
	CreatedBy    UserPartial               `json:"createdby"`
	UpdatedBy    UserPartial               `json:"updatedby"`
	Partner      PartnerPartialRespon      `json:"partner"`
}

type PaymentPartialRespont struct {
	UUID       uuid.UUID `json:"id"`
	BatchNo    string    `json:"batchno"`
	DocumentNo string    `json:"documentno"`
	Id         int       `json:"-"`
}

type PaymentLine struct {
	ID           int       `json:"-" gorm:"primaryKey;autoIncrement"`
	PaymentID    int       `gorm:"column:payment_id;index:idx_payment_id"`
	Payment      Payment   `gorm:"foreignKey:payment_id"`
	Price        float64   `gorm:"column:price"`
	Amount       float64   `gorm:"column:amount"`
	CreatedAt    time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP"`
	CreatedBy    string    `gorm:"column:created_by" json:"created_by"`
	User         User      `gorm:"foreignKey:created_by"`
	InvoiceID    int       `gorm:"column:invoice_id;not null;index:idx_invoice_id" json:"invoice_id"`
	Invoice      Invoice   `gorm:"foreignKey:invoice_id"`
	Discount     float64   `gorm:"column:discount" json:"discount"`
	IsPrecentage bool      `gorm:"column:isprecentage;default:false" json:"isprecentage"`
	UUID         uuid.UUID `json:"id" gorm:"type:uuid;default:uuid_generate_v4();index:idx_paymentLine_uuid"`
}

type PaymentLineRequest struct {
	PaymentID    int       `json:"payment_id"`
	Invoice_uuid uuid.UUID `json:"invoice_id"`
	Invoice_id   int       `json:"-"`
	Price        float64   `json:"price"`
	Discount     float64   `json:"discount"`
	IsPrecentage bool      `json:"isprecentage"`
	CreatedBy    string    `json:"createdby"`
}

type PaymentLineRespont struct {
	ID           uuid.UUID `json:"id"`
	Price        float64   `json:"price"`
	Amount       float64   `json:"amount"`
	BatchNo      string    `json:"batchno"`
	Invoice_id   int       `json:"invoice_id"`
	Discount     float64   `json:"discount"`
	IsPrecentage bool      `json:"isprecentage"`
	Payment      Payment
}

func GetSeatchParamPayment(dateFrom string, dateTo string, q string) string {
	//searchParam := []string{"batch_no", "documentno", "p.name"}
	var value string = " invoices.created_at >='" + dateFrom + "'::date and invoices.created_at <='" + dateTo + "'::date"

	if q != "" {
		q = "'%" + q + "%'"
		if IsIntegerVariable(q) {
			value = " lower(batch_no)  LIKE " + q + " OR lower(documentno) LIKE " + q + " OR grand_total::TEXT LIKE " + q
		} else {
			value = " lower(batch_no)  LIKE " + q + " OR lower(documentno) LIKE " + q
		}
	}

	return value
}

func GetSeatchParamPaymentLine(q string, paymentID int) string {
	id := strconv.Itoa(paymentID)
	value := " payment_id = " + id
	if IsIntegerVariable(q) {
		q = "'%" + q + "%'"
		value = value + " amount::TEXT LIKE " + q + " price::TEXT LIKE " + q
	} else {
		value = value + ""
	}

	return value
}
