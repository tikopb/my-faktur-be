package invoice

import (
	"bemyfaktur/internal/model"
	"bemyfaktur/internal/repository/invoice"
	"bemyfaktur/internal/repository/partner"
	"bemyfaktur/internal/repository/product"
	"errors"
)

type invoiceUsecase struct {
	invoiceRepo invoice.InvoiceRepositoryInterface
	parterRepo  partner.Repository
	productRepo product.Repository
}

func GetUsecase(invoiceRepo invoice.InvoiceRepositoryInterface) InvoiceUsecaseInterface {
	return &invoiceUsecase{
		invoiceRepo: invoiceRepo,
	}
}

// CreateInvoice implements InvoiceUsecaseInterface.
func (iu *invoiceUsecase) CreateInvoice(request model.Invoice) (model.InvoiceCreateRespon, error) {
	data := model.InvoiceCreateRespon{}
	partnerData, err := iu.parterRepo.Show(request.ID)
	if err != nil || !partnerData.Isactive {
		return data, errors.New("partner not exist")
	}

	return iu.invoiceRepo.Create(request)
}

// DeleteInvoice implements InvoiceUsecaseInterface.
func (iu *invoiceUsecase) DeleteInvoice(id int) (string, error) {
	return iu.invoiceRepo.Delete(id)
}

// GetInvoice implements InvoiceUsecaseInterface.
func (iu *invoiceUsecase) GetInvoice(id int) (model.Invoice, error) {
	return iu.invoiceRepo.Show(id)
}

// IndexInvoice implements InvoiceUsecaseInterface.
func (iu *invoiceUsecase) IndexInvoice(limit int, offset int) ([]model.InvoiceIndexRespont, error) {
	return iu.invoiceRepo.Index(limit, offset)
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
func (iu *invoiceUsecase) UpdatedInvoice(id int, request model.Invoice) (model.Invoice, error) {
	data := model.Invoice{}
	partnerData, err := iu.parterRepo.Show(request.ID)
	if err != nil || !partnerData.Isactive {
		return data, errors.New("partner not exist")
	}

	return iu.invoiceRepo.Update(id, request)
}

//-- invoice line part

// CreateInvoiceLine implements InvoiceUsecaseInterface.
func (iu *invoiceUsecase) CreateInvoiceLine(request model.InvoiceLine) (model.InvoiceLine, error) {
	data := model.InvoiceLine{}
	if !iu.validateProductIsActive(request.ProductID) {
		return data, errors.New("product is not activated")
	}
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
func (iu *invoiceUsecase) UpdatedInvoiceLine(id int, request model.InvoiceLine, productId int) (model.InvoiceLine, error) {
	data := model.InvoiceLine{}
	//validate product is active!
	if !iu.validateProductIsActive(productId) {
		return data, errors.New("product not activated")
	}

	return iu.invoiceRepo.UpdateLine(id, request)
}

// IndexLine implements InvoiceUsecaseInterface.
func (iu *invoiceUsecase) IndexLine(limit int, offset int, invoiceId int) ([]model.InvoiceLine, error) {
	return iu.invoiceRepo.IndexLine(limit, offset, invoiceId)
}

// validated product is activated on production
func (iu *invoiceUsecase) validateProductIsActive(id int) bool {
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
