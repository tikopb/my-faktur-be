package payment

import (
	"bemyfaktur/internal/model"
	"bemyfaktur/internal/model/constant"
)

// DocProcess implements PaymentRepositoryinterface.
func (pr *paymentRepo) DocProcess(data model.Payment) error {
	var err error

	switch data.DocAction {
	case constant.PaymentDocActionComplete:
		err = pr.CompleteIT(data)
	case constant.PaymentDocActionProcessed:
		err = nil
	case constant.PaymentDocActionVoid:
		err = pr.ReversedIt(data)
	}

	return err
}

// CompleteIT implements PaymentRepositoryinterface.
func (pr *paymentRepo) CompleteIT(data model.Payment) error {
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
		return err
	}

	return nil
}

// ReversedIt implements PaymentRepositoryinterface.
func (pr *paymentRepo) ReversedIt(data model.Payment) error {
	return nil
}
