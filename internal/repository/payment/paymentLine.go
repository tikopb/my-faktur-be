package payment

import (
	"bemyfaktur/internal/model"
	"errors"
)

func (pr *paymentRepo) CreateLine(paymentLine model.PaymentLineRequest) (model.PaymentLineRespont, error) {
	dataReturn := model.PaymentLineRespont{}
	data := model.PaymentLine{
		PaymentID:    paymentLine.PaymentID,
		Price:        paymentLine.Price,
		Discount:     paymentLine.Discount,
		CreatedBy:    "1", //##@ until the process security done deveop
		InvoiceID:    paymentLine.Invoice_id,
		IsPrecentage: paymentLine.IsPrecentage,
	}

	//beforesave validation
	data, err := pr.beforeSave(data)
	if err != nil {
		return dataReturn, err
	}

	//saved Data
	if err := pr.db.Create(&data).Error; err != nil {
		return dataReturn, err
	}

	//aftersave validation
	err = pr.afterSave(data)
	if err != nil {
		return dataReturn, err
	}

	//parsing data
	dataReturn, err = pr.parsingPaymentLineToRespont(data)
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
			ID:           paymentline.ID,
			BatchNo:      paymentline.Invoice.BatchNo,
			Invoice_id:   paymentline.InvoiceID,
			Price:        paymentline.Price,
			Discount:     paymentline.Discount,
			IsPrecentage: paymentline.IsPrecentage,
			Amount:       paymentline.Amount,
			Payment:      paymentline.Payment,
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

	//aftersave validation
	err = pr.afterSave(paymentLineData)
	if err != nil {
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

	//aftersave validation
	err = pr.afterSave(data)
	if err != nil {
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
		ID:           dataPreload.ID,
		Price:        dataPreload.Price,
		Amount:       dataPreload.Amount,
		BatchNo:      dataPreload.Invoice.BatchNo,
		Invoice_id:   dataPreload.InvoiceID,
		Discount:     data.Discount,
		IsPrecentage: data.IsPrecentage,
		Payment:      data.Payment,
	}

	return data, nil
}

func (pr *paymentRepo) beforeSave(data model.PaymentLine) (model.PaymentLine, error) {

	//validate price cant more than oustanding
	var outstandingPayment float64
	query := `
		select oustanding_payment from invoices i where id =?
	`
	if err := pr.db.Raw(query, data.InvoiceID).Scan(&outstandingPayment).Error; err != nil {
		return data, err
	}
	if data.Price > outstandingPayment {
		return data, errors.New("payment tidak boleh melebihi nilai oustanding")
	}

	//change amount validartion
	if data.IsPrecentage {
		data.Amount = data.Price - (data.Price * data.Discount / 100)
	} else {
		data.Amount = data.Price - data.Discount
	}

	return data, nil
}

func (pr *paymentRepo) afterSave(data model.PaymentLine) error {

	// after save update the oustanding payment invoice
	query := `
	update invoices as i
	SET oustanding_payment = (i.grand_total - COALESCE((SELECT SUM(price) FROM payment_lines AS pl WHERE pl.invoice_id = i.id), 0))
	where i.id = ?`

	if err := pr.db.Exec(query, data.InvoiceID).Error; err != nil {
		return err
	}

	//after save update grand total of payment header
	query = `
	update payments as p set grand_total = (
		select coalesce(sum(amount), 0) from payment_lines pl where pl.payment_id = p.id) 
		where p.id = ?
	`

	if err := pr.db.Exec(query, data.PaymentID).Error; err != nil {
		return err
	}

	return nil
}
