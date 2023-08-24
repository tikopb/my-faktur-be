package invoice

import "bemyfaktur/internal/model"

type InvoiceUsecaseInterface interface {
	IndexInvoice(limit int, offset int) ([]model.InvoiceIndexRespont, error)
	GetInvoice(id int) (model.Invoice, error)
	CreateInvoice(request model.Invoice) (model.InvoiceCreateRespon, error)
	UpdatedInvoice(id int, request model.Invoice) (model.Invoice, error)
	DeleteInvoice(id int) (string, error)
}
