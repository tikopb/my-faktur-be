package payment

import "bemyfaktur/internal/model"

type PaymentUsecaseInterface interface {
	Indexpayment(limit int, offset int, q string) ([]model.PaymentRespont, error)
	Getpayment(id int) (model.Payment, error)
	Createpayment(request model.PaymentRequest) (model.PaymentRespont, error)
	Updatedpayment(id int, request model.PaymentRequest) (model.PaymentRespont, error)
	Deletepayment(id int) (string, error)

	IndexLine(limit int, offset int, paymentId int, q string) ([]model.PaymentLineRespont, error)
	GetPaymentLine(id int) (model.PaymentLine, error)
	CreatePaymentLine(request model.PaymentLineRequest) (model.PaymentLineRespont, error)
	UpdatedPaymentLine(id int, request model.PaymentLineRequest) (model.PaymentLineRespont, error)
	DeletePaymentLine(id int) (string, error)

	HandlingPagination(q string, limit int, offset int) (int64, error)
}
