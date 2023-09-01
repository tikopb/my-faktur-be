package invoice

import (
	"bemyfaktur/internal/model"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

// CreateLine implements InvoiceRepositoryInterface.
func (ir *invoiceRepo) CreateLine(invoice model.InvoiceLine) (model.InvoiceLine, error) {
	data := model.InvoiceLine{
		InvoiceID:    invoice.InvoiceID,
		ProductID:    invoice.ProductID,
		Price:        invoice.Price,
		Discount:     invoice.Discount,
		Qty:          invoice.Qty,
		Amount:       ir.handlingAmount(invoice.Qty, invoice.Price, invoice.IsPrecentage, invoice.Discount),
		CreatedBy:    "1", //##@ until security module fixed
		IsPrecentage: invoice.IsPrecentage,
	}

	//saved data
	if err := ir.db.Create(&data).Error; err != nil {
		return data, err
	}

	//update the invoice header
	go ir.updateInvoiceHeader(data.InvoiceID)

	return data, nil
}

// IndexLine implements InvoiceRepositoryInterface.
func (ir *invoiceRepo) IndexLine(limit int, offset int, invoiceId int, q string) ([]model.InvoiceLineRespont, error) {
	data := []model.InvoiceLineRespont{}

	query := `
        SELECT 
            i.id AS invoice, 
            il.id AS invoice_line_id, 
            il.created_at AS invoice_line_created_at, 
            p.name AS product_name, 
			p.id as product_id,
			il.qty,
            il.price, 
            il.amount, 
            il.discount 
        FROM invoice_lines il
        JOIN invoices i ON il.invoice_id = i.id
        JOIN products p ON il.product_id = p.id
        WHERE il.invoice_id = ? 
    `
	if q != "" {
		query += `
			and p.name = ?
		`
	}

	query += ` limit ? offset ?`

	fmt.Println(query)
	if err := ir.db.Raw(query, invoiceId, q, limit, offset).Scan(&data).Error; err != nil {
		return data, err
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

	if err := ir.db.Save(&data).Error; err != nil {
		return data, err
	}

	//update the invoice header
	go ir.updateInvoiceHeader(data.InvoiceID)

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

func (ir *invoiceRepo) updateInvoiceHeader(invoiceId int) {
	//init the sql
	sql := `
		UPDATE invoices as i
		SET grand_total = (SELECT SUM(amount) FROM invoice_lines WHERE invoice_id = i.id)
		WHERE i.id = ?;
	`
	ir.db.Raw(sql, invoiceId)
}
