package invoice

import (
	"bemyfaktur/internal/model"
)

type InvoiceRepositoryInterface interface {
	// Index returns a list of invoices.
	// Call this method to retrieve a list of invoices.
	Index(limit int, offset int) ([]model.InvoiceRespont, error)

	// Create creates a new invoice.
	// Call this method to create a new invoice.
	Create(invoice model.InvoiceRequest, partner model.Partner) (model.InvoiceRespont, error)

	// Show retrieves an invoice by its ID.
	// Call this method to retrieve a specific invoice by ID.
	Show(id int) (model.Invoice, error)

	// Update updates an existing invoice by its ID.
	// Call this method to update an existing invoice by ID.
	Update(id int, updatedInvoice model.Invoice) (model.InvoiceRespont, error)

	// Delete deletes an invoice by its ID.
	// Call this method to delete an invoice by ID.
	Delete(id int) (string, error)

	// Index returns a list of InvoiceLine.
	// Call this method to retrieve a list of InvoiceLine.
	IndexLine(limit int, offset int, invoiceId int) ([]model.InvoiceLine, error)

	// Create creates a new InvoiceLine.
	// Call this method to create a new invoice.
	CreateLine(invoice model.InvoiceLine) (model.InvoiceLine, error)

	// Show retrieves an invoice by its ID.
	// Call this method to retrieve a specific invoice by ID.
	ShowLine(id int) (model.InvoiceLine, error)

	// Update updates an existing invoice by its ID.
	// Call this method to update an existing invoice by ID.
	UpdateLine(id int, updatedInvoiceLine model.InvoiceLine) (model.InvoiceLine, error)

	// Delete deletes an invoice by its ID.
	// Call this method to delete an invoice by ID.
	DeleteLine(id int) (string, error)
}
