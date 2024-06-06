package fakers

import (
	"bemyfaktur/internal/model"
	"time"

	"gorm.io/gorm"
)

//func ProductFaker(db *gorm.DB) *models.Product {

func PartnerFaker(db *gorm.DB, partner model.Partner) *model.Partner {

	return &model.Partner{
		Name:           partner.Name,
		CreatedAt:      time.Time{},
		CreatedBy:      partner.CreatedBy,
		DNAmount:       0,
		CNAmount:       0,
		Isactive:       true,
		Code:           partner.Code,
		UUID:           partner.UUID,
		OrganizationId: partner.OrganizationId,
	}
}
