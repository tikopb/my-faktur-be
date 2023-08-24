package invoice

import (
	"bemyfaktur/internal/model"
)

type InvoiceRepositoryInterface interface {
	// Index returns a list of invoices.
	// Call this method to retrieve a list of invoices.
	Index(limit int, offset int) ([]model.InvoiceIndexRespont, error)

	// Create creates a new invoice.
	// Call this method to create a new invoice.
	Create(invoice model.Invoice) (model.InvoiceCreateRespon, error)

	// Show retrieves an invoice by its ID.
	// Call this method to retrieve a specific invoice by ID.
	Show(id int) (model.Invoice, error)

	// Update updates an existing invoice by its ID.
	// Call this method to update an existing invoice by ID.
	Update(id int, updatedInvoice model.Invoice) (model.Invoice, error)

	// Delete deletes an invoice by its ID.
	// Call this method to delete an invoice by ID.
	Delete(id int) (string, error)
}
