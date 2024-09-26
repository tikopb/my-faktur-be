package model

import (
	"strconv"
	"time"

	"bemyfaktur/internal/model/constant"

	"github.com/google/uuid"
)

type Payment struct {
	ID             int                       `json:"-" gorm:"primaryKey;autoIncrement"`
	CreatedAt      time.Time                 `gorm:"column:created_at;default:CURRENT_TIMESTAMP"`
	UpdateAt       time.Time                 `gorm:"column:updated_at;default:current_timestamp"`
	CreatedBy      string                    `gorm:"column:created_by" json:"created_by"`
	User           User                      `gorm:"foreignKey:created_by"`
	UpdatedBy      string                    `gorm:"column:updated_by" json:"updated_by"`
	UserUpdated    User                      `gorm:"foreignKey:updated_by"`
	PartnerID      int                       `json:"partner_id" gorm:"column:partner_id;index:idx_payment_partner_id"`
	Partner        Partner                   `gorm:"foreignKey:partner_id"`
	TotalLine      float64                   `gorm:"column:total_line;default:0"`
	GrandTotal     float64                   `gorm:"column:grand_total;default:0"`
	Discount       float64                   `gorm:"column:discount;default:0"`
	BatchNo        string                    `json:"batchno" gorm:"column:batch_no"`
	Status         constant.PaymentStatus    `gorm:"column:status;default:DR;index:idx_payment_docstatus"`
	DocAction      constant.PaymentDocAction `gorm:"column:docaction;default:DR"`
	DocumentNo     string                    `json:"documentno" gorm:"column:documentno;not null;unique;index:idx_payment_documentno"`
	IsPrecentage   bool                      `gorm:"column:isprecentage;default:false" json:"isprecentage"`
	UUID           uuid.UUID                 `json:"id" gorm:"type:uuid;default:uuid_generate_v4();index:idx_payment_uuid"`
	PayDate        time.Time                 `gorm:"column:pay_date;default:null"`
	OrganizationId int                       `gorm:"column:org_id;index:idx_payment_org_id"`
	Organization   *Organization             `gorm:"foreignKey:org_id"`
}

type PaymentRequest struct {
	CreatedBy      string                    `json:"-"`
	UpdatedBy      string                    `json:"-"`
	PartnerID      int                       `json:"-"`
	PartnerUUID    uuid.UUID                 `json:"partner_id"`
	Discount       float64                   `json:"discount"`
	BatchNo        string                    `json:"batchno"`
	Status         constant.PaymentStatus    `json:"status"`
	DocAction      constant.PaymentDocAction `json:"docaction"`
	DocumentNo     string                    `json:"documentno"`
	IsPrecentage   bool                      `json:"isprecentage"`
	PayDateString  string                    `json:"pay_date"`
	PayDate        time.Time                 `json:"-"`
	OrganizationId int                       `json:"-"`
}

type PaymentRespont struct {
	ID           uuid.UUID                 `json:"id"`
	DocumentNo   string                    `json:"documentno"`
	BatchNo      string                    `json:"batchno"`
	IsPrecentage bool                      `json:"isprecentage"`
	Discount     float64                   `json:"discount"`
	TotalLine    float64                   `json:"total_line"`
	GrandTotal   float64                   `json:"grand_total"`
	Status       constant.PaymentStatus    `json:"status"`
	DoAction     constant.PaymentDocAction `json:"docaction"`
	CreatedAt    time.Time                 `json:"created_at"`
	UpdateAt     time.Time                 `json:"updated_at"`
	CreatedBy    UserPartial               `json:"createdby"`
	UpdatedBy    UserPartial               `json:"updatedby"`
	Partner      PartnerPartialRespon      `json:"partner"`
	UUID         uuid.UUID                 `json:"-"`
	PayDate      time.Time                 `json:"pay_date"`
	Line         []PaymentLineRespont      `json:"line"`
}

type PaymentPartialRespont struct {
	UUID       uuid.UUID `json:"id"`
	BatchNo    string    `json:"batchno"`
	DocumentNo string    `json:"documentno"`
	Id         int       `json:"-"`
}

type PaymentRequestV2 struct {
	Header PaymentRequest       `json:"header"`
	Line   []PaymentLineRequest `json:"line"`
}

type PaymentRequestV3 struct {
	Data PaymentRequestV2   `json:"data"`
	File FileServiceRespont `json:"file"`
}

type PaymentRespontV3 struct {
	Data PaymentRespont     `json:"data"`
	File FileServiceRespont `json:"file"`
}

type PaymentLine struct {
	ID             int           `json:"-" gorm:"primaryKey;autoIncrement"`
	PaymentID      int           `gorm:"column:payment_id;index:idx_payment_id"`
	Payment        Payment       `gorm:"foreignKey:payment_id;constraint:OnDelete:CASCADE"`
	Price          float64       `gorm:"column:price"`
	Amount         float64       `gorm:"column:amount"`
	CreatedAt      time.Time     `gorm:"column:created_at;default:CURRENT_TIMESTAMP"`
	CreatedBy      string        `gorm:"column:created_by" json:"created_by"`
	User           User          `gorm:"foreignKey:created_by"`
	UpdatedBy      string        `gorm:"column:updated_by" json:"updated_by"`
	UserUpdated    User          `gorm:"foreignKey:updated_by"`
	InvoiceID      int           `gorm:"column:invoice_id;not null;index:idx_invoice_id" json:"invoice_id"`
	Invoice        Invoice       `gorm:"foreignKey:invoice_id"`
	Discount       float64       `gorm:"column:discount" json:"discount"`
	IsPrecentage   bool          `gorm:"column:isprecentage;default:false" json:"isprecentage"`
	UUID           uuid.UUID     `json:"id" gorm:"type:uuid;default:uuid_generate_v4();index:idx_paymentLine_uuid"`
	OrganizationId int           `gorm:"column:org_id;index:idx_paymentline_org_id"`
	Organization   *Organization `gorm:"foreignKey:org_id"`
}

type PaymentLineRequest struct {
	PaymentUUID    uuid.UUID `json:"payment_id"`
	PaymentID      int       `json:"-"`
	Invoice_uuid   uuid.UUID `json:"invoice_id"`
	Invoice_id     int       `json:"-"`
	Price          float64   `json:"price"`
	Discount       float64   `json:"discount"`
	IsPrecentage   bool      `json:"isprecentage"`
	CreatedBy      string    `json:"-"`
	UpdatedBy      string    `json:"-"`
	OrganizationId int       `json:"-"`
}

type PaymentLineRespont struct {
	ID           uuid.UUID             `json:"id"`
	Price        float64               `json:"price"`
	Amount       float64               `json:"amount"`
	BatchNo      string                `json:"batchno"`
	Invoice_id   int                   `json:"-"`
	Discount     float64               `json:"discount"`
	IsPrecentage bool                  `json:"isprecentage"`
	Payment      PaymentPartialRespont `json:"payment"`
	CreatedBy    UserPartial           `json:"createdby"`
	UpdatedBy    UserPartial           `json:"updatedby"`
	Invoice      InvoicePartialRespont `json:"invoice"`
}

func GetSeatchParamPayment(dateFrom string, dateTo string, q string) string {
	//searchParam := []string{"batch_no", "documentno", "p.name"}
	var value string = " payments.created_at >='" + dateFrom + "'::date and payments.created_at <='" + dateTo + "'::date+1"

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
