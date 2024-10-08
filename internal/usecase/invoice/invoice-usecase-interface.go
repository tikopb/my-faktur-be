package invoice

import (
	"bemyfaktur/internal/model"
	"bemyfaktur/internal/model/constant"
	"mime/multipart"

	"github.com/google/uuid"
)

type InvoiceUsecaseInterface interface {
	IndexInvoice(limit int, offset int, q string, order []string, dateFrom string, dateTo string) ([]model.InvoiceRespont, error)
	GetInvoice(id uuid.UUID) (model.InvoiceRespont, error)
	CreateInvoice(request model.InvoiceRequest, userId string) (model.InvoiceRespont, error)
	UpdatedInvoice(id uuid.UUID, request model.InvoiceRequest, userId string) (model.InvoiceRespont, error)
	DeleteInvoice(id uuid.UUID) (string, error)
	Partial(partner_id uuid.UUID, q string) ([]model.InvoicePartialRespont, error)

	IndexLine(limit int, offset int, invoiceId uuid.UUID, q string, order []string) ([]model.InvoiceLineRespont, error)
	GetInvoiceLine(id uuid.UUID) (model.InvoiceLineRespont, error)
	CreateInvoiceLine(request model.InvoiceLineRequest, userId string) (model.InvoiceLineRespont, error)
	UpdatedInvoiceLine(id uuid.UUID, request model.InvoiceLineRequest) (model.InvoiceLineRespont, error)
	DeleteInvoiceLine(id uuid.UUID) (string, error)

	HandlingPagination(q string, limit int, offset int, dateFrom string, dateTo string) (int64, error)
	HandlingPaginationLine(q string, limit int, offset int, invoiceId int) (int64, error)

	//v2
	CreateInvoiceV2(request model.InvoiceRequestV2, userId string) (model.InvoiceRespontV2, error)

	//v3
	UpdateInvoiceV3(id uuid.UUID, request model.InvoiceRequest, form multipart.Form) (model.InvoiceRespontV3, error)
	StatusUpdateV3(id uuid.UUID, userId string, docAction constant.InvoiceDocAction) (model.InvoiceRespont, error)
}
