package payment

import (
	"bemyfaktur/internal/model"
	documentutil "bemyfaktur/internal/model/documentUtil"
	pgUtil "bemyfaktur/internal/model/paginationUtil"
	"errors"

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
func (pr *paymentRepo) Index(limit int, offset int, q string) ([]model.PaymentRespont, error) {
	data := []model.Payment{}
	dataReturn := []model.PaymentRespont{}

	//q param handler
	if q != "" {
		if err := pr.db.Joins("Partner", pr.db.Where(model.GetSeatchParamPartnerV2(q))).Where(model.GetSeatchParamPayment(q)).Limit(limit).Offset(offset).Find(&data).Error; err != nil {
			return dataReturn, err
		}
	} else {
		if err := pr.db.Order("created_at DESC").Limit(limit).Offset(offset).Find(&data).Error; err != nil {
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
			CreatedBy:    dataPreload.CreatedBy,
			PartnerID:    dataPreload.PartnerID,
			Partner_name: dataPreload.Partner_name,
			GrandTotal:   dataPreload.GrandTotal,
			Discount:     dataPreload.Discount,
			BatchNo:      dataPreload.BatchNo,
			Status:       dataPreload.Status,
			DoAction:     dataPreload.DoAction,
			Partner:      dataPreload.Partner,
			DocumentNo:   dataPreload.DocumentNo,
			IsPrecentage: dataPreload.IsPrecentage,
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
		CreatedBy:  "1", //##@ UNTIL SECURIT model DONE!
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
func (pr *paymentRepo) Show(id int) (model.Payment, error) {
	data := model.Payment{}

	if err := pr.db.Preload("Partner").Preload("User").First(&data, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return data, errors.New("data not found")
		}
	}
	return data, nil
}

// Update implements PaymentRepositoryinterface.
func (pr *paymentRepo) Update(id int, updatedPayment model.PaymentRequest) (model.PaymentRespont, error) {
	//set var
	data := model.PaymentRespont{}
	paymentData, err := pr.Show(id) //get payment data

	if err != nil {
		return data, err
	}

	paymentData.PartnerID = updatedPayment.PartnerID
	//paymentData.GrandTotal = updatedPayment.GrandTotal
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
func (pr *paymentRepo) Delete(id int) (string, error) {
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
	dataPreload, err := pr.Show(payment.ID)
	if err != nil {
		return data, err
	}

	data = model.PaymentRespont{
		ID:           dataPreload.ID,
		CreatedBy:    dataPreload.User.Username,
		PartnerID:    dataPreload.PartnerID,
		GrandTotal:   dataPreload.GrandTotal,
		Discount:     dataPreload.Discount,
		BatchNo:      dataPreload.BatchNo,
		Status:       dataPreload.Status,
		Partner:      dataPreload.Partner,
		DocumentNo:   dataPreload.DocumentNo,
		DoAction:     dataPreload.DocAction,
		IsPrecentage: data.IsPrecentage,
		Partner_name: data.Partner.Name,
	}
	return data, nil
}

func (pr *paymentRepo) getTableName() string {
	return "payments"
}

func (pr *paymentRepo) HandlingPagination(q string, limit int, offset int) (int64, error) {
	var count int64 = 0
	data := model.Invoice{}
	//q param handler
	if q != "" {
		if err := pr.db.Joins("Partner", pr.db.Where(model.GetSeatchParamPartnerV2(q))).Where(model.GetSeatchParamPayment(q)).Find(&data).Count(&count).Error; err != nil {
			return count, err
		}
	} else {
		if err := pr.db.Order("created_at DESC").Limit(limit).Offset(offset).Find(&data).Count(&count).Error; err != nil {
			return count, err
		}
	}
	return count, nil
}
