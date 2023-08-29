package invoice

import (
	"bemyfaktur/internal/model"
	"errors"

	"gorm.io/gorm"
)

// CreateLine implements InvoiceRepositoryInterface.
func (ir *invoiceRepo) CreateLine(invoice model.InvoiceLine) (model.InvoiceLine, error) {
	data := model.InvoiceLine{}

	if err := ir.db.Create(&data).Error; err != nil {
		return data, err
	}

	return data, nil
}

// IndexLine implements InvoiceRepositoryInterface.
func (ir *invoiceRepo) IndexLine(limit int, offset int, invoiceId int) ([]model.InvoiceLine, error) {
	data := []model.InvoiceLine{}

	if err := ir.db.Order("CreatedAt DESC").Where(model.InvoiceLine{InvoiceID: invoiceId}).Limit(limit).Offset(offset).Find(&data).Error; err != nil {
		return data, err
	}

	return data, nil
}

// ShowLine implements InvoiceRepositoryInterface.
func (ir *invoiceRepo) ShowLine(id int) (model.InvoiceLine, error) {
	data := model.InvoiceLine{}

	if err := ir.db.Where(model.InvoiceLine{ID: id}).Preload("InvoiceLine").First(&data).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return data, errors.New("data not found")
		}
	}
	return data, nil
}

// UpdateLine implements InvoiceRepositoryInterface.
func (ir *invoiceRepo) UpdateLine(id int, updatedInvoiceLine model.InvoiceLine) (model.InvoiceLine, error) {
	data, err := ir.ShowLine(id)

	if err != nil {
		return data, nil
	}

	data.ProductID = updatedInvoiceLine.ProductID
	data.Price = updatedInvoiceLine.Price
	data.Discount = updatedInvoiceLine.Discount
	data.Qty = updatedInvoiceLine.Qty
	data.Amount = updatedInvoiceLine.Amount

	if err := ir.db.Save(&data).Error; err != nil {
		return data, err
	}

	return data, nil
}

// DeleteLine implements InvoiceRepositoryInterface.
func (ir *invoiceRepo) DeleteLine(id int) (string, error) {
	data, err := ir.ShowLine(id)
	if err != nil {
		return "", nil
	}

	if err := ir.db.Delete(&data).Error; err != nil {
		return "", err
	}

	return "data deleted", nil
}
