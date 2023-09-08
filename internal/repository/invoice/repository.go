package invoice

import (
	"bemyfaktur/internal/model"
)

type InvoiceRepositoryInterface interface {

	//Header
	Index(limit int, offset int) ([]model.InvoiceRespont, error)
	Create(invoice model.InvoiceRequest, partner model.Partner) (model.InvoiceRespont, error)
	Show(id int) (model.Invoice, error)
	Update(id int, updatedInvoice model.Invoice) (model.InvoiceRespont, error)
	Delete(id int) (string, error)

	//Line
	IndexLine(limit int, offset int, invoiceId int, q string) ([]model.InvoiceLineRespont, error)
	CreateLine(invoice model.InvoiceLine) (model.InvoiceLine, error)
	ShowLine(id int) (model.InvoiceLine, error)
	UpdateLine(id int, updatedInvoiceLine model.InvoiceLine) (model.InvoiceLine, error)
	DeleteLine(id int) (string, error)

	//docValidation
	DocProcess(data model.Invoice) error
	CompleteIT(data model.Invoice) error
	ReversedIt(data model.Invoice) error
}
