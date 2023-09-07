// container.go
package usecase

import (
	paRepository "bemyfaktur/internal/repository/partner"
	paUsecase "bemyfaktur/internal/usecase/partner"

	productReposiftory "bemyfaktur/internal/repository/product"
	productUsecase "bemyfaktur/internal/usecase/product"

	invoiceReposiftory "bemyfaktur/internal/repository/invoice"
	invoiceUsecase "bemyfaktur/internal/usecase/invoice"

	usrRepository "bemyfaktur/internal/repository/user"

	paymentRepository "bemyfaktur/internal/repository/payment"
	paymentUsecase "bemyfaktur/internal/usecase/payment"

	"gorm.io/gorm"
)

type Container struct {
	PartnerUsecase paUsecase.Usecase
	ProductUsecase productUsecase.ProductUsecaseInterface
	InvoiceUsecase invoiceUsecase.InvoiceUsecaseInterface
	PaymentUsecase paymentUsecase.PaymentUsecaseInterface
}

func NewContainer(db *gorm.DB) *Container {
	partnerRepo := paRepository.GetRepository(db)
	partnerUsecase := paUsecase.GetUsecase(partnerRepo)

	productRepo := productReposiftory.GetRepository(db)
	productUsecase := productUsecase.GetUsecase(productRepo)

	userRepository := usrRepository.GetRepository(db)

	invoiceRepo := invoiceReposiftory.GetRepository(db)
	invoiceUsecase := invoiceUsecase.GetUsecase(invoiceRepo, partnerRepo, productRepo, userRepository)

	paymentRepo := paymentRepository.GetRepository(db)
	paymentUsecase := paymentUsecase.GetUsecase(paymentRepo, invoiceRepo)

	return &Container{
		PartnerUsecase: partnerUsecase,
		ProductUsecase: productUsecase,
		InvoiceUsecase: invoiceUsecase,
		PaymentUsecase: paymentUsecase,
	}
}
