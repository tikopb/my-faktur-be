package payment

import (
	"bemyfaktur/internal/model"

	"github.com/google/uuid"
)

type PaymentRepositoryinterface interface {
	Index(limit int, offset int, q string, order []string, dateFrom string, dateTo string) ([]model.PaymentRespont, error)
	Create(payment model.PaymentRequest) (model.PaymentRespont, error)
	CreateV2(payment model.PaymentRequestV2) (model.PaymentRespontV2, error)
	Show(id uuid.UUID) (model.PaymentRespont, error)
	ShowInternal(id uuid.UUID) (model.Payment, error)
	Update(id uuid.UUID, updatedPayment model.PaymentRequest) (model.PaymentRespont, error)
	Delete(id uuid.UUID) (string, error)

	IndexLine(limit int, offset int, paymentId int) ([]model.PaymentLineRespont, error)
	CreateLine(paymentLine model.PaymentLineRequest) (model.PaymentLineRespont, error)
	ShowLine(id uuid.UUID) (model.PaymentLineRespont, error)
	ShowLineInternal(id uuid.UUID) (model.PaymentLine, error)
	UpdateLine(id uuid.UUID, updatedPaymentLine model.PaymentLineRequest) (model.PaymentLineRespont, error)
	DeleteLine(id uuid.UUID) (string, error)

	//save validaition
	BeforeSave(data model.Payment) (model.Payment, error)

	//docValidation
	DocProcess(data model.Payment, docaction string) (model.Payment, error)
	CompleteIT(data model.Payment, docaction string) (model.Payment, error)
	ReversedIt(data model.Payment, docaction string) (model.Payment, error)

	//pagination handling
	HandlingPagination(q string, limit int, offset int, dateFrom string, dateTo string) (int64, error)
	HandlingPaginationLine(q string, limit int, offset int, paymentID int) (int64, error)
}
