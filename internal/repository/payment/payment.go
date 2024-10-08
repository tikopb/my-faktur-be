package payment

import (
	"bemyfaktur/internal/model"
	documentutil "bemyfaktur/internal/model/documentUtil"
	pgUtil "bemyfaktur/internal/model/paginationUtil"
	"encoding/json"
	"io"
	"os"
	"time"

	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type paymentRepo struct {
	db         *gorm.DB
	docUtil    documentutil.Repository
	pgUtilRepo pgUtil.Repository
}

func GetRepository(db *gorm.DB, docUtil documentutil.Repository, pgRepo pgUtil.Repository) PaymentRepositoryinterface {
	return &paymentRepo{
		db:         db,
		docUtil:    docUtil,
		pgUtilRepo: pgRepo,
	}
}

type mapperDoc struct {
	Status   string `json:"status"`
	Doaction string `json:"doaction"`
}

// Index implements PaymentRepositoryinterface.
func (pr *paymentRepo) Index(limit int, offset int, q string, order []string, dateFrom string, dateTo string) ([]model.PaymentRespont, error) {
	data := []model.Payment{}
	dataReturn := []model.PaymentRespont{}

	//order by handling
	orderParam := ""
	if len(order) != 0 {
		orderParam = fmt.Sprintf(" %s", string(order[0]))
	}

	//q param handler
	if orderParam != "" {
		if err := pr.db.Preload("Partner").Preload("User").Preload("UserUpdated").Joins("Partner", pr.db.Where(model.GetSearchParamPartnerV2(q))).Where(model.GetSeatchParamPayment(dateFrom, dateTo, q)).Limit(limit).Offset(offset).Order(orderParam).Find(&data).Error; err != nil {
			return dataReturn, err
		}
	} else {
		if err := pr.db.Preload("Partner").Preload("User").Preload("UserUpdated").Order("created_at DESC").Where(model.GetSeatchParamPayment(dateFrom, dateTo, q)).Limit(limit).Offset(offset).Find(&data).Error; err != nil {
			return dataReturn, err
		}
	}

	for _, payment := range data {
		dataPreload, err := pr.parsingPaymentToPaymentRespont(payment)
		if err != nil {
			return dataReturn, err
		}

		indexResponse := model.PaymentRespont{
			ID:           dataPreload.ID,
			DocumentNo:   dataPreload.DocumentNo,
			BatchNo:      dataPreload.BatchNo,
			IsPrecentage: dataPreload.IsPrecentage,
			GrandTotal:   dataPreload.GrandTotal,
			Discount:     dataPreload.Discount,
			Status:       dataPreload.Status,
			DoAction:     dataPreload.DoAction,
			CreatedBy:    dataPreload.CreatedBy,
			UpdatedBy:    dataPreload.UpdatedBy,
			Partner:      dataPreload.Partner,
		}
		dataReturn = append(dataReturn, indexResponse)
	}

	return dataReturn, nil
}

// Create implements PaymentRepositoryinterface.
func (pr *paymentRepo) Create(payment model.PaymentRequest) (model.PaymentRespont, error) {
	data := model.PaymentRespont{}

	//init for documentno
	documentno, err := pr.docUtil.GetDocumentNo(pr.getTableName())
	if err != nil {
		return data, err
	}

	//change the value from string to timestamp format
	dateStr := payment.PayDateString
	date, err := time.Parse("02-01-2006", dateStr)
	if err != nil {
		return model.PaymentRespont{}, err
	}
	payment.PayDate = date

	paymentData := model.Payment{
		CreatedBy:      payment.CreatedBy,
		PartnerID:      payment.PartnerID,
		UpdatedBy:      payment.UpdatedBy,
		GrandTotal:     0,
		Discount:       0,
		BatchNo:        payment.BatchNo,
		DocumentNo:     documentno,
		OrganizationId: payment.OrganizationId,
	}

	if err := pr.db.Create(&paymentData).Error; err != nil {
		return data, err
	}

	//set return data value
	dataPreload, err := pr.parsingPaymentToPaymentRespont(paymentData)
	if err != nil {
		return dataPreload, err
	}

	return dataPreload, nil
}

func (pr *paymentRepo) CreateV2(request model.PaymentRequestV2) (model.PaymentRespont, error) {
	//create header
	header, err := pr.Create(request.Header)
	if err != nil {
		return model.PaymentRespont{}, err
	}

	returnLine := []model.PaymentLineRespont{}

	payment, err := pr.ShowInternal(header.UUID)
	if err != nil {
		return model.PaymentRespont{}, err
	}
	for _, line := range request.Line {
		line.PaymentID = payment.ID
		lineCreated, err := pr.CreateLine(line)
		if err != nil {
			pr.Delete(payment.UUID)
			for _, lineGenerate := range returnLine {
				pr.Delete(lineGenerate.ID)
			}
			return model.PaymentRespont{}, err
		}
		returnLine = append(returnLine, lineCreated)
	}

	headerReturn, err := pr.Show(header.UUID)
	if err != nil {
		return model.PaymentRespont{}, err
	}

	return headerReturn, nil
}

// Show implements PaymentRepositoryinterface.
func (pr *paymentRepo) Show(id uuid.UUID) (model.PaymentRespont, error) {
	data := model.Payment{}

	if err := pr.db.Preload("Partner").Preload("User").Preload("UserUpdated").Where(&model.Payment{UUID: id}).First(&data).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.PaymentRespont{}, errors.New("data not found")
		}
	}

	return pr.parsingPaymentToPaymentRespont(data)
}

