package invoice

import (
	"bemyfaktur/internal/model"
	"bemyfaktur/internal/model/constant"
	"errors"

	"gorm.io/gorm"
)

type invoiceRepo struct {
	db *gorm.DB
}

func GetRepository(db *gorm.DB) InvoiceRepositoryInterface {
	return &invoiceRepo{
		db: db,
	}
}

// Create implements InvoiceRepositoryInterface.
func (ir *invoiceRepo) Create(invoice model.Invoice) (model.InvoiceCreateRespon, error) {
	data := model.InvoiceCreateRespon{}

	data = model.InvoiceCreateRespon{
		BatchNo:     invoice.BatchNo,
		CreatedAt:   invoice.CreatedAt,
		GrandTotal:  invoice.GrandTotal,
		Discount:    invoice.Discount,
		Status:      constant.InvoiceStatusComplete,
		PartnerName: invoice.Partner.Name,
		CreatedBy:   invoice.Partner.Name,
		User:        invoice.User,
		Partner:     invoice.Partner,
	}

	if err := ir.db.Create(&data).Error; err != nil {
		return data, err
	}

	return data, nil
}

// Delete implements InvoiceRepositoryInterface.
func (ir *invoiceRepo) Delete(id int) (string, error) {
	data, err := ir.Show(id)
	batchno := data.BatchNo
	if err != nil {
		return "", err
	}

	if err := ir.db.Delete(&data).Error; err != nil {
		return "", err
	}
	return batchno, nil
}

// Index implements InvoiceRepositoryInterface.
func (ir *invoiceRepo) Index(limit int, offset int) ([]model.InvoiceIndexRespont, error) {
	data := []model.Invoice{}
	dataReturn := []model.InvoiceIndexRespont{}
	if err := ir.db.Order("CreatedAt DESC").Limit(limit).Offset(offset).Find(&data).Error; err != nil {
		return dataReturn, err
	}

	for _, invoice := range data {
		indexResponse := model.InvoiceIndexRespont{
			ID:         invoice.ID,
			CreatedAt:  invoice.CreatedAt,
			CreatedBy:  invoice.CreatedBy,
			Partner:    invoice.Partner,
			GrandTotal: invoice.GrandTotal,
			Discount:   invoice.Discount,
			BatchNo:    invoice.BatchNo,
			Status:     invoice.Status,
		}
		dataReturn = append(dataReturn, indexResponse)
	}

	return dataReturn, nil
}

// Show implements InvoiceRepositoryInterface.
func (ir *invoiceRepo) Show(id int) (model.Invoice, error) {
	var data model.Invoice

	if err := ir.db.Where(model.Invoice{ID: id}).Preload("Invoice").First(&data).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return data, errors.New("data not found")
		}
	}
	return data, nil
}

// Update implements InvoiceRepositoryInterface.
func (ir *invoiceRepo) Update(id int, updatedInvoice model.Invoice) (model.Invoice, error) {
	data, err := ir.Show(id)

	if err != nil {
		return data, err
	}

	data.CreatedBy = updatedInvoice.CreatedBy
	data.PartnerID = updatedInvoice.PartnerID
	data.GrandTotal = updatedInvoice.GrandTotal
	data.Discount = updatedInvoice.Discount
	data.BatchNo = updatedInvoice.BatchNo
	data.Status = updatedInvoice.Status

	//save the data
	if err := ir.db.Save(&data).Error; err != nil {
		return data, nil
	}

	panic("unimplemented")
}
