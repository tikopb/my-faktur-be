package invoice

import (
	"bemyfaktur/internal/model"

	"github.com/google/uuid"
)

type InvoiceRepositoryInterface interface {
	//Header
	Index(limit int, offset int, q string, order []string, dateFrom string, dateTo string) ([]model.InvoiceRespont, error)
	Create(invoice model.InvoiceRequest, partner model.Partner) (model.InvoiceRespont, error)
	Show(id uuid.UUID) (model.InvoiceRespont, error)
	ShowInternal(id uuid.UUID) (model.Invoice, error)
	Update(id uuid.UUID, request model.InvoiceRequest) (model.InvoiceRespont, error)
	Delete(id uuid.UUID) (string, error)
	Partial(partner_id int, q string) ([]model.InvoicePartialRespont, error)

	//Line
	IndexLine(limit int, offset int, invoiceId int, q string, order []string) ([]model.InvoiceLineRespont, error)
	CreateLine(request model.InvoiceLineRequest) (model.InvoiceLineRespont, error)
	ShowLine(id uuid.UUID) (model.InvoiceLineRespont, error)
	ShowLineLinternal(id uuid.UUID) (model.InvoiceLine, error)
	UpdateLine(id uuid.UUID, request model.InvoiceLineRequest) (model.InvoiceLineRespont, error)
	DeleteLine(id uuid.UUID) (string, error)

	//save validation
	BeforeUpdate(data model.Invoice, docaction string) (model.Invoice, error)
	BeforeSave(data model.Invoice) (model.Invoice, error)

	//docValidation
	DocProcess(data model.Invoice, docaction string) (model.Invoice, error)
	CompleteIT(data model.Invoice, docaction string) (model.Invoice, error)
	ReversedIt(data model.Invoice, docaction string) (model.Invoice, error)

	//pagination
	HandlingPagination(q string, limit int, offset int, dateFrom string, dateTo string) (int64, error)
	HandlingPaginationLine(q string, limit int, offset int, invoiceId int) (int64, error)

	//create simultaneously header and line
	CreateInvoiceV2(request model.InvoiceRequest, requestLines []model.InvoiceLineRequest, partner model.Partner) (model.InvoiceRespont, []model.InvoiceLineRespont, error)

	//parsing
	ParsingInvoiceToInvoiceRespont(invoice model.Invoice) (model.InvoiceRespont, error)
	ParsingInvoiceToInvoiceRequest(invoice model.Invoice) (model.InvoiceRequest, error)
}
