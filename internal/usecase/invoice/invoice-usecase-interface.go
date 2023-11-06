package invoice

import (
	"bemyfaktur/internal/model"

	"github.com/google/uuid"
)

type InvoiceUsecaseInterface interface {
	IndexInvoice(limit int, offset int, q string) ([]model.InvoiceRespont, error)
	GetInvoice(id int) (model.Invoice, error)
	CreateInvoice(request model.InvoiceRequest, userId string) (model.InvoiceRespont, error)
	UpdatedInvoice(id int, request model.Invoice) (model.InvoiceRespont, error)
	DeleteInvoice(id int) (string, error)

	IndexLine(limit int, offset int, invoiceId int, q string) ([]model.InvoiceLineRespont, error)
	GetInvoiceLine(id int) (model.InvoiceLine, error)
	CreateInvoiceLine(request model.InvoiceLine, userId string) (model.InvoiceLine, error)
	UpdatedInvoiceLine(id int, request model.InvoiceLine, productId uuid.UUID) (model.InvoiceLine, error)
	DeleteInvoiceLine(id int) (string, error)

	HandlingPagination(q string, limit int, offset int) (int64, error)
	HandlingPaginationLine(q string, limit int, offset int, invoiceId int) (int64, error)
}
