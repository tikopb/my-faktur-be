package payment

import (
	"bemyfaktur/internal/model"
	"bemyfaktur/internal/model/constant"
	"bemyfaktur/internal/repository/invoice"
	"bemyfaktur/internal/repository/payment"
	"errors"
)

type paymentUsecase struct {
	paymentRepo payment.PaymentRepositoryinterface
	invoiceRepo invoice.InvoiceRepositoryInterface
}

func GetUsecase(paymentRepo payment.PaymentRepositoryinterface, invoiceRepo invoice.InvoiceRepositoryInterface) PaymentUsecaseInterface {
	return &paymentUsecase{
		paymentRepo: paymentRepo,
		invoiceRepo: invoiceRepo,
	}
}

// ehader payment part
// Createpayment implements PaymentUsecaseInterface.
func (pu *paymentUsecase) Createpayment(request model.PaymentRequest, userId string) (model.PaymentRespont, error) {
	request.CreatedBy = userId
	return pu.paymentRepo.Create(request)
}

// Deletepayment implements PaymentUsecaseInterface.
func (pu *paymentUsecase) Deletepayment(id int) (string, error) {
	return pu.paymentRepo.Delete(id)
}

// Getpayment implements PaymentUsecaseInterface.
func (pu *paymentUsecase) Getpayment(id int) (model.Payment, error) {
	return pu.paymentRepo.Show(id)
}

// Indexpayment implements PaymentUsecaseInterface.
func (pu *paymentUsecase) Indexpayment(limit int, offset int, q string) ([]model.PaymentRespont, error) {
	return pu.paymentRepo.Index(limit, offset, q)
}

// Updatedpayment implements PaymentUsecaseInterface.
func (pu *paymentUsecase) Updatedpayment(id int, request model.PaymentRequest) (model.PaymentRespont, error) {
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
func (pu *paymentUsecase) GetPaymentLine(id int) (model.PaymentLine, error) {
	return pu.paymentRepo.ShowLine(id)
}

// IndexLine implements PaymentUsecaseInterface.
func (pu *paymentUsecase) IndexLine(limit int, offset int, paymentId int, q string) ([]model.PaymentLineRespont, error) {
	return pu.paymentRepo.IndexLine(limit, offset, paymentId)
}

// UpdatedInvoiceLine implements PaymentUsecaseInterface.
func (pu *paymentUsecase) UpdatedPaymentLine(id int, request model.PaymentLineRequest) (model.PaymentLineRespont, error) {
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
func (pu *paymentUsecase) DeletePaymentLine(id int) (string, error) {
	return pu.paymentRepo.DeleteLine(id)
}

func (pu *paymentUsecase) HandlingPagination(q string, limit int, offset int) (int64, error) {
	count, err := pu.paymentRepo.HandlingPagination(q, limit, offset)
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
