package payment

import (
	"bemyfaktur/internal/model"
	"bemyfaktur/internal/model/constant"
	"mime/multipart"

	"github.com/google/uuid"
)

type PaymentUsecaseInterface interface {
	Indexpayment(limit int, offset int, q string, order []string, dateFrom string, dateTo string) ([]model.PaymentRespont, error)
	Getpayment(id uuid.UUID) (model.PaymentRespont, error)
	Createpayment(request model.PaymentRequest, userId string) (model.PaymentRespont, error)
	CreatePaymentV2(request model.PaymentRequestV2, userId string) (model.PaymentRespont, error)
	Updatedpayment(id uuid.UUID, request model.PaymentRequest) (model.PaymentRespont, error)
	Deletepayment(id uuid.UUID) (string, error)

	IndexLine(limit int, offset int, paymentId uuid.UUID, q string) ([]model.PaymentLineRespont, int, error)
	GetPaymentLine(id uuid.UUID) (model.PaymentLineRespont, error)
	CreatePaymentLine(request model.PaymentLineRequest, userId string) (model.PaymentLineRespont, error)
	UpdatedPaymentLine(id uuid.UUID, request model.PaymentLineRequest) (model.PaymentLineRespont, error)
	DeletePaymentLine(id uuid.UUID) (string, error)

	HandlingPagination(q string, limit int, offset int, dateFrom string, dateTo string) (int64, error)
	HandlingPaginationLine(q string, limit int, offset int, paymentID int) (int64, error)

	//v3 update function
	PostPaymentV3(request model.PaymentRequestV2, userID string, form *multipart.Form) (model.PaymentRespontV3, error)
	UpdatePaymentV3(id uuid.UUID, request model.PaymentRequest, form *multipart.Form) (model.PaymentRespontV3, error)
	StatusUpdateV3(id uuid.UUID, userId string, docAction constant.PaymentDocAction) (model.PaymentRespont, error)
}
