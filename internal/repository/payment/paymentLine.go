package payment

import (
	"bemyfaktur/internal/model"
	"errors"

	"github.com/google/uuid"
)

func (pr *paymentRepo) CreateLine(paymentLine model.PaymentLineRequest) (model.PaymentLineRespont, error) {
	dataReturn := model.PaymentLineRespont{}
	data := model.PaymentLine{
		PaymentID:      paymentLine.PaymentID,
		Price:          paymentLine.Price,
		Discount:       paymentLine.Discount,
		InvoiceID:      paymentLine.Invoice_id,
		IsPrecentage:   paymentLine.IsPrecentage,
		CreatedBy:      paymentLine.CreatedBy,
		UpdatedBy:      paymentLine.UpdatedBy,
		OrganizationId: paymentLine.OrganizationId,
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

	if err := pr.db.Preload("Payment").Preload("Invoice").Preload("User").Preload("UserUpdated").Where(model.PaymentLine{PaymentID: paymentId}).Order("created_at").Limit(limit).Offset(offset).Find(&paymentLine).Error; err != nil {
		return data, err
	}

	for _, paymentline := range paymentLine {
		parsing, err := pr.parsingPaymentLineToRespont(paymentline)
		if err != nil {
			return []model.PaymentLineRespont{}, err
		}
		data = append(data, parsing)

	}

	return data, nil

}

// ShowLine implements PaymentRepositoryinterface.
func (pr *paymentRepo) ShowLine(id uuid.UUID) (model.PaymentLineRespont, error) {
	data := model.PaymentLine{}

	if err := pr.db.Preload("Payment").Preload("Invoice").Preload("User").Preload("UserUpdated").Where(model.PaymentLine{UUID: id}).First(&data).Error; err != nil {
		return model.PaymentLineRespont{}, errors.New("data not found")
	}

	return pr.parsingPaymentLineToRespont(data)
}

// ShowLine implements PaymentRepositoryinterface.
func (pr *paymentRepo) ShowLineInternal(id uuid.UUID) (model.PaymentLine, error) {
	data := model.PaymentLine{}

	if err := pr.db.Preload("Payment").Preload("Invoice").Preload("User").Preload("UserUpdated").Where(model.PaymentLine{UUID: id}).First(&data).Error; err != nil {
		return data, errors.New("data not found")
	}

	return data, nil
}

// UpdateLine implements PaymentRepositoryinterface.
func (pr *paymentRepo) UpdateLine(id uuid.UUID, updatedPaymentLine model.PaymentLineRequest) (model.PaymentLineRespont, error) {
	//set var
	data := model.PaymentLineRespont{}
	paymentLineData, err := pr.ShowLineInternal(id)

	if err != nil {
		return data, err
	}

	paymentLineData.Price = updatedPaymentLine.Price
	paymentLineData.InvoiceID = updatedPaymentLine.Invoice_id
	paymentLineData.IsPrecentage = updatedPaymentLine.IsPrecentage
	paymentLineData.Discount = updatedPaymentLine.Discount

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
func (pr *paymentRepo) DeleteLine(id uuid.UUID) (string, error) {
	data, err := pr.ShowLineInternal(id)

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
	dataPreload, err := pr.ShowLineInternal(paymentLine.UUID)
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
	payment := model.PaymentPartialRespont{
		UUID:       dataPreload.Payment.UUID,
		BatchNo:    dataPreload.Payment.BatchNo,
		DocumentNo: dataPreload.Payment.DocumentNo,
	}
	invoice := model.InvoicePartialRespont{
		UUID:              dataPreload.Invoice.UUID,
		BatchNo:           dataPreload.Invoice.BatchNo,
		Documentno:        dataPreload.Invoice.DocumentNo,
		OustandingPayment: dataPreload.Invoice.OustandingPayment,
	}

	data = model.PaymentLineRespont{
		ID:           dataPreload.UUID,
		Price:        dataPreload.Price,
		Amount:       dataPreload.Amount,
		BatchNo:      dataPreload.Invoice.BatchNo,
		Invoice_id:   dataPreload.InvoiceID,
		Discount:     dataPreload.Discount,
		IsPrecentage: dataPreload.IsPrecentage,
		Payment:      payment,
		CreatedBy:    createdBy,
		UpdatedBy:    updateBy,
		Invoice:      invoice,
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

	//after save update total line of payment header
	query = `
	update payments as p set total_line = (
		select coalesce(sum(amount), 0) from payment_lines pl where pl.payment_id = p.id) 
		where p.id = ?
	`

	if err := pr.db.Exec(query, data.PaymentID).Error; err != nil {
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

func (pr *paymentRepo) HandlingPaginationLine(q string, limit int, offset int, paymentID int) (int64, error) {
	var count int64 = 0
	data := model.Invoice{}
	//q param handler
	if q != "" {
		if err := pr.db.Joins("Partner", pr.db.Where(model.GetSearchParamPartnerV2(q))).Where(model.GetSeatchParamPaymentLine(q, paymentID)).Find(&data).Count(&count).Error; err != nil {
			return count, err
		}
	} else {
		if err := pr.db.Order("created_at DESC").Where(model.GetSeatchParamPaymentLine(q, paymentID)).Find(&data).Count(&count).Error; err != nil {
			return count, err
		}
	}
	return count, nil
}
