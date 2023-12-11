package invoice

import (
	"bemyfaktur/internal/model"
	"bemyfaktur/internal/model/constant"
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

	//set Partner
	partnerData, err := iu.partnerRepo.ShowInternal(request.PartnerUUID)
	if err != nil || !partnerData.Isactive {
		return data, errors.New("partner not exist")
	}
	request.PartnerId = partnerData.ID

	//setCreatedBy && updateBy
	request.CreatedById = userID
	request.UpdatedById = userID

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
func (iu *invoiceUsecase) IndexInvoice(limit int, offset int, q string, order []string, dateFrom string, dateTo string) ([]model.InvoiceRespont, error) {
	return iu.invoiceRepo.Index(limit, offset, q, order, dateFrom, dateTo)
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
func (iu *invoiceUsecase) UpdatedInvoice(id uuid.UUID, request model.InvoiceRequest, userId string) (model.InvoiceRespont, error) {

	//set Partner
	partnerData, err := iu.partnerRepo.ShowInternal(request.PartnerUUID)
	if err != nil || !partnerData.Isactive {
		return model.InvoiceRespont{}, errors.New("partner not exist or inactive")
	}
	request.PartnerId = partnerData.ID

	//set updated
	request.UpdatedById = userId
	return iu.invoiceRepo.Update(id, request)
}

//-- invoice line part

// CreateInvoiceLine implements InvoiceUsecaseInterface.
func (iu *invoiceUsecase) CreateInvoiceLine(request model.InvoiceLineRequest, userId string) (model.InvoiceLineRespont, error) {
	//set createdby
	request.CreatedById = userId

	//validated the product
	product, err := iu.validateProduct(request.ProductUUID)
	if err != nil {
		return model.InvoiceLineRespont{}, err
	}
	request.ProductID = product.ID //declare the product id of int

	//validated header not void or not complete
	invoice, err := iu.validateInvoice(request.InvoiceUUId)
	if err != nil {
		return model.InvoiceLineRespont{}, err
	}
	request.InvoiceId = invoice.ID

	return iu.invoiceRepo.CreateLine(request)
}

// DeleteInvoiceLine implements InvoiceUsecaseInterface.
func (iu *invoiceUsecase) DeleteInvoiceLine(id uuid.UUID) (string, error) {
	return iu.invoiceRepo.DeleteLine(id)
}

// GetInvoiceLine implements InvoiceUsecaseInterface.
func (iu *invoiceUsecase) GetInvoiceLine(id uuid.UUID) (model.InvoiceLineRespont, error) {
	return iu.invoiceRepo.ShowLine(id)
}

// UpdatedInvoiceLine implements InvoiceUsecaseInterface.
/*
	Definition Of Done
	- Validate Product:
		- Before using a product, ensure that its 'isActive' flag is set to true.
		- If the product is not valid, an expected error function should be triggered.
	- validate invoice:
		- invoice must be active and in draft posisition
	- Calculate Line Amount:
		- This function calculates the line amount based on the quantity and discounts applied.
		- The calculation considers whether the discount is a percentage or a fixed amount.
		- The result is the product of (quantity * unit price) minus the discount.
*/
func (iu *invoiceUsecase) UpdatedInvoiceLine(id uuid.UUID, request model.InvoiceLineRequest) (model.InvoiceLineRespont, error) {
	//validated the product
	//validated the product
	product, err := iu.validateProduct(request.ProductUUID)
	if err != nil {
		return model.InvoiceLineRespont{}, err
	}
	request.ProductID = product.ID //declare the product id of int

	invoice, err := iu.validateInvoice(request.InvoiceUUId)
	if err != nil {
		return model.InvoiceLineRespont{}, err
	}
	request.InvoiceId = invoice.ID

	return iu.invoiceRepo.UpdateLine(id, request)
}

// IndexLine implements InvoiceUsecaseInterface.
func (iu *invoiceUsecase) IndexLine(limit int, offset int, invoiceId uuid.UUID, q string, order []string) ([]model.InvoiceLineRespont, error) {
	//validated header not void or not complete
	invoice, err := iu.invoiceRepo.ShowInternal(invoiceId)
	if err != nil {
		return []model.InvoiceLineRespont{}, err
	}
	return iu.invoiceRepo.IndexLine(limit, offset, invoice.ID, q, order)
}

// validated product is activated on production
func (iu *invoiceUsecase) validateProduct(id uuid.UUID) (model.Product, error) {
	product, err := iu.productRepo.ShowInternal(id)
	if err != nil {
		return model.Product{}, err
	}
	if !product.IsActive {
		return model.Product{}, errors.New("product is not activated, please review the data")
	}
	return product, nil
}

// validatet the invoiced header data
func (iu *invoiceUsecase) validateInvoice(id uuid.UUID) (model.Invoice, error) {
	invoice, err := iu.invoiceRepo.ShowInternal(id)
	if err != nil {
		return model.Invoice{}, err
	}
	if invoice.Status == constant.InvoiceStatusVoid {
		return model.Invoice{}, errors.New("invoice status is void, please review the data")
	} else if invoice.Status == constant.InvoiceStatusProcessed {
		return model.Invoice{}, errors.New("invoice status is in progress, please review the data")
	} else if invoice.Status == constant.InvoiceStatusComplete {
		return model.Invoice{}, errors.New("invoice status is complete, please review the data")
	}

	return invoice, nil
}

func (iu *invoiceUsecase) HandlingPagination(q string, limit int, offset int, dateFrom string, dateTo string) (int64, error) {
	count, err := iu.invoiceRepo.HandlingPagination(q, limit, offset, dateFrom, dateTo)
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