func (pr *paymentRepo) ShowInternal(id uuid.UUID) (model.Payment, error) {
	data := model.Payment{}

	if err := pr.db.Preload("Partner").Preload("User").Preload("UserUpdated").Where(&model.Payment{UUID: id}).First(&data).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.Payment{}, errors.New("data not found")
		}
	}
	return data, nil
}

// Update implements PaymentRepositoryinterface.
func (pr *paymentRepo) Update(id uuid.UUID, updatedPayment model.PaymentRequest) (model.PaymentRespont, error) {
	//set var
	data := model.PaymentRespont{}
	paymentData, err := pr.ShowInternal(id) //get payment data
	if err != nil {
		return data, err
	}

	paymentData, err = pr.BeforeUpdate(paymentData, string(updatedPayment.DocAction))
	if err != nil {
		return data, err
	}

	//updated updatedby field becaue gorm overwrite the updatedby colum ## not fixet yed!
	err = pr.db.Exec("UPDATE payments SET updated_by = ? WHERE uuid = ?", updatedPayment.UpdatedBy, id).Error
	if err != nil {
		return model.PaymentRespont{}, err
	}

	//change the value from string to timestamp format
	if updatedPayment.PayDateString != "" {
		dateStr := updatedPayment.PayDateString
		date, err := time.Parse("02-01-2006", dateStr)
		if err != nil {
			return model.PaymentRespont{}, err
		}
		updatedPayment.PayDate = date
	}

	paymentData.PartnerID = updatedPayment.PartnerID
	paymentData.Discount = updatedPayment.Discount
	paymentData.IsPrecentage = updatedPayment.IsPrecentage
	paymentData.BatchNo = updatedPayment.BatchNo
	paymentData.PayDate = updatedPayment.PayDate

	//validate before save
	paymentData, err = pr.BeforeSave(paymentData)
	if err != nil {
		return data, err
	}

	//validation docaction
	paymentData, err = pr.DocProcess(paymentData, string(updatedPayment.DocAction))
	if err != nil {
		return data, err
	}

	//save the data
	if err := pr.db.Omit("updated_by").Save(&paymentData).Error; err != nil {
		return data, err
	}

	data, err = pr.parsingPaymentToPaymentRespont(paymentData)
	if err != nil {
		return data, err
	}

	return data, nil
}

// Delete implements PaymentRepositoryinterface.
func (pr *paymentRepo) Delete(id uuid.UUID) (string, error) {
	data, err := pr.ShowInternal(id)
	batchNo := data.BatchNo
	if err != nil {
		return "", err
	}

	if err := pr.db.Delete(&data).Error; err != nil {
		return "", err
	}

	return batchNo, nil
}

func (pr *paymentRepo) BeforeSave(data model.Payment) (model.Payment, error) {
	//change grand total to sum of line first!
	var totalAmount float64
	//searching the value off grandtotal
	err := pr.db.Table("payment_lines").Where("payment_id = ?", data.ID).Select("COALESCE(SUM(amount), 0)").Scan(&totalAmount).Error

	if err != nil {
		return data, err
	}

	//set the value of grandtotal
	data.GrandTotal = totalAmount

	//run handling GrandTotal
	data = pr.handlingGrandTotal(data)

	return data, nil
}

