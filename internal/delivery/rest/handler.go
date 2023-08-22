package rest

import (
	"bemyfaktur/internal/usecase/partner"
	"bemyfaktur/internal/usecase/product"
)

type handler struct {
	partnerUsecase partner.Usecase
	productUsecase product.ProductUsecaseInterface
}

func NewHandler(partnerUsecase partner.Usecase, productUsecase product.ProductUsecaseInterface) *handler {

	return &handler{
		partnerUsecase: partnerUsecase,
		productUsecase: productUsecase,
	}
}
