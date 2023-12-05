package product

import (
	"bemyfaktur/internal/model"

	"github.com/google/uuid"
)

type Repository interface {
	// Index returns a list of Partners.
	// Call this method to retrieve a list of Partners.
	Index(limit int, offset int, q string, order []string) ([]model.ProductRespon, error)

	// Create creates a new Partner.
	// Call this method to create a new Partner.
	Create(product model.Product) (model.ProductRespon, error)

	// Show retrieves an Partner by its ID.
	// Call this method to retrieve a specific Partner by ID.
	Show(id uuid.UUID) (model.ProductRespon, error)
	ShowInternal(id uuid.UUID) (model.Product, error)

	// Update updates an existing Partner by its ID.
	// Call this method to update an existing Partner by ID.
	Update(id uuid.UUID, updatedProduct model.Product) (model.ProductRespon, error)

	// Delete deletes an Partner by its ID.
	// Call this method to delete an Partner by ID.
	Delete(id uuid.UUID) (string, error)

	//partial get for getting data of product related
	Partial(q string) ([]model.ProductPartialRespon, error)
}
