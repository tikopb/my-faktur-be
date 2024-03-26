package model

import (
	"strconv"
	"time"

	"bemyfaktur/internal/model/constant"

	"github.com/google/uuid"
)

// Invoice -- invoice
type Invoice struct {
	ID                int                       `json:"-" gorm:"primaryKey;autoIncrement"`
	UUID              uuid.UUID                 `json:"id" gorm:"type:uuid;default:uuid_generate_v4();index:idx_invoice_uuid"`
	CreatedAt         time.Time                 `gorm:"column:created_at;default:current_timestamp"`
	UpdateAt          time.Time                 `gorm:"column:updated_at;default:current_timestamp"`
	CreatedBy         string                    `gorm:"column:created_by;" json:"created_by"`
	User              User                      `gorm:"foreignKey:created_by"`
	UpdatedBy         string                    `gorm:"column:updated_by" json:"updated_by"`
	UserUpdated       User                      `gorm:"foreignKey:updated_by"`
	PartnerID         int                       `json:"partner_id" gorm:"column:partner_id;index:idx_invoice_partner_id"`
	Partner           Partner                   `gorm:"foreignKey:partner_id"`
	TotalLine         float64                   `gorm:"column:total_line;default:0"`
	GrandTotal        float64                   `gorm:"column:grand_total;default:0"`
	Discount          float64                   `json:"discount" gorm:"column:discount"`
	BatchNo           string                    `json:"batchno" gorm:"column:batch_no;index:idx_invoice_batchno"`
	Status            constant.InvoiceStatus    `gorm:"column:status;default:DR"`
	DocAction         constant.InvoiceDocAction `json:"docaction" gorm:"column:docaction;default:DR"`
	OustandingPayment float64                   `json:"outstanding" gorm:"column:oustanding_payment;default:0"`
	DocumentNo        string                    `json:"documentno" gorm:"column:documentno;not null;unique"`
	IsPrecentage      bool                      `gorm:"column:isprecentage;default:false" json:"isprecentage"`
	PayDate           time.Time                 `gorm:"column:pay_date"`
}

type InvoiceRequest struct {
	Discount      float64                   `json:"discount"`
	BatchNo       string                    `json:"batchno"`
	Status        constant.InvoiceStatus    `json:"status"`
	DocAction     constant.InvoiceDocAction `json:"docaction"`
	IsPrecentage  bool                      `json:"ispercentage"`
	PartnerUUID   uuid.UUID                 `json:"partner_id"`
	PartnerId     int                       `json:"-"`
	CreatedById   string                    `json:"-"`
	UpdatedById   string                    `json:"-"`
	DateFrom      time.Time                 `json:"date_from"`
	DateTo        time.Time                 `json:"date_to"`
	PayDateString string                    `json:"pay_date"`
	PayDate       time.Time                 `json:"-"`
}

type InvoiceRespont struct {
	ID                uuid.UUID                 `json:"id"`
	InvoiceId         int                       `json:"-"`
	CreatedAt         time.Time                 `json:"created_at"`
	UpdatedAt         time.Time                 `json:"updated_at"`
	BatchNo           string                    `json:"batchno"`
	Status            constant.InvoiceStatus    `json:"status"`
	DocAction         constant.InvoiceDocAction `json:"docaction"`
	OustandingPayment float64                   `json:"oustanding"`
	DocumentNo        string                    `json:"documentno"`
	IsPrecentage      bool                      `json:"ispercentage"`
	PayDate           time.Time                 `json:"pay_date"`
	TotalLine         float64                   `json:"total_line"`
	Discount          float64                   `json:"discount"`
	GrandTotal        float64                   `json:"grand_total"`
	CreatedBy         UserPartial               `json:"createdby"`
	UpdatedBy         UserPartial               `json:"updatedby"`
	Partner           PartnerPartialRespon      `json:"partner"`
	Line              []InvoiceLineRespont      `json:"line"`
}

type InvoicePartialRespont struct {
	UUID              uuid.UUID `json:"id"`
	BatchNo           string    `json:"batchno"`
	Documentno        string    `json:"documentno"`
	OustandingPayment float64   `json:"outstanding" gorm:"oustanding_payment"`
	Id                int       `json:"-"`
}

