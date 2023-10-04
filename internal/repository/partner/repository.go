package partner

import (
	"bemyfaktur/internal/model"
)

type Repository interface {
	// Index returns a list of Partners.
	// Call this method to retrieve a list of Partners.
	Index(limit int, offset int, q string) ([]model.Partner, error)

	// Create creates a new Partner.
	// Call this method to create a new Partner.
	Create(partner model.Partner) (model.PartnerRespon, error)

	// Show retrieves an Partner by its ID.
	// Call this method to retrieve a specific Partner by ID.
	Show(id int) (model.Partner, error)

	// Update updates an existing Partner by its ID.
	// Call this method to update an existing Partner by ID.
	Update(id int, updatedPartner model.Partner) (model.PartnerRespon, error)

	// Delete deletes an Partner by its ID.
	// Call this method to delete an Partner by ID.
	Delete(id int) (string, error)
}
