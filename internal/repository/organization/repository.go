package organization

import (
	"bemyfaktur/internal/model"

	"github.com/google/uuid"
)

type Repository interface {
	Create(request model.OrganizationRequest) (model.OrganizationRespont, error)
	Show(request model.OrganizationRequest) (model.OrganizationRespont, error)
	ShowInternal(id uuid.UUID) (model.Organization, error)
	Update(request model.OrganizationRequest) (model.OrganizationRespont, error)
	Delete(request model.OrganizationRequest) error
	Partial(request model.OrganizationRequest) (model.OrganizationRespont, error)
}
