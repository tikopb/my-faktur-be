package fakers

import (
	"bemyfaktur/internal/model"
	"time"

	"gorm.io/gorm"
)

func PaymentFaker(db *gorm.DB, payment model.Payment) *model.Payment {
	return &model.Payment{
		ID:             payment.ID,
		CreatedAt:      time.Time{},
		CreatedBy:      payment.CreatedBy,
		UpdatedBy:      payment.UpdatedBy,
		PartnerID:      payment.PartnerID,
		GrandTotal:     payment.GrandTotal,
		Discount:       payment.Discount,
		BatchNo:        payment.BatchNo,
		Status:         payment.Status,
		DocAction:      payment.DocAction,
		DocumentNo:     payment.DocumentNo,
		IsPrecentage:   payment.IsPrecentage,
		UUID:           payment.UUID,
		OrganizationId: payment.OrganizationId,
	}
}

func PaymentLineFaker(db *gorm.DB, line model.PaymentLine) *model.PaymentLine {
	return &model.PaymentLine{
		ID:             line.ID,
		PaymentID:      line.PaymentID,
		Price:          line.Price,
		Amount:         line.Amount,
		CreatedAt:      time.Time{},
		CreatedBy:      line.CreatedBy,
		UpdatedBy:      line.UpdatedBy,
		UserUpdated:    line.UserUpdated,
		InvoiceID:      line.InvoiceID,
		Discount:       line.Discount,
		IsPrecentage:   line.IsPrecentage,
		UUID:           line.UUID,
		OrganizationId: line.OrganizationId,
	}
}
