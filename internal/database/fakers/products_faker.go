package fakers

import (
	"bemyfaktur/internal/model"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

//func ProductFaker(db *gorm.DB) *models.Product {

func ProductFaker(db *gorm.DB, product model.Product) *model.Product {
	return &model.Product{
		ID:          product.ID,
		Name:        product.Name,
		Value:       product.Value,
		Upc:         product.Upc,
		Description: product.Description,
		CreatedBy:   product.CreatedBy,
		CreatedAt:   time.Time{},
		UpdateAt:    time.Time{},
		IsActive:    product.IsActive,
		UUID:        uuid.New(),
	}
}
