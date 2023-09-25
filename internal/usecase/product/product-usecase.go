package product

import (
	"bemyfaktur/internal/model"
	"bemyfaktur/internal/repository/product"
)

type productUsecase struct {
	productRepo product.Repository
}

func GetUsecase(productRepo product.Repository) ProductUsecaseInterface {
	return &productUsecase{
		productRepo: productRepo,
	}
}

// IndexPartner implements Usecase.
func (m *productUsecase) IndexProduct(limit int, offset int, q string) ([]model.Product, error) {
	return m.productRepo.Index(limit, offset, q)
}

// CreateProduct implements Usecase.
func (m *productUsecase) CreateProduct(request model.Product) (model.ProductRespon, error) {
	return m.productRepo.Create(request)
}

// DeleteProduct implements Usecase.
func (m *productUsecase) DeleteProduct(id int) (string, error) {
	return m.productRepo.Delete(id)
}

// GetProduct implements Usecase.
func (m *productUsecase) GetProduct(id int) (model.Product, error) {
	return m.productRepo.Show(id)
}

// UpdatedProduct implements Usecase.
func (m *productUsecase) UpdatedProduct(id int, request model.Product) (model.ProductRespon, error) {
	return m.productRepo.Update(id, request)
}
