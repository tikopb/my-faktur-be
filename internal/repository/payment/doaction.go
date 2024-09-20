package payment

import (
	"bemyfaktur/internal/model"
	"bemyfaktur/internal/model/constant"
)

// DocProcess implements PaymentRepositoryinterface.
func (pr *paymentRepo) DocProcess(data model.Payment, docaction string) (model.Payment, error) {
	var err error

	//if nothing change then dont do anything
	if data.Status == constant.PaymentStatus(docaction) {
		return data, nil
	}

	switch docaction {
	case "CO":
		data, err = pr.CompleteIT(data, docaction)
	case "IP":
		data.Status = constant.PaymentStatusProcessed
		data.DocAction = constant.PaymentDocActionProcessed
		err = nil
	case "VO":
		data, err = pr.ReversedIt(data, docaction)
	}

	return data, err
}

// CompleteIT implements PaymentRepositoryinterface.
func (pr *paymentRepo) CompleteIT(data model.Payment, docaction string) (model.Payment, error) {
	var err error

	query := `UPDATE invoices
		SET oustanding_payment = (
			SELECT SUM(il.amount)
			FROM payment_lines il
			WHERE il.invoice_id = invoices.id
		)
		WHERE invoices.id IN (
			SELECT il.invoice_id
			FROM payment_lines il
			WHERE il.payment_id = ?
		);
	`

	if err = pr.db.Raw(query, data.ID).Error; err != nil {
		return data, err
	}

	data.Status = constant.PaymentStatus(docaction)
	data.DocAction = constant.PaymentDocAction(docaction)

	return data, nil
}

// ReversedIt implements PaymentRepositoryinterface.
func (pr *paymentRepo) ReversedIt(data model.Payment, docaction string) (model.Payment, error) {

	//update line
	query := `
	update payment_lines set price = 0, amount = 0, discount = 0 where payment_id  = ? 
	`
	if err := pr.db.Raw(query, data.ID).Error; err != nil {
		return data, err
	}

	//update header
	data.GrandTotal = 0
	data.Discount = 0

	//update oustanding
	query = `UPDATE invoices
		SET oustanding_payment = (
			SELECT SUM(il.amount)
			FROM payment_lines il
			WHERE il.invoice_id = invoices.id
		)
		WHERE invoices.id IN (
			SELECT il.invoice_id
			FROM payment_lines il
			WHERE il.payment_id = ?
		);
	`

	if err := pr.db.Raw(query, data.ID).Error; err != nil {
		return data, err
	}

	//change status
	data.Status = constant.PaymentStatus(docaction)
	data.DocAction = constant.PaymentDocAction(docaction)

	return data, nil
}
