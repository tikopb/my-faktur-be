package fakers

import (
	"bemyfaktur/internal/model"
	"time"

	"gorm.io/gorm"
)

func PaymentFaker(db *gorm.DB, invoice model.Payment) *model.Payment {
	return &model.Payment{
		ID:           invoice.ID,
		CreatedAt:    time.Time{},
		CreatedBy:    invoice.CreatedBy,
		UpdatedBy:    invoice.UpdatedBy,
		PartnerID:    invoice.PartnerID,
		GrandTotal:   invoice.GrandTotal,
		Discount:     invoice.Discount,
		BatchNo:      invoice.BatchNo,
		Status:       invoice.Status,
		DocAction:    invoice.DocAction,
		DocumentNo:   invoice.DocumentNo,
		IsPrecentage: invoice.IsPrecentage,
	}
}

func PaymentLineFaker(db *gorm.DB, line model.PaymentLine) *model.PaymentLine {
	return &model.PaymentLine{
		ID:           line.ID,
		PaymentID:    line.PaymentID,
		Price:        line.Price,
		Amount:       line.Amount,
		CreatedAt:    time.Time{},
		CreatedBy:    line.CreatedBy,
		UpdatedBy:    line.UpdatedBy,
		UserUpdated:  line.UserUpdated,
		InvoiceID:    line.InvoiceID,
		Discount:     line.Discount,
		IsPrecentage: line.IsPrecentage,
	}
}
