package payment

import (
	"bemyfaktur/internal/model"
	"bemyfaktur/internal/model/constant"
	"bemyfaktur/internal/repository/invoice"
	"bemyfaktur/internal/repository/partner"
	"bemyfaktur/internal/repository/payment"
	"errors"

	"github.com/google/uuid"
)

type paymentUsecase struct {
	paymentRepo payment.PaymentRepositoryinterface
	invoiceRepo invoice.InvoiceRepositoryInterface
	partnerRepo partner.Repository
}

func GetUsecase(paymentRepo payment.PaymentRepositoryinterface, invoiceRepo invoice.InvoiceRepositoryInterface, partnerRepo partner.Repository) PaymentUsecaseInterface {
	return &paymentUsecase{
		paymentRepo: paymentRepo,
		invoiceRepo: invoiceRepo,
		partnerRepo: partnerRepo,
	}
}

// ehader payment part
// Createpayment implements PaymentUsecaseInterface.
func (pu *paymentUsecase) Createpayment(request model.PaymentRequest, userId string) (model.PaymentRespont, error) {
	request.CreatedBy = userId
	request.UpdatedBy = userId

	//validate the partner
	partner, err := pu.partnerRepo.ShowInternal(request.PartnerUUID)
	if err != nil || !partner.Isactive {
		return model.PaymentRespont{}, errors.New("partner not exist or inactived")
	}

	request.PartnerID = partner.ID

	return pu.paymentRepo.Create(request)
}

// Deletepayment implements PaymentUsecaseInterface.
func (pu *paymentUsecase) Deletepayment(id uuid.UUID) (string, error) {
	return pu.paymentRepo.Delete(id)
}

// Getpayment implements PaymentUsecaseInterface.
func (pu *paymentUsecase) Getpayment(id uuid.UUID) (model.PaymentRespont, error) {
	return pu.paymentRepo.Show(id)
}

// Indexpayment implements PaymentUsecaseInterface.
func (pu *paymentUsecase) Indexpayment(limit int, offset int, q string, order []string, dateFrom string, dateTo string) ([]model.PaymentRespont, error) {
	return pu.paymentRepo.Index(limit, offset, q, order, dateFrom, dateTo)
}

// Updatedpayment implements PaymentUsecaseInterface.
func (pu *paymentUsecase) Updatedpayment(id uuid.UUID, request model.PaymentRequest) (model.PaymentRespont, error) {

	return pu.paymentRepo.Update(id, request)
}

// invoice line part
// CreateInvoiceLine implements PaymentUsecaseInterface.
func (pu *paymentUsecase) CreatePaymentLine(request model.PaymentLineRequest, userId string) (model.PaymentLineRespont, error) {
	data := model.PaymentLineRespont{}
	invoice, err := pu.invoiceRepo.ShowInternal(request.Invoice_uuid)
	if err != nil {
		//set the value to invoice_id because relation key used with id int not uuid
		request.Invoice_id = invoice.ID
		return data, err
	} else if invoice.Status != constant.InvoiceStatusComplete {
		return data, errors.New("invoice not in completed")
	}

	request.CreatedBy = userId

	//return value
	return pu.paymentRepo.CreateLine(request)
}

// GetInvoiceLine implements PaymentUsecaseInterface.
func (pu *paymentUsecase) GetPaymentLine(id uuid.UUID) (model.PaymentLineRespont, error) {
	return pu.paymentRepo.ShowLine(id)
}

// IndexLine implements PaymentUsecaseInterface.
func (pu *paymentUsecase) IndexLine(limit int, offset int, paymentId int, q string) ([]model.PaymentLineRespont, error) {
	return pu.paymentRepo.IndexLine(limit, offset, paymentId)
}

// UpdatedInvoiceLine implements PaymentUsecaseInterface.
func (pu *paymentUsecase) UpdatedPaymentLine(id uuid.UUID, request model.PaymentLineRequest) (model.PaymentLineRespont, error) {
	data := model.PaymentLineRespont{}
	invoice, err := pu.invoiceRepo.Show(request.Invoice_uuid)
	if err != nil {
		return data, err
	} else if invoice.Status != constant.InvoiceStatusComplete {
		return data, errors.New("invoice not in completed")
	}

	//return value
	return pu.paymentRepo.UpdateLine(id, request)
}

// DeleteInvoiceLine implements PaymentUsecaseInterface.
func (pu *paymentUsecase) DeletePaymentLine(id uuid.UUID) (string, error) {
	return pu.paymentRepo.DeleteLine(id)
}

// CreateV2 implements PaymentUsecaseInterface.
func (pu *paymentUsecase) CreateV2(request model.PaymentRequestV2, userId string) (model.PaymentRespontV2, error) {
	//header validation
	request.Header.CreatedBy = userId

	//line validation
	for _, line := range request.Line {
		//validate invoice
		invoice, err := pu.invoiceRepo.ShowInternal(line.Invoice_uuid)
		if err != nil {
			//set the value to invoice_id because relation key used with id int not uuid
			line.Invoice_id = invoice.ID
		} else if invoice.Status != constant.InvoiceStatusComplete {
			return model.PaymentRespontV2{}, errors.New("invoice not in completed")
		}

		//set created by
		line.CreatedBy = userId
		line.UpdatedBy = userId
	}

	//hit the request
	data, err := pu.paymentRepo.CreateV2(request)
	if err != nil {
		return model.PaymentRespontV2{}, err
	}

	//set return value of data
	return data, nil
}

func (pu *paymentUsecase) HandlingPagination(q string, limit int, offset int, dateFrom string, dateTo string) (int64, error) {
	count, err := pu.paymentRepo.HandlingPagination(q, limit, offset, dateFrom, dateTo)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (pu *paymentUsecase) HandlingPaginationLine(q string, limit int, offset int, paymentID int) (int64, error) {
	count, err := pu.paymentRepo.HandlingPaginationLine(q, limit, offset, paymentID)
	if err != nil {
		return 0, err
	}
	return count, nil
}
