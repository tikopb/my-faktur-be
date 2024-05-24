package organization

import (
	"bemyfaktur/internal/model"
	"errors"
	"fmt"
	"time"

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

// On delete system cant delete the data system just will inactive the orgazniation data
func (o *OrganizationRepo) Delete(request model.OrganizationRequest) (bool, error) {
	data, err := o.ShowInternal(request.ID)
	if err != nil {
		return false, err
	}

	//updarte data to inactive
	data.IsActive = false
	if err := o.db.Updates(&data).Where(&model.Organization{UUID: request.ID}).Error; err != nil {
		return false, err
	}

	errorMsg := fmt.Sprintf("the data of %s success to inactive ", request.Name)
	return true, errors.New(errorMsg)
}

// Partial implements Repository.
func (o *OrganizationRepo) Partial(request model.OrganizationRequest) (model.OrganizationRespont, error) {
	panic("unimplemented")
}

// Show implements Repository.
func (o *OrganizationRepo) Show(request model.OrganizationRequest) (model.OrganizationRespont, error) {
	var data model.Organization

	if err := o.db.Preload("User").Where(model.Organization{UUID: request.ID}).First(&data).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.OrganizationRespont{}, errors.New("data not found")
		}
	}

	return o.ParsingOrganizationToRespont(data), nil
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
	//get data!
	data, err := o.ShowInternal(request.ID)
	if err != nil {
		return model.OrganizationRespont{}, err
	}

	data.UpdateAt = time.Now()
	data.UpdatedBy = request.UserId
	data.Name = request.Name
	data.Description = request.Description

	if err := o.db.Updates(&data).Where(&model.Organization{UUID: data.UUID}).Error; err != nil {
		return model.OrganizationRespont{}, err
	}

	//set return data
	dataReturn := o.ParsingOrganizationToRespont(data)
	return dataReturn, nil

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

// GetOrgByUserId implements Repository.
func (o *OrganizationRepo) GetOrgByUserId(userId string) (model.OrganizationRespont, error) {
	var data model.Organization
	if err := o.db.Preload("User").Where(model.Organization{CreatedBy: userId}).First(&data).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.OrganizationRespont{}, errors.New("data not found")
		}
	}

	return o.ParsingOrganizationToRespont(data), nil
}
