// container.go
package usecase

import (
	paRepository "bemyfaktur/internal/repository/partner"
	paUsecase "bemyfaktur/internal/usecase/partner"

	productReposiftory "bemyfaktur/internal/repository/product"
	productUsecase "bemyfaktur/internal/usecase/product"

	invoiceReposiftory "bemyfaktur/internal/repository/invoice"
	invoiceUsecase "bemyfaktur/internal/usecase/invoice"

	"gorm.io/gorm"
)

type Container struct {
	PartnerUsecase paUsecase.Usecase
	ProductUsecase productUsecase.ProductUsecaseInterface
	InvoiceUsecase invoiceUsecase.InvoiceUsecaseInterface
}

func NewContainer(db *gorm.DB) *Container {
	partnerRepo := paRepository.GetRepository(db)
	partnerUsecase := paUsecase.GetUsecase(partnerRepo)

	productRepo := productReposiftory.GetRepository(db)
	productUsecase := productUsecase.GetUsecase(productRepo)

	invoiceRepo := invoiceReposiftory.GetRepository(db)
	invoiceUsecase := invoiceUsecase.GetUsecase(invoiceRepo)

	return &Container{
		PartnerUsecase: partnerUsecase,
		ProductUsecase: productUsecase,
		InvoiceUsecase: invoiceUsecase,
	}
}