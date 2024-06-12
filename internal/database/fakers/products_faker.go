package fakers

import (
	"bemyfaktur/internal/model"

	"gorm.io/gorm"
)

//func ProductFaker(db *gorm.DB) *models.Product {

func ProductFaker(db *gorm.DB, product model.Product) *model.Product {
	return &model.Product{
		Name:           product.Name,
		Value:          product.Value,
		Upc:            product.Upc,
		Description:    product.Description,
		CreatedBy:      product.CreatedBy,
		CreatedAt:      product.CreatedAt,
		UpdateAt:       product.UpdateAt,
		IsActive:       product.IsActive,
		UUID:           product.UUID,
		OrganizationId: product.OrganizationId,
	}
}
