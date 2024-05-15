package organization

import (
	"bemyfaktur/internal/model"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrganizationRepo struct {
	db *gorm.DB
}

func GetRepository(db *gorm.DB) Repository {
	return &OrganizationRepo{
		db: db,
	}
}

// Create implements Repository.
func (o *OrganizationRepo) Create(request model.OrganizationRequest) (model.OrganizationRespont, error) {
	//translate from request to organization related code
	orgData := model.Organization{
		OrgCode:     request.OrgCode,
		Name:        request.Name,
		Description: request.Description,
		IsActive:    true,
	}

	if err := o.db.Create(&orgData).Error; err != nil {
		return model.OrganizationRespont{}, err
	}

	//parsing to respont
	dataPreload := o.ParsingOrganizationToRespont(orgData)

	return dataPreload, nil
}

// Delete implements Repository.
func (o *OrganizationRepo) Delete(request model.OrganizationRequest) error {
	panic("unimplemented")
}

// Partial implements Repository.
func (o *OrganizationRepo) Partial(request model.OrganizationRequest) (model.OrganizationRespont, error) {
	panic("unimplemented")
}

// Show implements Repository.
func (o *OrganizationRepo) Show(request model.OrganizationRequest) (model.OrganizationRespont, error) {
	panic("unimplemented")
}

// ShowInternal implements Repository.
func (o *OrganizationRepo) ShowInternal(id uuid.UUID) (model.Organization, error) {
	data := model.Organization{}

	if err := o.db.Where(model.Organization{UUID: id}).First(&data).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return data, errors.New("data not found")
		}
	}

	return data, nil
}

// Update implements Repository.
func (o *OrganizationRepo) Update(request model.OrganizationRequest) (model.OrganizationRespont, error) {
	panic("unimplemented")
}

func (o *OrganizationRepo) ParsingOrganizationToRespont(organization model.Organization) model.OrganizationRespont {
	respont := model.OrganizationRespont{
		ID:          organization.UUID,
		OrgCode:     organization.OrgCode,
		Name:        organization.Name,
		Description: organization.Description,
		IsActive:    organization.IsActive,
	}

	return respont
}
