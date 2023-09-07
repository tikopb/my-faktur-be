package payment

import (
	"bemyfaktur/internal/model"
	"errors"
)

func (pr *paymentRepo) CreateLine(paymentLine model.PaymentLineRequest) (model.PaymentLineRespont, error) {
	dataReturn := model.PaymentLineRespont{}
	data := model.PaymentLine{
		PaymentID: paymentLine.PaymentID,
		Price:     paymentLine.Price,
		Amount:    paymentLine.Price,
		CreatedBy: "1", //##@ until the process security done deveop
		InvoiceID: paymentLine.Invoice_id,
	}

	//saved Data
	if err := pr.db.Create(&data).Error; err != nil {
		return dataReturn, err
	}

	//parsing data
	dataReturn, err := pr.parsingPaymentLineToRespont(data)
	if err != nil {
		return dataReturn, err
	}

	return dataReturn, nil
}

// IndexLine implements PaymentRepositoryinterface.
func (pr *paymentRepo) IndexLine(limit int, offset int, paymentId int) ([]model.PaymentLineRespont, error) {
	paymentLine := []model.PaymentLine{}
	data := []model.PaymentLineRespont{}

	if err := pr.db.Preload("Payment").Preload("Invoice").Where(model.PaymentLine{PaymentID: paymentId}).Order("created_at").Limit(limit).Offset(offset).Find(&paymentLine).Error; err != nil {
		return data, err
	}

	for _, paymentline := range paymentLine {
		indexResponse := model.PaymentLineRespont{
			ID:         paymentline.ID,
			Price:      paymentline.Price,
			Amount:     paymentline.Amount,
			BatchNo:    paymentline.Invoice.BatchNo,
			Payment:    paymentline.Payment,
			Invoice_id: paymentline.InvoiceID,
		}
		data = append(data, indexResponse)
	}

	return data, nil

}

// ShowLine implements PaymentRepositoryinterface.
func (pr *paymentRepo) ShowLine(id int) (model.PaymentLine, error) {
	data := model.PaymentLine{}

	if err := pr.db.Preload("Payment").Preload("Invoice").Preload("User").First(&data, id).Error; err != nil {
		return data, errors.New("daat not found")
	}

	return data, nil
}

// UpdateLine implements PaymentRepositoryinterface.
func (pr *paymentRepo) UpdateLine(id int, updatedPaymentLine model.PaymentLineRequest) (model.PaymentLineRespont, error) {
	//set var
	data := model.PaymentLineRespont{}
	paymentLineData, err := pr.ShowLine(id)

	if err != nil {
		return data, err
	}

	paymentLineData.Price = updatedPaymentLine.Price
	paymentLineData.InvoiceID = updatedPaymentLine.Invoice_id

	//save the data
	if err := pr.db.Save(&paymentLineData).Error; err != nil {
		return data, err
	}

	data, err = pr.parsingPaymentLineToRespont(paymentLineData)
	if err != nil {
		return data, err
	}

	return data, nil
}

// DeleteLine implements PaymentRepositoryinterface.
func (pr *paymentRepo) DeleteLine(id int) (string, error) {
	data, err := pr.ShowLine(id)

	if err != nil {
		return "", nil
	}

	if err := pr.db.Delete(&data).Error; err != nil {
		return "", err
	}

	return data.Payment.BatchNo, nil
}

func (pr *paymentRepo) parsingPaymentLineToRespont(paymentLine model.PaymentLine) (model.PaymentLineRespont, error) {
	data := model.PaymentLineRespont{}
	dataPreload, err := pr.ShowLine(paymentLine.ID)
	if err != nil {
		return data, err
	}

	data = model.PaymentLineRespont{
		ID:         dataPreload.ID,
		Price:      dataPreload.Price,
		Amount:     dataPreload.Amount,
		BatchNo:    dataPreload.Invoice.BatchNo,
		Invoice_id: dataPreload.InvoiceID,
	}

	return data, nil
}
