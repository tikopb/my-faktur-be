package invoice

import (
	"bemyfaktur/internal/model"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

// CreateLine implements InvoiceRepositoryInterface.
func (ir *invoiceRepo) CreateLine(request model.InvoiceLine) (model.InvoiceLine, error) {
	data := model.InvoiceLine{
		InvoiceID:    request.InvoiceID,
		ProductID:    request.ProductID,
		Price:        request.Price,
		Discount:     request.Discount,
		Qty:          request.Qty,
		Amount:       ir.handlingAmount(request.Qty, request.Price, request.IsPrecentage, request.Discount),
		CreatedBy:    request.CreatedBy,
		IsPrecentage: request.IsPrecentage,
	}

	//saved data
	if err := ir.db.Create(&data).Error; err != nil {
		return data, err
	}

	//update the invoice header
	if err := ir.AfterSave(data); err != nil {
		return data, err
	}

	return data, nil
}

// IndexLine implements InvoiceRepositoryInterface.
func (ir *invoiceRepo) IndexLine(limit int, offset int, invoiceId int, q string) ([]model.InvoiceLineRespont, error) {
	data := []model.InvoiceLineRespont{}
	invoiceLine := []model.InvoiceLine{}

	//q param handler
	if q != "" {
		if err := ir.db.Joins("Product", ir.db.Where(model.GetSeatchParamPartnerV2(q))).Where(model.GetSeatchParamInvoiceLine(q, invoiceId)).Limit(limit).Offset(offset).Find(&data).Error; err != nil {
			return data, err
		}
	} else {
		if err := ir.db.Preload("Product").Where(model.GetSeatchParamInvoiceLine("", invoiceId)).Find(&invoiceLine).Error; err != nil {
			return data, err
		}
	}

	for _, line := range invoiceLine {
		indexResponse := model.InvoiceLineRespont{
			Invoice_id:      line.ID,
			Invoice_line_id: line.InvoiceID,
			Created_at:      line.CreatedAt,
			Product_name:    line.Product.Name,
			Product_id:      line.Product.ID,
			Qty:             line.Qty,
			Price:           line.Price,
			Amount:          line.Amount,
			Discount:        line.Discount,
			IsPrecentage:    line.IsPrecentage,
		}
		data = append(data, indexResponse)
	}

	return data, nil
}

// ShowLine implements InvoiceRepositoryInterface.
func (ir *invoiceRepo) ShowLine(id int) (model.InvoiceLine, error) {
	data := model.InvoiceLine{}

	if err := ir.db.Preload("User").Preload("Product").Preload("Invoice").First(&data, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return data, errors.New("data not found")
		}
		return data, err
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
	data.Amount = ir.handlingAmount(updatedInvoiceLine.Qty, updatedInvoiceLine.Price, updatedInvoiceLine.IsPrecentage, updatedInvoiceLine.Discount)

	if err := ir.db.Updates(&data).Error; err != nil {
		return data, err
	}

	//update the invoice header
	go ir.AfterSave(data)

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

	err = ir.AfterSave(data)
	if err != nil {
		return "ERORR", err
	}

	return "data deleted", nil
}

func (ir *invoiceRepo) handlingAmount(qty float64, price float64, isPrecentage bool, discount float64) float64 {
	var amount float64 = 0
	var discountValue float64 = 0

	if !isPrecentage {
		amount = (qty * price) - discount
	} else {
		discountValue = (discount / 100) * (qty * price)
		amount = (qty * price) - discountValue
	}

	return amount
}

func (ir *invoiceRepo) AfterSave(request model.InvoiceLine) error {

	//init the sql
	query := `
		UPDATE invoices as i 
		SET grand_total = (SELECT coalesce(SUM(amount),0) FROM invoice_lines WHERE invoice_id = i.id)
		WHERE i.id = ?;
	`
	err := ir.db.Exec(query, request.InvoiceID).Error
	if err != nil {
		return err
	}
	fmt.Println(err)
	return nil
}

func (ir *invoiceRepo) HandlingPaginationLine(q string, limit int, offset int, paymentID int) (int64, error) {
	var count int64 = 0
	data := model.InvoiceLine{}
	//q param handler
	if q != "" {
		if err := ir.db.Joins("Product", ir.db.Where(model.GetSeatchParamProductV2(q))).Where(model.GetSeatchParamInvoiceLine(q, paymentID)).Find(&data).Count(&count).Error; err != nil {
			return count, err
		}
	} else {
		if err := ir.db.Where(model.GetSeatchParamInvoiceLine("", paymentID)).Find(&data).Count(&count).Error; err != nil {
			return count, err
		}
	}
	return count, nil
}
