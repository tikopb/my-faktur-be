package invoice

import (
	"bemyfaktur/internal/model"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// CreateLine implements InvoiceRepositoryInterface.
func (ir *invoiceRepo) CreateLine(request model.InvoiceLineRequest) (model.InvoiceLineRespont, error) {

	tx := ir.db.Begin()

	data := model.InvoiceLine{
		InvoiceID:      request.InvoiceId,
		ProductID:      request.ProductID,
		Price:          request.Price,
		Discount:       request.Discount,
		Qty:            request.Qty,
		IsPrecentage:   request.IsPrecentage,
		Amount:         ir.handlingAmount(request.Qty, request.Price, request.IsPrecentage, request.Discount),
		CreatedBy:      request.CreatedById,
		UpdatedBy:      request.CreatedById,
		OrganizationId: request.OrganizationId,
	}

	//saved data
	if err := tx.Create(&data).Error; err != nil {
		return model.InvoiceLineRespont{}, err
	}

	tx.Commit()

	//update the invoice header
	if err := ir.AfterSave(data); err != nil {
		return model.InvoiceLineRespont{}, err
	}

	//return value set
	return ir.ParsingInvoiceLineToInvoiceRequest(data)
}

// IndexLine implements InvoiceRepositoryInterface.
func (ir *invoiceRepo) IndexLine(limit int, offset int, invoiceId int, q string, order []string) ([]model.InvoiceLineRespont, error) {
	data := []model.InvoiceLineRespont{}
	invoiceLine := []model.InvoiceLine{}

	//q param handler
	if q != "" {
		if len(order) > 0 {
			if err := ir.db.Preload("Invoice").Preload("Product").Preload("User").Preload("UserUpdated").Joins("Product", ir.db.Where(model.GetSeatchParamProductV2(q))).Where(model.GetSeatchParamInvoiceLine(q, invoiceId)).Limit(limit).Offset(offset).Order(order[0]).Find(&invoiceLine).Error; err != nil {
				return []model.InvoiceLineRespont{}, err
			}
		} else {
			if err := ir.db.Preload("Invoice").Preload("Product").Preload("User").Preload("UserUpdated").Joins("INNER JOIN products Product ON invoice_lines.product_id = Product.id AND ( " + model.GetSeatchParamProductV2(q) + " )").Where(model.GetSeatchParamInvoiceLine(q, invoiceId)).Limit(limit).Offset(offset).Find(&invoiceLine).Error; err != nil {
				return []model.InvoiceLineRespont{}, err
			}
		}
	} else {
		if len(order) > 0 {
			if err := ir.db.Preload("Invoice").Preload("Product").Preload("User").Preload("UserUpdated").Where(model.GetSeatchParamInvoiceLine("", invoiceId)).Limit(limit).Offset(offset).Order(order[0]).Find(&invoiceLine).Error; err != nil {
				return []model.InvoiceLineRespont{}, err
			}
		} else {
			if err := ir.db.Preload("Invoice").Preload("Product").Preload("User").Preload("UserUpdated").Where(model.GetSeatchParamInvoiceLine("", invoiceId)).Limit(limit).Offset(offset).Find(&invoiceLine).Error; err != nil {
				return []model.InvoiceLineRespont{}, err
			}
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

	if err := ir.db.Preload("User").Preload("UserUpdated").Preload("Product").Preload("Invoice").Where(model.InvoiceLine{UUID: id}).First(&data).Error; err != nil {
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

	if err := ir.db.Preload("User").Preload("UserUpdated").Preload("Product").Preload("Invoice").Where(model.InvoiceLine{UUID: id}).First(&data).Error; err != nil {
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

	//updated updatedby field becaue gorm overwrite the updatedby colum ## not fixet yed!
	err = ir.db.Exec("UPDATE invoice_lines SET updated_by = ? WHERE uuid = ?", request.UpdatedById, id).Error
	if err != nil {
		return model.InvoiceLineRespont{}, err
	}

	data.UpdateAt = time.Now()
	data.ProductID = request.ProductID
	data.Price = request.Price
	data.Discount = request.Discount
	data.Qty = request.Qty
	data.IsPrecentage = request.IsPrecentage
	data.Amount = ir.handlingAmount(request.Qty, request.Price, request.IsPrecentage, request.Discount)

	if err := ir.db.Omit("UpdatedBy").Updates(&data).Error; err != nil {
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

// calculated the amount base on qty, price and discount
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

// bring updated after saved is success
func (ir *invoiceRepo) AfterSave(request model.InvoiceLine) error {

	//set the total_line
	query := `
		UPDATE invoices as i 
		SET 
		total_line = (
			SELECT coalesce(SUM(amount), 0) 
			FROM invoice_lines 
			WHERE invoice_id = i.id
		)
		WHERE i.id = ?;
	`
	err := ir.db.Exec(query, request.InvoiceID).Error
	if err != nil {
		return err
	}
	//set the grand total
	query = `
		UPDATE invoices as i 
		SET 
		grand_total = CASE
			WHEN i.isprecentage = true THEN i.total_line - (i.total_line * i.discount / 100)
			ELSE i.total_line - i.discount
		END
		WHERE i.id = ?;
	`
	err = ir.db.Exec(query, request.InvoiceID).Error
	if err != nil {
		return err
	}

	return nil
}

func (ir *invoiceRepo) HandlingPaginationLine(q string, limit int, offset int, invoiceId int) (int64, error) {
	var count int64 = 0
	data := model.InvoiceLine{}
	//q param handler
	if q != "" {
		if err := ir.db.Joins("Product", ir.db.Where(model.GetSeatchParamProductV2(q))).Where(model.GetSeatchParamInvoiceLine(q, invoiceId)).Find(&data).Count(&count).Error; err != nil {
			return count, err
		}
	} else {
		if err := ir.db.Where(model.GetSeatchParamInvoiceLine("", invoiceId)).Find(&data).Count(&count).Error; err != nil {
			return count, err
		}
	}
	return count, nil
}

func (ir *invoiceRepo) ParsingInvoiceLineToInvoiceRequest(line model.InvoiceLine) (model.InvoiceLineRespont, error) {
	dataRe, err := ir.ShowLineLinternal(line.UUID)
	if err != nil {
		panic("erorr parsing")
	}
	createdBy := model.UserPartial{
		UserId:   dataRe.User.ID,
		Username: dataRe.User.Username,
	}
	updateBy := model.UserPartial{
		UserId:   dataRe.UserUpdated.ID,
		Username: dataRe.UserUpdated.Username,
	}

	invoice := model.InvoicePartialRespont{
		Id:         dataRe.Invoice.ID,
		UUID:       dataRe.Invoice.UUID,
		BatchNo:    dataRe.Invoice.BatchNo,
		Documentno: dataRe.Invoice.DocumentNo,
	}

	product := model.ProductPartialRespon{
		UUID: dataRe.Product.UUID,
		Name: dataRe.Product.Name,
	}

	data := model.InvoiceLineRespont{
		ID:           dataRe.UUID,
		CreatedAt:    dataRe.CreatedAt,
		UpdatedAt:    dataRe.UpdateAt,
		Qty:          dataRe.Qty,
		Price:        dataRe.Price,
		Amount:       dataRe.Amount,
		Discount:     dataRe.Discount,
		IsPrecentage: dataRe.IsPrecentage,
		CreatedBy:    createdBy,
		UpdatedBy:    updateBy,
		Invoice:      invoice,
		Product:      product,
	}

	return data, nil
}
