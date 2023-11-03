package partner

import (
	"bemyfaktur/internal/model"
	"errors"
	"strconv"

	pgUtil "bemyfaktur/internal/model/paginationUtil"

	"github.com/google/uuid"
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
	dataValue := model.PartnerRespon{}

	if partner.Code == "" {
		var count int64
		pr.db.Table("partners").Count(&count)

		partner.Code = "BP - " + strconv.FormatInt(count, 10)
	}

	if err := pr.db.Create(&partner).Error; err != nil {
		return dataValue, err
	}

	dataValue = pr.parsingPartnerToParnerRespond(partner)
	return dataValue, nil
}

// Index implements Repository.
func (pr *partnerRepo) Index(limit int, offset int, q string) ([]model.PartnerRespon, error) {
	var dataReturn []model.PartnerRespon
	var data []model.Partner

	if q != "" {
		query := " select * from partners " + pr.pgUtilRepo.HandlingPaginationWhere(model.GetSeatchParamPartner(), q, "", "")
		if err := pr.db.Raw(query).Scan(&data).Error; err != nil {
			return dataReturn, err
		}
	} else {
		if err := pr.db.Order("name").Find(&data).Error; err != nil {
			return dataReturn, err
		}
	}

	//parsing to responFormat
	for _, partner := range data {
		dataReturn = append(dataReturn, pr.parsingPartnerToParnerRespond(partner))

	}

	return dataReturn, nil
}

// Show implements Repository.
func (pr *partnerRepo) Show(id uuid.UUID) (model.Partner, error) {
	var data model.Partner

	if err := pr.db.Preload("User").Where(&model.Partner{UUID: id}).First(&data).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return data, errors.New("data not found")
		}
	}

	return data, nil
}

// Update implements Repository.
func (pr *partnerRepo) Update(id uuid.UUID, updatedPartner model.Partner) (model.PartnerRespon, error) {
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
	data.Code = updatedPartner.Code

	//save the data
	if err := pr.db.Updates(&data).Error; err != nil {
		return dataUpdated, err
	}

	//inisiate data updated return
	dataUpdated = pr.parsingPartnerToParnerRespond(data)

	return dataUpdated, nil
}

// Delete implements Repository.
func (pr *partnerRepo) Delete(id uuid.UUID) (string, error) {
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

func (pr *partnerRepo) parsingPartnerToParnerRespond(partner model.Partner) model.PartnerRespon {
	data := model.PartnerRespon{
		ID:        partner.UUID,
		Name:      partner.Name,
		Code:      partner.Code,
		DNAmount:  partner.DNAmount,
		CNAmount:  partner.CNAmount,
		Isactive:  partner.Isactive,
		CreatedBy: partner.User.Username,
		CreatedAt: partner.CreatedAt,
	}

	return data
}
