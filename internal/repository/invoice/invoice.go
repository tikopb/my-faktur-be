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
func (ir *invoiceRepo) Create(request model.InvoiceRequest, partner model.Partner, user model.User) (model.InvoiceRespont, error) {
	data := model.InvoiceRespont{}

	invoiceData := model.Invoice{
		BatchNo:    request.BatchNo,
		CreatedAt:  request.CreatedAt,
		GrandTotal: request.GrandTotal,
		Discount:   request.Discount,
		Status:     constant.InvoiceStatusComplete,
		CreatedBy:  "1", //##@ UNTIL SECURITY MODULE DONE
		PartnerID:  partner.ID,
	}

	if err := ir.db.Create(&invoiceData).Error; err != nil {
		return data, err
	}

	//set return data
	data = model.InvoiceRespont{
		ID:         invoiceData.ID,
		CreatedAt:  invoiceData.CreatedAt,
		GrandTotal: invoiceData.GrandTotal,
		Discount:   invoiceData.Discount,
		BatchNo:    invoiceData.BatchNo,
		Status:     invoiceData.Status,
		CreatedBy:  invoiceData.User,
		Partner:    invoiceData.Partner,
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
func (ir *invoiceRepo) Index(limit int, offset int) ([]model.InvoiceRespont, error) {
	data := []model.Invoice{}
	dataReturn := []model.InvoiceRespont{}
	if err := ir.db.Order("CreatedAt DESC").Limit(limit).Offset(offset).Find(&data).Error; err != nil {
		return dataReturn, err
	}

	for _, invoice := range data {
		indexResponse := model.InvoiceRespont{
			ID:         invoice.ID,
			CreatedAt:  invoice.CreatedAt,
			GrandTotal: invoice.GrandTotal,
			Discount:   invoice.Discount,
			BatchNo:    invoice.BatchNo,
			Status:     invoice.Status,
			CreatedBy:  invoice.User,
			Partner:    invoice.Partner,
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
func (ir *invoiceRepo) Update(id int, updatedInvoice model.Invoice) (model.InvoiceRespont, error) {
	//set var
	data := model.InvoiceRespont{}
	invoiceData, err := ir.Show(id) //get invoice Data

	if err != nil {
		return data, err
	}

	invoiceData.CreatedBy = updatedInvoice.CreatedBy
	invoiceData.PartnerID = updatedInvoice.PartnerID
	invoiceData.GrandTotal = updatedInvoice.GrandTotal
	invoiceData.Discount = updatedInvoice.Discount
	invoiceData.BatchNo = updatedInvoice.BatchNo
	invoiceData.Status = updatedInvoice.Status

	//save the data
	if err := ir.db.Save(&invoiceData).Error; err != nil {
		return data, err
	}

	//set return data
	data = model.InvoiceRespont{
		ID:         invoiceData.ID,
		CreatedAt:  invoiceData.CreatedAt,
		GrandTotal: invoiceData.GrandTotal,
		Discount:   invoiceData.Discount,
		BatchNo:    invoiceData.BatchNo,
		Status:     invoiceData.Status,
		CreatedBy:  invoiceData.User,
		Partner:    invoiceData.Partner,
	}

	return data, nil
}
