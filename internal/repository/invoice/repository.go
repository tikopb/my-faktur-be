package invoice

import (
	"bemyfaktur/internal/model"

	"github.com/google/uuid"
)

type InvoiceRepositoryInterface interface {
	//Header
	Index(limit int, offset int, q string, order []string) ([]model.InvoiceRespont, error)
	Create(invoice model.InvoiceRequest, partner model.Partner) (model.InvoiceRespont, error)
	Show(id uuid.UUID) (model.InvoiceRespont, error)
	ShowInternal(id uuid.UUID) (model.Invoice, error)
	Update(id uuid.UUID, updatedInvoice model.Invoice) (model.InvoiceRespont, error)
	Delete(id uuid.UUID) (string, error)

	//Line
	IndexLine(limit int, offset int, invoiceId int, q string) ([]model.InvoiceLineRespont, error)
	CreateLine(request model.InvoiceLine) (model.InvoiceLine, error)
	ShowLine(id int) (model.InvoiceLine, error)
	UpdateLine(id int, updatedInvoiceLine model.InvoiceLine) (model.InvoiceLine, error)
	DeleteLine(id int) (string, error)

	//docValidation
	DocProcess(data model.Invoice, docaction string) (model.Invoice, error)
	CompleteIT(data model.Invoice, docaction string) (model.Invoice, error)
	ReversedIt(data model.Invoice, docaction string) (model.Invoice, error)

	//pagination
	HandlingPagination(q string, limit int, offset int) (int64, error)
	HandlingPaginationLine(q string, limit int, offset int, invoiceId int) (int64, error)
}
