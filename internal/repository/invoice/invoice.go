package invoice

import (
	"bemyfaktur/internal/model"
	"bemyfaktur/internal/model/constant"
	documentutil "bemyfaktur/internal/model/documentUtil"
	pgUtil "bemyfaktur/internal/model/paginationUtil"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
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

	//init for documentno
	documentno, err := ir.docUtil.GetDocumentNo(ir.getTableName())
	if err != nil {
		return model.InvoiceRespont{}, err
	}

	invoiceData := model.Invoice{
		CreatedBy:         request.CreatedById,
		UpdatedBy:         request.UpdatedById,
		PartnerID:         request.PartnerId,
		GrandTotal:        0, // all new invoice data is 0
		Discount:          request.Discount,
		BatchNo:           request.BatchNo,
		Status:            constant.InvoiceStatusDraft, // all new data set to draft
		DocAction:         constant.InvoiceActionDraft, // all new data set to draft
		OustandingPayment: 0,
		DocumentNo:        documentno,
		IsPrecentage:      request.IsPrecentage,
		PayDate:           request.PayDate,
	}

	if err := ir.db.Create(&invoiceData).Error; err != nil {
		return model.InvoiceRespont{}, err
	}

	//set return data
	//parsing the data return
	dataPreload, err := ir.ParsingInvoiceToInvoiceRequest(invoiceData)
	if err != nil {
		return dataPreload, err
	}

	return dataPreload, nil
}

// create data with header and line simultaneously
func (ir *invoiceRepo) CreateInvoiceV2(request model.InvoiceRequest, requestLines []model.InvoiceLineRequest, partner model.Partner) (model.InvoiceRespont, []model.InvoiceLineRespont, error) {
	tx := ir.db.Begin()

	//create header
	header, err := ir.Create(request, partner)
	if err != nil {
		return model.InvoiceRespont{}, []model.InvoiceLineRespont{}, err
	}
	tx.Commit()

	returnLine := []model.InvoiceLineRespont{}
	invoice, err := ir.ShowInternal(header.ID)
	if err != nil {
		return model.InvoiceRespont{}, []model.InvoiceLineRespont{}, err

	}
	for _, line := range requestLines {
		line.InvoiceId = invoice.ID
		dataLine, err := ir.CreateLine(line)
		if err != nil {
			ir.Delete(header.ID)
			for _, lineGenerate := range returnLine {
				ir.Delete(lineGenerate.ID)
			}

			return model.InvoiceRespont{}, []model.InvoiceLineRespont{}, err
		}

		returnLine = append(returnLine, dataLine)
	}

	headerReturn, err := ir.Show(invoice.UUID)
	if err != nil {
		return model.InvoiceRespont{}, []model.InvoiceLineRespont{}, err

	}

	return headerReturn, returnLine, nil
}

// Delete implements InvoiceRepositoryInterface.
func (ir *invoiceRepo) Delete(id uuid.UUID) (string, error) {
	data, err := ir.ShowInternal(id)
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
func (ir *invoiceRepo) Index(limit int, offset int, q string, order []string, dateFrom string, dateTo string) ([]model.InvoiceRespont, error) {
	data := []model.Invoice{}
	dataReturn := []model.InvoiceRespont{}

	//order by handling
	orderParam := ""
	if len(order) != 0 {
		orderParam = fmt.Sprintf(" %s", string(order[0]))
	}

	if orderParam != "" {
		if err := ir.db.Preload("Partner").Preload("User").Preload("UserUpdated").Joins("Partner", ir.db.Where(model.GetSearchParamPartnerV2(q))).Where(model.GetSeatchParamInvoiceV2(dateFrom, dateTo, q)).Limit(limit).Offset(offset).Order(orderParam).Find(&data).Error; err != nil {
			return dataReturn, err
		}
	} else {
		if err := ir.db.Preload("Partner").Preload("User").Preload("UserUpdated").Joins("Partner", ir.db.Where(model.GetSearchParamPartnerV2(q))).Where(model.GetSeatchParamInvoiceV2(dateFrom, dateTo, q)).Limit(limit).Offset(offset).Find(&data).Error; err != nil {
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
			ID:                dataPreload.ID,
			CreatedAt:         dataPreload.CreatedAt,
			GrandTotal:        dataPreload.GrandTotal,
			Discount:          dataPreload.Discount,
			BatchNo:           dataPreload.BatchNo,
			Status:            dataPreload.Status,
			DocAction:         dataPreload.DocAction,
			OustandingPayment: dataPreload.OustandingPayment,
			DocumentNo:        dataPreload.DocumentNo,
			IsPrecentage:      dataPreload.IsPrecentage,
			PayDate:           dataPreload.PayDate,
			CreatedBy:         dataPreload.CreatedBy,
			UpdatedBy:         dataPreload.UpdatedBy,
			Partner:           dataPreload.Partner,
		}
		dataReturn = append(dataReturn, indexResponse)
	}

	return dataReturn, nil
}

// Show implements InvoiceRepositoryInterface.
func (ir *invoiceRepo) Show(id uuid.UUID) (model.InvoiceRespont, error) {
	var data model.Invoice

	if err := ir.db.Preload("Partner").Preload("User").Preload("UserUpdated").Where(model.Invoice{UUID: id}).First(&data).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.InvoiceRespont{}, errors.New("data not found")
		}
	}

	return ir.ParsingInvoiceToInvoiceRequest(data)
}

