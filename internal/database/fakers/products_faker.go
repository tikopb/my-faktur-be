package fakers

import (
	"bemyfaktur/internal/model"
	"time"

	"github.com/bxcodec/faker/v4"
	"github.com/google/uuid"

	"gorm.io/gorm"
)

//func ProductFaker(db *gorm.DB) *models.Product {

func ProductFaker(db *gorm.DB) *model.Product {
	return &model.Product{
		Name:        faker.Name(),
		Value:       faker.Name(),
		Upc:         faker.Name(),
		Description: faker.Paragraph(),
		CreatedBy:   "38fa80ce-fe0c-403b-9d45-d8e1d15682a0",
		CreatedAt:   time.Time{},
		UpdateAt:    time.Time{},
		IsActive:    true,
		UUID:        uuid.New(),
	}
}
