package invoice

import (
	"bemyfaktur/internal/model"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// CreateLine implements InvoiceRepositoryInterface.
func (ir *invoiceRepo) CreateLine(request model.InvoiceLineRequest) (model.InvoiceLineRespont, error) {
	data := model.InvoiceLine{
		InvoiceID:    request.InvoiceId,
		ProductID:    request.ProductID,
		Price:        request.Price,
		Discount:     request.Discount,
		Qty:          request.Qty,
		Amount:       ir.handlingAmount(request.Qty, request.Price, request.IsPrecentage, request.Discount),
		CreatedBy:    request.CreatedById,
		UpdatedBy:    request.UpdatedById,
		IsPrecentage: request.IsPrecentage,
	}

	//saved data
	if err := ir.db.Create(&data).Error; err != nil {
		return model.InvoiceLineRespont{}, err
	}

	//update the invoice header
	if err := ir.AfterSave(data); err != nil {
		return model.InvoiceLineRespont{}, err
	}

	//return value set
	return ir.ParsingInvoiceLineToInvoiceRequest(data)
}

// IndexLine implements InvoiceRepositoryInterface.
func (ir *invoiceRepo) IndexLine(limit int, offset int, invoiceId int, q string) ([]model.InvoiceLineRespont, error) {
	data := []model.InvoiceLineRespont{}
	invoiceLine := []model.InvoiceLine{}

	//q param handler
	if q != "" {
		if err := ir.db.Joins("Product", ir.db.Where(model.GetSearchParamPartnerV2(q))).Where(model.GetSeatchParamInvoiceLine(q, invoiceId)).Limit(limit).Offset(offset).Find(&data).Error; err != nil {
			return data, err
		}
	} else {
		if err := ir.db.Preload("Product").Where(model.GetSeatchParamInvoiceLine("", invoiceId)).Find(&invoiceLine).Error; err != nil {
			return data, err
		}
	}

	for _, line := range invoiceLine {
		parsing, err := ir.ParsingInvoiceLineToInvoiceRequest(line)
		if err != nil {
			return []model.InvoiceLineRespont{}, nil
		}
		data = append(data, parsing)
	}

	return data, nil
}

// ShowLine implements InvoiceRepositoryInterface.
func (ir *invoiceRepo) ShowLine(id uuid.UUID) (model.InvoiceLineRespont, error) {
	data := model.InvoiceLine{}

	if err := ir.db.Preload("User").Preload("Product").Preload("Invoice").First(&data, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.InvoiceLineRespont{}, errors.New("data not found")
		}
		return model.InvoiceLineRespont{}, err
	}
	return ir.ParsingInvoiceLineToInvoiceRequest(data)
}

// ShowLine implements InvoiceRepositoryInterface.
func (ir *invoiceRepo) ShowLineLinternal(id uuid.UUID) (model.InvoiceLine, error) {
	data := model.InvoiceLine{}

	if err := ir.db.Preload("User").Preload("UserUpdated").Preload("Product").Preload("Invoice").First(&data, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return data, errors.New("data not found")
		}
		return data, err
	}
	return data, nil
}

// UpdateLine implements InvoiceRepositoryInterface.
func (ir *invoiceRepo) UpdateLine(id uuid.UUID, request model.InvoiceLineRequest) (model.InvoiceLineRespont, error) {
	data, err := ir.ShowLineLinternal(id)

	if err != nil {
		return model.InvoiceLineRespont{}, nil
	}

	data.UpdateAt = time.Now()
	data.ProductID = request.ProductID
	data.Price = request.Price
	data.Discount = request.Discount
	data.Qty = request.Qty
	data.Amount = ir.handlingAmount(request.Qty, request.Price, request.IsPrecentage, request.Discount)

	if err := ir.db.Save(&data).Error; err != nil {
		return model.InvoiceLineRespont{}, err
	}

	//update the invoice header
	go ir.AfterSave(data)

	return ir.ParsingInvoiceLineToInvoiceRequest(data)
}

// DeleteLine implements InvoiceRepositoryInterface.
func (ir *invoiceRepo) DeleteLine(id uuid.UUID) (string, error) {
	data, err := ir.ShowLineLinternal(id)
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

func (ir *invoiceRepo) ParsingInvoiceLineToInvoiceRequest(line model.InvoiceLine) (model.InvoiceLineRespont, error) {

	createdBy := model.UserPartial{
		UserId:   line.User.ID,
		Username: line.User.Username,
	}
	updateBy := model.UserPartial{
		UserId:   line.UserUpdated.ID,
		Username: line.UserUpdated.Username,
	}

	invoice := model.InvoicePartialRespont{
		ID:         line.Invoice.UUID,
		BatchNo:    line.Invoice.BatchNo,
		DocumentNo: line.Invoice.DocumentNo,
	}

	product := model.ProductPartialRespon{
		UUID: line.Product.UUID,
		Name: line.Product.Name,
	}

	data := model.InvoiceLineRespont{
		ID:           line.UUID,
		CreatedAt:    line.CreatedAt,
		UpdatedAt:    line.UpdateAt,
		Qty:          line.Qty,
		Price:        line.Price,
		Amount:       line.Amount,
		Discount:     line.Discount,
		IsPrecentage: line.IsPrecentage,
		CreatedBy:    createdBy,
		UpdatedBy:    updateBy,
		Invoice:      invoice,
		Product:      product,
	}

	return data, nil
}
