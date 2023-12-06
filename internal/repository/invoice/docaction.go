package invoice

import (
	"bemyfaktur/internal/model"
	"bemyfaktur/internal/model/constant"
	"errors"
)

func (ir *invoiceRepo) DocProcess(data model.Invoice, docaction string) (model.Invoice, error) {
	var err error

	//if nothing change then dont do anything
	if data.Status == constant.InvoiceStatus(docaction) {
		return data, nil
	}

	switch docaction {
	case "CO":
		data, err = ir.CompleteIT(data, docaction)
	case "IP":
		data.DocAction = constant.InvoiceActionProcessed
		data.Status = constant.InvoiceStatusProcessed
		err = nil
	case "VO":
		data, err = ir.ReversedIt(data, docaction)
	}

	return data, err
}

func (ir *invoiceRepo) CompleteIT(data model.Invoice, docaction string) (model.Invoice, error) {
	data.Status = constant.InvoiceStatus(docaction)
	data.DocAction = constant.InvoiceDocAction(docaction)
	data.OustandingPayment = data.GrandTotal
	return data, nil
}

func (ir *invoiceRepo) ReversedIt(data model.Invoice, docaction string) (model.Invoice, error) {
	payment := []model.Payment{}
	data = model.Invoice{}

	query := `
		select i.documentno as documentno from payment_lines pl 
		join payments p on pl.invoice_id = p.id 
		join invoices i on pl.invoice_id = i.id 
		where pl.invoice_id = ? or i.status = ?
	`

	if err := ir.db.Raw(query, data.ID, constant.PaymentDocActionVoid).Scan(&payment).Error; err != nil {
		return data, err
	}

	// Check if there is more than one payment record.
	if len(payment) > 1 {
		errorMsg := ("Multiple payment records exist with DocumentNo: ")
		for _, paymentData := range payment {
			errorMsg += paymentData.DocumentNo + ", "
		}
		errorMsg = errorMsg[:len(errorMsg)-2]

		return data, errors.New(errorMsg)
	}

	return data, nil
}
