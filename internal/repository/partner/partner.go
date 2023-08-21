package partner

import (
	"bemyfaktur/internal/model"
	"errors"

	"gorm.io/gorm"
)

type partnerRepo struct {
	db *gorm.DB
}

func GetRepository(db *gorm.DB) Repository {
	return &partnerRepo{
		db: db,
	}
}

// Create implements Repository.
func (pr *partnerRepo) Create(partner model.Partner) (model.PartnerRespon, error) {
	data := model.PartnerRespon{}
	if err := pr.db.Create(&partner).Error; err != nil {
		return data, err
	}

	data = model.PartnerRespon{
		Name:     partner.Name,
		DNAmount: partner.DNAmount,
		CNAmount: partner.CNAmount,
		Isactive: partner.Isactive,
	}
	return data, nil
}

// Index implements Repository.
func (pr *partnerRepo) Index() ([]model.Partner, error) {
	var data []model.Partner

	if err := pr.db.Find(&data).Error; err != nil {
		return data, err
	}

	return data, nil
}

// Show implements Repository.
func (pr *partnerRepo) Show(id int) (model.Partner, error) {
	var data model.Partner

	if err := pr.db.Where(model.Partner{ID: id}).Preload("Partner").First(&data).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return data, errors.New("data not found")
		}
	}
	return data, nil
}

// Update implements Repository.
func (pr *partnerRepo) Update(id int, updatedPartner model.Partner) (model.Partner, error) {
	data, err := pr.Show(id)

	if err != nil {
		return data, err
	}
	data.Name = updatedPartner.Name
	data.CreatedAt = updatedPartner.CreatedAt
	data.CreatedBy = updatedPartner.CreatedBy
	data.DNAmount = updatedPartner.DNAmount
	data.CNAmount = updatedPartner.CNAmount
	data.Isactive = updatedPartner.Isactive

	//save the data
	if err := pr.db.Save(&data).Error; err != nil {
		return data, err
	}

	return data, nil
}

// Delete implements Repository.
func (pr *partnerRepo) Delete(id int) (string, error) {
	data, err := pr.Show(id)
	name := data.Name

	if err != nil {
		return "", err
	}
	if err := pr.db.Delete(&data).Error; err != nil {
		return "", err
	}
	return name, nil
}
