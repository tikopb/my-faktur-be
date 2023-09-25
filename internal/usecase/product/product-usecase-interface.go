package product

import "bemyfaktur/internal/model"

type ProductUsecaseInterface interface {
	IndexProduct(limit int, offset int, q string) ([]model.Product, error)
	GetProduct(id int) (model.Product, error)
	CreateProduct(request model.Product) (model.ProductRespon, error)
	UpdatedProduct(id int, request model.Product) (model.ProductRespon, error)
	DeleteProduct(id int) (string, error)
}
