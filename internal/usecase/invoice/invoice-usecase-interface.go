package invoice

import "bemyfaktur/internal/model"

type InvoiceUsecaseInterface interface {
	IndexInvoice(limit int, offset int) ([]model.InvoiceIndexRespont, error)
	GetInvoice(id int) (model.Invoice, error)
	CreateInvoice(request model.Invoice) (model.InvoiceCreateRespon, error)
	UpdatedInvoice(id int, request model.Invoice) (model.Invoice, error)
	DeleteInvoice(id int) (string, error)

	IndexLine(limit int, offset int, invoiceId int) ([]model.InvoiceLine, error)
	GetInvoiceLine(id int) (model.InvoiceLine, error)
	CreateInvoiceLine(request model.InvoiceLine) (model.InvoiceLine, error)
	UpdatedInvoiceLine(id int, request model.InvoiceLine, productId int) (model.InvoiceLine, error)
	DeleteInvoiceLine(id int) (string, error)
}
