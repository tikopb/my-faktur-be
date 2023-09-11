package payment

import (
	"bemyfaktur/internal/model"
)

type PaymentRepositoryinterface interface {
	Index(limit int, offset int) ([]model.PaymentRespont, error)
	Create(payment model.PaymentRequest) (model.PaymentRespont, error)
	Show(id int) (model.Payment, error)
	Update(id int, updatedPayment model.PaymentRequest) (model.PaymentRespont, error)
	Delete(id int) (string, error)

	IndexLine(limit int, offset int, paymentId int) ([]model.PaymentLineRespont, error)
	CreateLine(paymentLine model.PaymentLineRequest) (model.PaymentLineRespont, error)
	ShowLine(id int) (model.PaymentLine, error)
	UpdateLine(id int, updatedPaymentLine model.PaymentLineRequest) (model.PaymentLineRespont, error)
	DeleteLine(id int) (string, error)

	//docValidation
	DocProcess(data model.Payment) error
	CompleteIT(data model.Payment) error
	ReversedIt(data model.Payment) error
}
