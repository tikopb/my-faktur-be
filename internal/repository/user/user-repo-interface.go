package user

import "bemyfaktur/internal/model"

type Repository interface {
	// Index returns a list of Partners.
	// Call this method to retrieve a list of Partners.
	Index(limit int, offset int) ([]model.User, error)

	// Create creates a new Partner.
	// Call this method to create a new Partner.
	Create(user model.User) (model.User, error)

	// Show retrieves an Partner by its ID.
	// Call this method to retrieve a specific Partner by ID.
	Show(id string) (model.User, error)

	// Update updates an existing Partner by its ID.
	// Call this method to update an existing Partner by ID.
	Update(id string, updateduser model.User) (model.User, error)

	// Delete deletes an Partner by its ID.
	// Call this method to delete an Partner by ID.
	Delete(id string) (string, error)
}