// InvoiceLine -- invoice line
type InvoiceLine struct {
	ID           int       `json:"-" gorm:"primaryKey;autoIncrement"`
	CreatedAt    time.Time `gorm:"column:created_at;default:current_timestamp"`
	UpdateAt     time.Time `gorm:"column:updated_at;default:current_timestamp"`
	CreatedBy    string    `gorm:"column:created_by;" json:"created_by"`
	User         User      `gorm:"foreignKey:created_by"`
	UpdatedBy    string    `gorm:"column:updated_by" json:"updated_by"`
	UserUpdated  User      `gorm:"foreignKey:updated_by"`
	Price        float64   `gorm:"column:price" json:"price"`
	Discount     float64   `gorm:"column:discount" json:"discount"`
	Qty          float64   `gorm:"column:qty" json:"qty"`
	Amount       float64   `gorm:"column:amount"`
	IsPrecentage bool      `gorm:"column:isprecentage;default:false" json:"ispercentage"`
	ProductID    int       `gorm:"column:product_id;index:idx_invoiceline_productId" json:"product_id"`
	Product      Product   `gorm:"foreignKey:ProductID"`
	InvoiceID    int       `gorm:"column:invoice_id;not null;index:idx_invoiceline_invoiceId" json:"invoice_id"`
	Invoice      Invoice   `gorm:"foreignKey:invoice_id;constraint:OnDelete:CASCADE"`
	UUID         uuid.UUID `json:"id" gorm:"type:uuid;default:uuid_generate_v4();index:idx_invoiceline_uuid"`
}

type InvoiceLineRequest struct {
	InvoiceUUId  uuid.UUID `json:"invoice_id"`
	InvoiceId    int       `json:"-"`
	ProductUUID  uuid.UUID `json:"product_id"`
	ProductID    int       `json:"-"`
	Qty          float64   `json:"qty"`
	Price        float64   `json:"price"`
	Amount       float64   `json:"-"`
	Discount     float64   `json:"discount"`
	IsPrecentage bool      `json:"ispercentage"`
	CreatedById  string    `json:"-"`
	UpdatedById  string    `json:"-"`
}

type InvoiceLineRespont struct {
	ID           uuid.UUID             `json:"id"`
	CreatedAt    time.Time             `json:"created_at"`
	UpdatedAt    time.Time             `json:"updated_at"`
	Qty          float64               `json:"qty"`
	Price        float64               `json:"price"`
	Amount       float64               `json:"amount"`
	Discount     float64               `json:"discount"`
	IsPrecentage bool                  `json:"ispercentage"`
	CreatedBy    UserPartial           `json:"createdby"`
	UpdatedBy    UserPartial           `json:"updatedby"`
	Invoice      InvoicePartialRespont `json:"invoice"`
	Product      ProductPartialRespon  `json:"product"`
}

type InvoiceRequestV2 struct {
	Header InvoiceRequest       `json:"header"`
	Line   []InvoiceLineRequest `json:"line"`
}
type InvoiceRespontV2 struct {
	Header InvoiceRespont       `json:"header"`
	Line   []InvoiceLineRespont `json:"line"`
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

func GetSeatchParamInvoiceV2(dateFrom string, dateTo string, q string) string {
	var value string = " invoices.created_at >='" + dateFrom + "'::date and invoices.created_at <='" + dateTo + "'::date+1"

	if q != "" { // if q not nill then add information of documentno
		q = "'%" + q + "%'"
		if IsIntegerVariable(q) {
			value = value + " and lower(batch_no)  LIKE " + q + " OR lower(documentno) LIKE " + q + " OR grand_total::TEXT LIKE " + q
		} else {
			value = value + " and lower(batch_no)  LIKE " + q + " OR lower(documentno) LIKE " + q
		}
	}

	return value
}

func GetSeatchParamInvoiceLine(q string, invoiceId int) string {
	id := strconv.Itoa(invoiceId)
	value := " invoice_id = " + id
	if IsIntegerVariable(q) {
		q = "'%" + q + "%'"
		value = value + " AND discount::TEXT LIKE " + q + " OR price::TEXT LIKE " + q + " OR amount::TEXT LIKE " + q
	} else {
		value = value + ""
	}

	return value
}
