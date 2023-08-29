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
func (ir *invoiceRepo) Create(request model.InvoiceRequest, partner model.Partner) (model.InvoiceRespont, error) {
	data := model.InvoiceRespont{}

	invoiceData := model.Invoice{
		ID:          0,
		CreatedAt:   request.CreatedAt,
		CreatedBy:   "1",
		PartnerID:   request.PartnerID,
		GrandTotal:  request.GrandTotal,
		Discount:    request.Discount,
		BatchNo:     request.BatchNo,
		InvoiceLine: []model.InvoiceLine{},
		Status:      constant.InvoiceStatusDraft,
	}

	if err := ir.db.Create(&invoiceData).Error; err != nil {
		return data, err
	}

	//set return data
	//get user return value
	var userData model.User
	if err := ir.db.First(&userData, invoiceData.CreatedBy).Error; err != nil {
		return data, err
	}
	//get partner return value
	var partnerData model.Partner
	if err := ir.db.First(&partnerData, invoiceData.PartnerID).Preload("User").Error; err != nil {
		return data, err
	}

	//get partner user value

	data = model.InvoiceRespont{
		ID:         invoiceData.ID,
		CreatedAt:  invoiceData.CreatedAt,
		GrandTotal: invoiceData.GrandTotal,
		Discount:   invoiceData.Discount,
		BatchNo:    invoiceData.BatchNo,
		Status:     invoiceData.Status,
		CreatedBy:  userData,
		Partner:    partnerData,
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
	if err := ir.db.Order("created_at DESC").Limit(limit).Offset(offset).Find(&data).Error; err != nil {
		return dataReturn, err
	}

	for _, invoice := range data {

		//get user return value
		//get user return value
		var userData model.User
		if err := ir.db.First(&userData, invoice.CreatedBy).Error; err != nil {
			return dataReturn, err
		}
		//get partner return value
		var partnerData model.Partner
		if err := ir.db.First(&partnerData, invoice.PartnerID).Preload("User").Error; err != nil {
			return dataReturn, err
		}

		indexResponse := model.InvoiceRespont{
			ID:         invoice.ID,
			CreatedAt:  invoice.CreatedAt,
			GrandTotal: invoice.GrandTotal,
			Discount:   invoice.Discount,
			BatchNo:    invoice.BatchNo,
			Status:     invoice.Status,
			CreatedBy:  userData,
			Partner:    partnerData,
		}
		dataReturn = append(dataReturn, indexResponse)
	}

	return dataReturn, nil
}

// Show implements InvoiceRepositoryInterface.
func (ir *invoiceRepo) Show(id int) (model.Invoice, error) {
	var data model.Invoice

	if err := ir.db.Preload("Invoice").Preload("Partner").Preload("User").First(&data, id).Error; err != nil {
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
	//get user return value
	var userData model.User
	if err := ir.db.First(&userData, invoiceData.CreatedBy).Error; err != nil {
		return data, err
	}
	//get partner return value
	var partnerData model.Partner
	if err := ir.db.First(&partnerData, invoiceData.PartnerID).Preload("User").Error; err != nil {
		return data, err
	}

	data = model.InvoiceRespont{
		ID:         invoiceData.ID,
		CreatedAt:  invoiceData.CreatedAt,
		GrandTotal: invoiceData.GrandTotal,
		Discount:   invoiceData.Discount,
		BatchNo:    invoiceData.BatchNo,
		Status:     invoiceData.Status,
		CreatedBy:  userData,
		Partner:    partnerData,
	}

	return data, nil
}
