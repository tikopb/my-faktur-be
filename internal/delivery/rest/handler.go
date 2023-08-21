package rest

import (
	"bemyfaktur/internal/usecase/partner"
)

type handler struct {
	partnerUsecase partner.Usecase
}

func NewHandler(partnerUsecase partner.Usecase) *handler {
	return &handler{
		partnerUsecase: partnerUsecase,
	}
}
