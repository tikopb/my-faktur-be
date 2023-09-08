package invoice

import (
	"bemyfaktur/internal/model"
	"bemyfaktur/internal/model/constant"
	"errors"
	"fmt"
)

func (ir *invoiceRepo) DocProcess(data model.Invoice) error {
	var err error

	switch data.DocAction {
	case constant.InvoiceActionComplete:
		err = ir.CompleteIT(data)
	case constant.InvoiceActionProcessed:
		fmt.Println("roolback")
	case constant.InvoiceActionVoid:
		err = ir.ReversedIt(data)
	}

	return err
}

func (ir *invoiceRepo) CompleteIT(data model.Invoice) error {
	return nil
}

func (ir *invoiceRepo) ReversedIt(data model.Invoice) error {
	payment := []model.Payment{}

	query := `
		select i.documentno as documentno from payment_lines pl 
		join payments p on pl.invoice_id = p.id 
		join invoices i on pl.invoice_id = i.id 
		where pl.invoice_id = ? or i.status = ?
	`

	if err := ir.db.Raw(query, data.ID, constant.PaymentDocActionVoid).Scan(&payment).Error; err != nil {
		return err
	}

	// Check if there is more than one payment record.
	if len(payment) > 1 {
		errorMsg := ("Multiple payment records exist with DocumentNo: ")
		for _, paymentData := range payment {
			errorMsg += paymentData.DocumentNo + ", "
		}
		errorMsg = errorMsg[:len(errorMsg)-2]

		return errors.New(errorMsg)
	}

	return nil
}
