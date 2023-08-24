package invoice

import (
	"bemyfaktur/internal/model"
	"bemyfaktur/internal/repository/invoice"
	"bemyfaktur/internal/repository/partner"
	"errors"
)

type invoiceUsecase struct {
	invoiceRepo invoice.InvoiceRepositoryInterface
	parterRepo  partner.Repository
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