// Show implements InvoiceRepositoryInterface.
func (ir *invoiceRepo) ShowInternal(id uuid.UUID) (model.Invoice, error) {
	var data model.Invoice

	if err := ir.db.Preload("Partner").Preload("User").Preload("UserUpdated").Where(model.Invoice{UUID: id}).First(&data).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return data, errors.New("data not found")
		}
	}

	return data, nil
}

// Update implements InvoiceRepositoryInterface.
func (ir *invoiceRepo) Update(id uuid.UUID, request model.InvoiceRequest) (model.InvoiceRespont, error) {
	//set var
	data := model.InvoiceRespont{}
	invoiceData, err := ir.ShowInternal(id) //get invoice Data

	if err != nil {
		return data, err
	}

	invoiceData.UpdateAt = time.Now()
	invoiceData.UpdatedBy = request.UpdatedById
	invoiceData.PartnerID = request.PartnerId
	invoiceData.Discount = request.Discount
	invoiceData.BatchNo = request.BatchNo

	//handling Grand Total
	invoiceData = ir.handlingGrandTotal(invoiceData)

	//validation docaction
	invoiceData, err = ir.DocProcess(invoiceData, string(request.DocAction))
	if err != nil {
		return data, err
	}

	//save the data
	if err := ir.db.Updates(&invoiceData).Error; err != nil {
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
	//parsing to patial verstion first!
	createdBy := model.UserPartial{
		UserId:   invoice.User.ID,
		Username: invoice.User.Username,
	}
	updateBy := model.UserPartial{
		UserId:   invoice.UserUpdated.ID,
		Username: invoice.UserUpdated.Username,
	}
	partner := model.PartnerPartialRespon{
		UUID: invoice.Partner.UUID,
		Name: invoice.Partner.Name,
	}

	data := model.InvoiceRespont{
		ID:                invoice.UUID,
		CreatedAt:         invoice.CreatedAt,
		UpdatedAt:         invoice.UpdateAt,
		GrandTotal:        invoice.GrandTotal,
		Discount:          invoice.Discount,
		BatchNo:           invoice.BatchNo,
		Status:            invoice.Status,
		DocAction:         invoice.DocAction,
		OustandingPayment: invoice.OustandingPayment,
		DocumentNo:        invoice.DocumentNo,
		IsPrecentage:      invoice.IsPrecentage,
		PayDate:           invoice.PayDate,
		CreatedBy:         createdBy,
		UpdatedBy:         updateBy,
		Partner:           partner,
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

func (ir *invoiceRepo) HandlingPagination(q string, limit int, offset int, dateFrom string, dateTo string) (int64, error) {
	var count int64 = 0
	data := model.Invoice{}

	if err := ir.db.Joins("Partner", ir.db.Where(model.GetSearchParamPartnerV2(q))).Where(model.GetSeatchParamInvoiceV2(dateFrom, dateTo, q)).Find(&data).Count(&count).Error; err != nil {
		return count, err
	}

	return count, nil
}
