package invoice

import (
	"bemyfaktur/internal/model"
	"bemyfaktur/internal/repository/invoice"
	"bemyfaktur/internal/repository/partner"
	"bemyfaktur/internal/repository/product"
	"errors"

	"github.com/google/uuid"
)

type invoiceUsecase struct {
	invoiceRepo invoice.InvoiceRepositoryInterface
	partnerRepo partner.Repository
	productRepo product.Repository
	//userRepo    user.Repository
}

func GetUsecase(invoiceRepo invoice.InvoiceRepositoryInterface, partnerRepo partner.Repository, productRepo product.Repository) InvoiceUsecaseInterface {
	return &invoiceUsecase{
		invoiceRepo: invoiceRepo,
		partnerRepo: partnerRepo,
		productRepo: productRepo,
	}
}

// CreateInvoice implements InvoiceUsecaseInterface.
func (iu *invoiceUsecase) CreateInvoice(request model.InvoiceRequest, userID string) (model.InvoiceRespont, error) {
	//# Get all data First
	data := model.InvoiceRespont{}

	//get Partner
	partnerData, err := iu.partnerRepo.ShowInternal(request.PartnerIdUU)
	if err != nil || !partnerData.Isactive {
		return data, errors.New("partner not exist")
	}

	//setCreatedBy
	request.CreatedBy = userID
	return iu.invoiceRepo.Create(request, partnerData)
}

// DeleteInvoice implements InvoiceUsecaseInterface.
func (iu *invoiceUsecase) DeleteInvoice(id uuid.UUID) (string, error) {
	return iu.invoiceRepo.Delete(id)
}

// GetInvoice implements InvoiceUsecaseInterface.
func (iu *invoiceUsecase) GetInvoice(id uuid.UUID) (model.InvoiceRespont, error) {
	return iu.invoiceRepo.Show(id)
}

// IndexInvoice implements InvoiceUsecaseInterface.
func (iu *invoiceUsecase) IndexInvoice(limit int, offset int, q string, order []string) ([]model.InvoiceRespont, error) {
	return iu.invoiceRepo.Index(limit, offset, q, order)
}

// UpdatedInvoice implements InvoiceUsecaseInterface.
/**
DOD (Definition Of Done)
1. 	validated partner
1.a validated for partner first
1.b if partner exist and isactive == true then passed
2.	insert data of invoice with struct invoice
3. 	fill invoice respon for return
*/
func (iu *invoiceUsecase) UpdatedInvoice(id uuid.UUID, request model.Invoice) (model.InvoiceRespont, error) {
	data := model.InvoiceRespont{}

	//get and set partner
	partnerData, err := iu.partnerRepo.Show(request.PartnerIdUU)
	if err != nil || !partnerData.Isactive {
		return data, errors.New("partner not exist or partnet not active")
	}

	return iu.invoiceRepo.Update(id, request)
}

//-- invoice line part

// CreateInvoiceLine implements InvoiceUsecaseInterface.
func (iu *invoiceUsecase) CreateInvoiceLine(request model.InvoiceLine, userId string) (model.InvoiceLine, error) {
	//set createdby
	request.CreatedBy = userId

	// data := model.InvoiceLine{}
	// if !iu.validateProductIsActive(request.Product.UUID) {
	// 	return data, errors.New("product is not activated")
	// }
	return iu.invoiceRepo.CreateLine(request)
}

// DeleteInvoiceLine implements InvoiceUsecaseInterface.
func (iu *invoiceUsecase) DeleteInvoiceLine(id int) (string, error) {
	return iu.invoiceRepo.DeleteLine(id)
}

// GetInvoiceLine implements InvoiceUsecaseInterface.
func (iu *invoiceUsecase) GetInvoiceLine(id int) (model.InvoiceLine, error) {
	return iu.invoiceRepo.ShowLine(id)
}

// UpdatedInvoiceLine implements InvoiceUsecaseInterface.
/*
	Definition Of Done
	- Validate Product:
		- Before using a product, ensure that its 'isActive' flag is set to true.
		- If the product is not valid, an expected error function should be triggered.

	- Calculate Line Amount:
		- This function calculates the line amount based on the quantity and discounts applied.
		- The calculation considers whether the discount is a percentage or a fixed amount.
		- The result is the product of (quantity * unit price) minus the discount.
*/
func (iu *invoiceUsecase) UpdatedInvoiceLine(id int, request model.InvoiceLine, productId uuid.UUID) (model.InvoiceLine, error) {
	data := model.InvoiceLine{}
	//validate product is active!
	if !iu.validateProductIsActive(productId) {
		return data, errors.New("product not activated")
	}

	return iu.invoiceRepo.UpdateLine(id, request)
}

// IndexLine implements InvoiceUsecaseInterface.
func (iu *invoiceUsecase) IndexLine(limit int, offset int, invoiceId int, q string) ([]model.InvoiceLineRespont, error) {
	return iu.invoiceRepo.IndexLine(limit, offset, invoiceId, q)
}

// validated product is activated on production
func (iu *invoiceUsecase) validateProductIsActive(id uuid.UUID) bool {
	//validate product is active!
	data, err := iu.productRepo.Show(id)
	if err != nil {
		return false
	}

	if data.IsActive {
		return true
	}
	return false
}

func (iu *invoiceUsecase) HandlingPagination(q string, limit int, offset int) (int64, error) {
	count, err := iu.invoiceRepo.HandlingPagination(q, limit, offset)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (iu *invoiceUsecase) HandlingPaginationLine(q string, limit int, offset int, invoiceId int) (int64, error) {
	count, err := iu.invoiceRepo.HandlingPaginationLine(q, limit, offset, invoiceId)
	if err != nil {
		return 0, err
	}
	return count, nil
}
