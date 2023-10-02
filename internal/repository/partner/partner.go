package partner

import (
	"bemyfaktur/internal/model"
	"errors"

	pgUtil "bemyfaktur/internal/model/paginationUtil"

	"gorm.io/gorm"
)

type partnerRepo struct {
	db         *gorm.DB
	pgUtilRepo pgUtil.Repository
}

func GetRepository(db *gorm.DB, pgRepo pgUtil.Repository) Repository {
	return &partnerRepo{
		db:         db,
		pgUtilRepo: pgRepo,
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
func (pr *partnerRepo) Index(limit int, offset int, q string) ([]model.Partner, error) {
	var data []model.Partner

	if q != "" {
		query := " select * from partners " + pr.pgUtilRepo.HandlingPaginationWhere(model.GetSeatchParamPartner(), q, "", "")
		if err := pr.db.Raw(query).Scan(&data).Error; err != nil {
			return data, err
		}
	} else {
		if err := pr.db.Order("name").Find(&data).Error; err != nil {
			return data, err
		}
	}

	return data, nil
}

// Show implements Repository.
func (pr *partnerRepo) Show(id int) (model.Partner, error) {
	var data model.Partner

	if err := pr.db.First(&data, id).Preload("Partner").Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return data, errors.New("data not found")
		}
	}

	return data, nil
}

// Update implements Repository.
func (pr *partnerRepo) Update(id int, updatedPartner model.Partner) (model.PartnerRespon, error) {
	dataUpdated := model.PartnerRespon{}
	data, err := pr.Show(id)

	if err != nil {
		return dataUpdated, err
	}
	data.Name = updatedPartner.Name
	data.CreatedAt = updatedPartner.CreatedAt
	data.CreatedBy = updatedPartner.CreatedBy
	data.DNAmount = updatedPartner.DNAmount
	data.CNAmount = updatedPartner.CNAmount
	data.Isactive = updatedPartner.Isactive

	//save the data
	if err := pr.db.Save(&data).Error; err != nil {
		return dataUpdated, err
	}

	//inisiate data updated return
	dataUpdated = model.PartnerRespon{
		Name:     data.Name,
		DNAmount: data.DNAmount,
		CNAmount: data.CNAmount,
		Isactive: data.Isactive,
	}
	return dataUpdated, nil
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
