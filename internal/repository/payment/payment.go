package payment

import (
	"bemyfaktur/internal/model"
	documentutil "bemyfaktur/internal/model/documentUtil"
	pgUtil "bemyfaktur/internal/model/paginationUtil"
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
		if err := pr.db.Preload("Partner").Preload("User").Preload("UserUpdated").Joins("Partner", pr.db.Where(model.GetSearchParamPartnerV2(q))).Where(model.GetSeatchParamPayment(q, dateFrom, dateTo)).Limit(limit).Offset(offset).Order(orderParam).Find(&data).Error; err != nil {
			return dataReturn, err
		}
	} else {
		if err := pr.db.Preload("Partner").Preload("User").Preload("UserUpdated").Order("created_at DESC").Where(model.GetSeatchParamPayment(q, dateFrom, dateTo)).Limit(limit).Offset(offset).Find(&data).Error; err != nil {
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

	paymentData := model.Payment{
		CreatedBy:  payment.CreatedBy,
		PartnerID:  payment.PartnerID,
		GrandTotal: 0,
		Discount:   0,
		BatchNo:    payment.BatchNo,
		DocumentNo: documentno,
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

// Show implements PaymentRepositoryinterface.
func (pr *paymentRepo) Show(id uuid.UUID) (model.PaymentRespont, error) {
	data := model.PaymentRespont{}

	if err := pr.db.Preload("Partner").Preload("User").Preload("UserUpdated").Where(&model.Payment{UUID: id}).First(&data).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.PaymentRespont{}, errors.New("data not found")
		}
	}
	return data, nil
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

	paymentData.PartnerID = updatedPayment.PartnerID
	paymentData.Discount = updatedPayment.Discount
	paymentData.IsPrecentage = updatedPayment.IsPrecentage
	paymentData.BatchNo = updatedPayment.BatchNo

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
	if err := pr.db.Save(&paymentData).Error; err != nil {
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
	data, err := pr.Show(id)
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
	var grandTotal float64 = 0
	query := `
    	select coalesce(sum(amount), 0) from payment_lines pl where payment_id = ?
	`
	if err := pr.db.Raw(query, data.ID).Scan(&grandTotal).Error; err != nil {
		return data, err
	}
	data.GrandTotal = grandTotal

	//run handling GrandTotal
	data = pr.handlingGrandTotal(data)

	return data, nil
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

	data = model.PaymentRespont{
		ID:           dataPreload.UUID,
		GrandTotal:   dataPreload.GrandTotal,
		Discount:     dataPreload.Discount,
		BatchNo:      dataPreload.BatchNo,
		Status:       dataPreload.Status,
		DocumentNo:   dataPreload.DocumentNo,
		DoAction:     dataPreload.DocAction,
		IsPrecentage: data.IsPrecentage,
		CreatedBy:    createdBy,
		UpdatedBy:    updateBy,
		Partner:      partner,
	}
	return data, nil
}

func (pr *paymentRepo) getTableName() string {
	return "payments"
}

func (pr *paymentRepo) HandlingPagination(q string, limit int, offset int, dateFrom string, dateTo string) (int64, error) {
	var count int64 = 0
	data := model.Invoice{}
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
