package product

import (
	"bemyfaktur/internal/model"
	"bemyfaktur/internal/repository/product"

	"github.com/google/uuid"
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
func (m *productUsecase) IndexProduct(limit int, offset int, q string, order []string) ([]model.ProductRespon, error) {
	return m.productRepo.Index(limit, offset, q, order)
}

// CreateProduct implements Usecase.
func (m *productUsecase) CreateProduct(request model.Product, userId string) (model.ProductRespon, error) {
	request.CreatedBy = userId
	request.OrganizationId = 1
	return m.productRepo.Create(request)
}

// DeleteProduct implements Usecase.
func (m *productUsecase) DeleteProduct(id uuid.UUID) (string, error) {
	return m.productRepo.Delete(id)
}

// GetProduct implements Usecase.
func (m *productUsecase) GetProduct(id uuid.UUID) (model.ProductRespon, error) {
	return m.productRepo.Show(id)
}

// UpdatedProduct implements Usecase.
func (m *productUsecase) UpdatedProduct(id uuid.UUID, request model.Product) (model.ProductRespon, error) {
	return m.productRepo.Update(id, request)
}

// Partial implements ProductUsecaseInterface.
func (m *productUsecase) Partial(q string) ([]model.ProductPartialRespon, error) {
	return m.productRepo.Partial(q)
}
