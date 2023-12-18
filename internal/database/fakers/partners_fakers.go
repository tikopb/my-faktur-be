package fakers

import (
	"bemyfaktur/internal/model"
	"time"

	"github.com/bxcodec/faker/v4"
	"github.com/google/uuid"

	"gorm.io/gorm"
)

//func ProductFaker(db *gorm.DB) *models.Product {

func PartnerFaker(db *gorm.DB) *model.Partner {

	return &model.Partner{
		Name:      faker.Name(),
		CreatedAt: time.Time{},
		CreatedBy: "38fa80ce-fe0c-403b-9d45-d8e1d15682a0",
		DNAmount:  0,
		CNAmount:  0,
		Isactive:  true,
		Code:      faker.Name(),
		UUID:      uuid.New(),
	}
}
