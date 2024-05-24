// container.go
package usecase

import (
	paRepository "bemyfaktur/internal/repository/partner"
	"bemyfaktur/internal/usecase/fileservice"

	paUsecase "bemyfaktur/internal/usecase/partner"
	"time"

	productReposiftory "bemyfaktur/internal/repository/product"
	productUsecase "bemyfaktur/internal/usecase/product"

	invoiceReposiftory "bemyfaktur/internal/repository/invoice"
	invoiceUsecase "bemyfaktur/internal/usecase/invoice"

	paymentRepository "bemyfaktur/internal/repository/payment"
	paymentUsecase "bemyfaktur/internal/usecase/payment"

	documentutil "bemyfaktur/internal/model/documentUtil"
	pgUtil "bemyfaktur/internal/model/paginationUtil"

	userRepo "bemyfaktur/internal/repository/user"
	authUsecase "bemyfaktur/internal/usecase/auth"

	orgRepo "bemyfaktur/internal/repository/organization"
	orgUsecase "bemyfaktur/internal/usecase/organization"

	fileserviceRepo "bemyfaktur/internal/repository/fileService"

	midUtil "bemyfaktur/internal/delivery/auth"

	"crypto/rand"
	"crypto/rsa"

	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type Container struct {
	PartnerUsecase      paUsecase.Usecase
	ProductUsecase      productUsecase.ProductUsecaseInterface
	InvoiceUsecase      invoiceUsecase.InvoiceUsecaseInterface
	PaymentUsecase      paymentUsecase.PaymentUsecaseInterface
	FileserviceUsecase  fileservice.Repository
	DocumentUtil        documentutil.Repository
	AuthUsecase         authUsecase.Usecase
	PgUtil              pgUtil.Repository
	Middleware          midUtil.MidlewareInterface
	OrganizationUsecase orgUsecase.Usecase
}

func NewContainer(db *gorm.DB) *Container {

	organizationRepo := orgRepo.GetRepository(db)
	organizationUsecase := orgUsecase.GetUsecase(organizationRepo)

	secret := GetEnv("key_secret")
	signKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		panic(err)
	}
	userRepo, err := userRepo.GetRepository(db, secret, 1, 64*1024, 4, 32, signKey, 60*time.Second, 48*time.Hour, organizationUsecase)
	if err != nil {
		panic("errorr repo")
	}
	authUsecase := authUsecase.GetUsecase(userRepo)

	middleware := midUtil.GetAuthMiddleware(authUsecase)

	fileserviceRepo := fileserviceRepo.GetRepository(db)
	fileserviceUsecase := fileservice.GetRepository(fileserviceRepo)

	documentUtilRepo := documentutil.GetRepository(db)
	pgUtilRepo := pgUtil.GetRepository(db)

	partnerRepo := paRepository.GetRepository(db, pgUtilRepo)
	partnerUsecase := paUsecase.GetUsecase(partnerRepo)

	productRepo := productReposiftory.GetRepository(db, pgUtilRepo)
	productUsecase := productUsecase.GetUsecase(productRepo)

	invoiceRepo := invoiceReposiftory.GetRepository(db, documentUtilRepo, pgUtilRepo)
	invoiceUsecase := invoiceUsecase.GetUsecase(invoiceRepo, partnerRepo, productRepo, fileserviceUsecase)

	paymentRepo := paymentRepository.GetRepository(db, documentUtilRepo, pgUtilRepo)
	paymentUsecase := paymentUsecase.GetUsecase(paymentRepo, invoiceRepo, partnerRepo)

	return &Container{
		PartnerUsecase:      partnerUsecase,
		ProductUsecase:      productUsecase,
		InvoiceUsecase:      invoiceUsecase,
		PaymentUsecase:      paymentUsecase,
		DocumentUtil:        documentUtilRepo,
		AuthUsecase:         authUsecase,
		PgUtil:              pgUtilRepo,
		Middleware:          middleware,
		FileserviceUsecase:  fileserviceUsecase,
		OrganizationUsecase: organizationUsecase,
	}
}

func GetEnv(param string) string {
	var value string

	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()
	if err != nil {
		panic("config environment not found!")
	}

	value = viper.GetString("key_secret")
	return value
}