// BeforeUpdate implements PaymentRepositoryinterface.
func (pr *paymentRepo) BeforeUpdate(data model.Payment, docaction string) (model.Payment, error) {
	// Open the JSON file
	file, err := os.Open("internal/repository/invoice/mapper.json")
	if err != nil {
		return model.Payment{}, err
	}
	defer file.Close() // Ensure file is closed

	// Read the file content
	mapperFile, err := io.ReadAll(file)
	if err != nil {
		return model.Payment{}, err
	}

	// Unmarshal the JSON data
	var mapper []mapperDoc
	err = json.Unmarshal(mapperFile, &mapper)
	if err != nil {
		return model.Payment{}, err
	}

	// Filter elements with "DR" status
	for _, element := range mapper {
		if element.Status == string(data.Status) && element.Doaction == docaction {
			return data, nil
		}
	}

	return model.Payment{}, errors.New("document invalid: status" + string(data.Status) + "docaction" + string(data.DocAction) + "NOT FOUND")
}

func (pr *paymentRepo) handlingGrandTotal(data model.Payment) model.Payment {
	if data.IsPrecentage {
		data.GrandTotal = data.GrandTotal - (data.GrandTotal * data.Discount / 100)
	} else {
		data.GrandTotal = data.GrandTotal - data.Discount
	}
	return data
}

func (pr *paymentRepo) parsingPaymentToPaymentRespont(payment model.Payment) (model.PaymentRespont, error) {
	data := model.PaymentRespont{}
	dataPreload, err := pr.ShowInternal(payment.UUID)
	if err != nil {
		return data, err
	}

	createdBy := model.UserPartial{
		UserId:   dataPreload.User.ID,
		Username: dataPreload.User.Username,
	}
	updateBy := model.UserPartial{
		UserId:   dataPreload.UserUpdated.ID,
		Username: dataPreload.UserUpdated.Username,
	}
	partner := model.PartnerPartialRespon{
		UUID: dataPreload.Partner.UUID,
		Name: dataPreload.Partner.Name,
	}

	line, err := pr.IndexLine(15, 0, dataPreload.ID)
	if err != nil {
		return model.PaymentRespont{}, err
	}

	data = model.PaymentRespont{
		ID:           dataPreload.UUID,
		DocumentNo:   dataPreload.DocumentNo,
		BatchNo:      dataPreload.BatchNo,
		IsPrecentage: dataPreload.IsPrecentage,
		Discount:     dataPreload.Discount,
		TotalLine:    dataPreload.TotalLine,
		GrandTotal:   dataPreload.GrandTotal,
		Status:       dataPreload.Status,
		DoAction:     dataPreload.DocAction,
		CreatedAt:    dataPreload.CreatedAt,
		UpdateAt:     dataPreload.UpdateAt,
		CreatedBy:    createdBy,
		UpdatedBy:    updateBy,
		Partner:      partner,
		UUID:         dataPreload.UUID,
		PayDate:      data.PayDate,
		Line:         line,
	}
	return data, nil
}

func (pr *paymentRepo) ParsingPaymentToPaymentRequest(payment model.Payment) (model.PaymentRequest, error) {
	data := model.PaymentRequest{
		UpdatedBy:    payment.CreatedBy,
		PartnerID:    payment.PartnerID,
		Discount:     payment.Discount,
		BatchNo:      payment.BatchNo,
		Status:       payment.Status,
		DocAction:    payment.DocAction,
		IsPrecentage: payment.IsPrecentage,
		PayDate:      payment.PayDate,
	}

	return data, nil
}

func (pr *paymentRepo) getTableName() string {
	return "payments"
}

func (pr *paymentRepo) HandlingPagination(q string, limit int, offset int, dateFrom string, dateTo string) (int64, error) {
	var count int64 = 0
	data := model.Payment{}
	//q param handler
	if q != "" {
		if err := pr.db.Joins("Partner", pr.db.Where(model.GetSearchParamPartnerV2(q))).Where(model.GetSeatchParamPayment(q, dateFrom, dateTo)).Find(&data).Count(&count).Error; err != nil {
			return count, err
		}
	} else {
		if err := pr.db.Order("created_at DESC").Limit(limit).Offset(offset).Find(&data).Count(&count).Error; err != nil {
			return count, err
		}
	}
	return count, nil
}
