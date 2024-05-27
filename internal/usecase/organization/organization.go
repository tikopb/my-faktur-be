package organization

import (
	"bemyfaktur/internal/model"
	"bemyfaktur/internal/repository/organization"

	"github.com/google/uuid"
)

type Usecase interface {
	Create(request model.OrganizationRequest) (model.OrganizationRespont, error)
	Show(request model.OrganizationRequest) (model.OrganizationRespont, error)
	ShowInternal(id uuid.UUID) (model.Organization, error)
	Update(request model.OrganizationRequest) (model.OrganizationRespont, error)
	Delete(request model.OrganizationRequest) (bool, error)
	Partial(request model.OrganizationRequest) (model.OrganizationRespont, error)
	GetOrgByUserId(userId string) (model.OrganizationRespont, error)
}

type organizationUsecase struct {
	organizationRepo organization.Repository
}

func GetUsecase(organizationRepo organization.Repository) Usecase {
	return &organizationUsecase{
		organizationRepo: organizationRepo,
	}
}

// Create implements Usecase.
func (o *organizationUsecase) Create(request model.OrganizationRequest) (model.OrganizationRespont, error) {
	return o.organizationRepo.Create(request)
}

// Delete implements Usecase.
func (o *organizationUsecase) Delete(request model.OrganizationRequest) (bool, error) {
	return o.organizationRepo.Delete(request)
}

// Partial implements Usecase.
func (o *organizationUsecase) Partial(request model.OrganizationRequest) (model.OrganizationRespont, error) {
	panic("unimplemented")
}

// Show implements Usecase.
func (o *organizationUsecase) Show(request model.OrganizationRequest) (model.OrganizationRespont, error) {
	return o.organizationRepo.Show(request)
}

// ShowInternal implements Usecase.
func (o *organizationUsecase) ShowInternal(id uuid.UUID) (model.Organization, error) {
	panic("unimplemented")
}

// Update implements Usecase.
func (o *organizationUsecase) Update(request model.OrganizationRequest) (model.OrganizationRespont, error) {
	return o.organizationRepo.Update(request)
}

// GetOrgByUserId implements Usecase.
func (o *organizationUsecase) GetOrgByUserId(userId string) (model.OrganizationRespont, error) {
	return o.organizationRepo.GetOrgByUserId(userId)
}
