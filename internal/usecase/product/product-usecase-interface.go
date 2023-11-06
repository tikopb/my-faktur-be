package product

import (
	"bemyfaktur/internal/model"

	"github.com/google/uuid"
)

type ProductUsecaseInterface interface {
	IndexProduct(limit int, offset int, q string) ([]model.ProductRespon, error)
	GetProduct(id uuid.UUID) (model.Product, error)
	CreateProduct(request model.Product, userId string) (model.ProductRespon, error)
	UpdatedProduct(id uuid.UUID, request model.Product) (model.ProductRespon, error)
	DeleteProduct(id uuid.UUID) (string, error)
}
