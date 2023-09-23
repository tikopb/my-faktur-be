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

	documentutil "bemyfaktur/internal/model/documentUtil"
	documentutilRepo "bemyfaktur/internal/model/documentUtil"

	"gorm.io/gorm"
)

type Container struct {
	PartnerUsecase paUsecase.Usecase
	ProductUsecase productUsecase.ProductUsecaseInterface
	InvoiceUsecase invoiceUsecase.InvoiceUsecaseInterface
	PaymentUsecase paymentUsecase.PaymentUsecaseInterface
	DocumentUtil   documentutilRepo.Repository
}

func NewContainer(db *gorm.DB) *Container {
	documnetUtilRepo := documentutil.GetRepository(db)

	partnerRepo := paRepository.GetRepository(db)
	partnerUsecase := paUsecase.GetUsecase(partnerRepo)

	productRepo := productReposiftory.GetRepository(db)
	productUsecase := productUsecase.GetUsecase(productRepo)

	userRepository := usrRepository.GetRepository(db)

	invoiceRepo := invoiceReposiftory.GetRepository(db, documnetUtilRepo)
	invoiceUsecase := invoiceUsecase.GetUsecase(invoiceRepo, partnerRepo, productRepo, userRepository)

	paymentRepo := paymentRepository.GetRepository(db, documnetUtilRepo)
	paymentUsecase := paymentUsecase.GetUsecase(paymentRepo, invoiceRepo)

	return &Container{
		PartnerUsecase: partnerUsecase,
		ProductUsecase: productUsecase,
		InvoiceUsecase: invoiceUsecase,
		PaymentUsecase: paymentUsecase,
	}
}
