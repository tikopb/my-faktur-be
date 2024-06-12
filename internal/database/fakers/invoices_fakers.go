package fakers

import (
	"bemyfaktur/internal/model"
	"time"

	"gorm.io/gorm"
)

func InvoiceFaker(db *gorm.DB, invoice model.Invoice) *model.Invoice {
	return &model.Invoice{
		UUID:              invoice.UUID,
		CreatedAt:         time.Time{},
		UpdateAt:          time.Time{},
		CreatedBy:         invoice.CreatedBy,
		UpdatedBy:         invoice.UpdatedBy,
		PartnerID:         invoice.PartnerID,
		GrandTotal:        invoice.GrandTotal,
		Discount:          invoice.Discount,
		BatchNo:           invoice.BatchNo,
		Status:            invoice.Status,
		DocAction:         invoice.DocAction,
		OustandingPayment: invoice.OustandingPayment,
		DocumentNo:        invoice.DocumentNo,
		IsPrecentage:      invoice.IsPrecentage,
		PayDate:           invoice.PayDate,
		OrganizationId:    invoice.OrganizationId,
	}
}

func InvoiceLineFaker(db *gorm.DB, line model.InvoiceLine) *model.InvoiceLine {
	return &model.InvoiceLine{
		CreatedAt:      time.Time{},
		UpdateAt:       time.Time{},
		CreatedBy:      line.CreatedBy,
		UpdatedBy:      line.UpdatedBy,
		Price:          line.Price,
		Discount:       line.Discount,
		Qty:            line.Qty,
		Amount:         line.Amount,
		IsPrecentage:   line.IsPrecentage,
		ProductID:      line.ProductID,
		InvoiceID:      line.InvoiceID,
		UUID:           line.UUID,
		OrganizationId: line.OrganizationId,
	}
}
