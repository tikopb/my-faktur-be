package invoice

import (
	"bemyfaktur/internal/model"
	"bemyfaktur/internal/model/constant"
	documentutil "bemyfaktur/internal/model/documentUtil"
	pgUtil "bemyfaktur/internal/model/paginationUtil"
	"errors"

	"gorm.io/gorm"
)

type invoiceRepo struct {
	db         *gorm.DB
	docUtil    documentutil.Repository
	pgUtilRepo pgUtil.Repository
}

func GetRepository(db *gorm.DB, docUtil documentutil.Repository, pgRepo pgUtil.Repository) InvoiceRepositoryInterface {
	return &invoiceRepo{
		db:         db,
		docUtil:    docUtil,
		pgUtilRepo: pgRepo,
	}
}

// Create implements InvoiceRepositoryInterface.
func (ir *invoiceRepo) Create(request model.InvoiceRequest, partner model.Partner) (model.InvoiceRespont, error) {
	data := model.InvoiceRespont{}

	//init for documentno
	documentno, err := ir.docUtil.GetDocumentNo(ir.getTableName())
	if err != nil {
		return data, err
	}

	invoiceData := model.Invoice{
		CreatedAt:         request.CreatedAt,
		CreatedBy:         request.CreatedBy,
		PartnerID:         request.PartnerID,
		BatchNo:           request.BatchNo,
		InvoiceLine:       []model.InvoiceLine{},
		Status:            constant.InvoiceStatusDraft, //every new document default as draft
		DocumentNo:        documentno,
		DocAction:         constant.InvoiceActionDraft,
		OustandingPayment: 0,
		Discount:          request.Discount,
		IsPrecentage:      request.IsPrecentage,
	}

	if err := ir.db.Create(&invoiceData).Error; err != nil {
		return data, err
	}

	//set return data
	//set data preload for return
	dataPreload, err := ir.ParsingInvoiceToInvoiceRequest(invoiceData)
	if err != nil {
		return dataPreload, err
	}

	return dataPreload, nil
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
func (ir *invoiceRepo) Index(limit int, offset int, q string) ([]model.InvoiceRespont, error) {
	data := []model.Invoice{}
	dataReturn := []model.InvoiceRespont{}

	//q param handler
	if q != "" {
		if err := ir.db.Joins("Partner", ir.db.Where(model.GetSearchParamPartnerV2(q))).Where(model.GetSeatchParamInvoice(q)).Limit(limit).Offset(offset).Find(&data).Error; err != nil {
			return dataReturn, err
		}
	} else {
		if err := ir.db.Order("created_at DESC").Limit(limit).Offset(offset).Find(&data).Error; err != nil {
			return dataReturn, err
		}
	}

	for _, invoice := range data {

		//get user return value
		dataPreload, err := ir.ParsingInvoiceToInvoiceRequest(invoice)
		if err != nil {
			return dataReturn, err
		}

		indexResponse := model.InvoiceRespont{
			ID:                invoice.ID,
			CreatedAt:         invoice.CreatedAt,
			DocumentNo:        dataPreload.DocumentNo,
			BatchNo:           invoice.BatchNo,
			Status:            invoice.Status,
			CreatedBy:         dataPreload.Partner.User,
			Discount:          invoice.Discount,
			IsPrecentage:      dataPreload.IsPrecentage,
			GrandTotal:        invoice.GrandTotal,
			OustandingPayment: invoice.OustandingPayment,
			Partner:           dataPreload.Partner,
		}
		dataReturn = append(dataReturn, indexResponse)
	}

	return dataReturn, nil
}

// Show implements InvoiceRepositoryInterface.
func (ir *invoiceRepo) Show(id int) (model.InvoiceRespont, error) {
	var data model.Invoice

	if err := ir.db.Preload("Partner").Preload("User").First(&data, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.InvoiceRespont{}, errors.New("data not found")
		}
	}

	return ir.ParsingInvoiceToInvoiceRequest(data)
}

// Show implements InvoiceRepositoryInterface.
func (ir *invoiceRepo) ShowInternal(id int) (model.Invoice, error) {
	var data model.Invoice

	if err := ir.db.Preload("Partner").Preload("User").First(&data, id).Error; err != nil {
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
	invoiceData, err := ir.ShowInternal(id) //get invoice Data

	if err != nil {
		return data, err
	}

	invoiceData.PartnerID = updatedInvoice.PartnerID
	invoiceData.Discount = updatedInvoice.Discount
	invoiceData.BatchNo = updatedInvoice.BatchNo

	//handling Grand Total
	invoiceData = ir.handlingGrandTotal(invoiceData)

	//validation docaction
	invoiceData, err = ir.DocProcess(invoiceData, string(updatedInvoice.DocAction))
	if err != nil {
		return data, err
	}

	//save the data
	if err := ir.db.Save(&invoiceData).Error; err != nil {
		return data, err
	}

	//set return data
	dataReturn, err := ir.ParsingInvoiceToInvoiceRequest(invoiceData)
	if err != nil {
		return data, err
	}

	return dataReturn, nil
}

func (ir *invoiceRepo) ParsingInvoiceToInvoiceRequest(invoice model.Invoice) (model.InvoiceRespont, error) {
	data := model.InvoiceRespont{}
	dataPreload, err := ir.ShowInternal(invoice.ID)
	if err != nil {
		return data, err
	}

	data = model.InvoiceRespont{
		ID:           dataPreload.ID,
		CreatedAt:    dataPreload.CreatedAt,
		GrandTotal:   dataPreload.GrandTotal,
		Discount:     dataPreload.Discount,
		BatchNo:      dataPreload.BatchNo,
		Status:       dataPreload.Status,
		CreatedBy:    dataPreload.User,
		Partner:      dataPreload.Partner,
		DocumentNo:   dataPreload.DocumentNo,
		DocAction:    dataPreload.DocAction,
		IsPrecentage: data.IsPrecentage,
	}

	return data, nil
}

func (pr *invoiceRepo) handlingGrandTotal(data model.Invoice) model.Invoice {
	if data.IsPrecentage {
		data.GrandTotal = data.GrandTotal - (data.GrandTotal * data.Discount / 100)
	} else {
		data.GrandTotal = data.GrandTotal - data.Discount
	}
	return data
}

func (ir *invoiceRepo) getTableName() string {
	return "invoices"
}

func (ir *invoiceRepo) HandlingPagination(q string, limit int, offset int) (int64, error) {
	var count int64 = 0
	data := model.Invoice{}
	//q param handler
	if q != "" {
		if err := ir.db.Joins("Partner", ir.db.Where(model.GetSearchParamPartnerV2(q))).Where(model.GetSeatchParamInvoice(q)).Find(&data).Count(&count).Error; err != nil {
			return count, err
		}
	} else {
		if err := ir.db.Order("created_at DESC").Limit(limit).Offset(offset).Find(&data).Count(&count).Error; err != nil {
			return count, err
		}
	}
	return count, nil
}
