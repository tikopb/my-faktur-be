package constant

type InvoiceStatus string

const (
	InvoiceStatusDraft     InvoiceStatus = "DR"
	InvoiceStatusProcessed InvoiceStatus = "IP"
	InvoiceStatusComplete  InvoiceStatus = "CO"
	InvoiceStatusVoid      InvoiceStatus = "VO"
)

type InvoiceDocAction string

const (
	InvoiceActionDraft     InvoiceDocAction = "DR"
	InvoiceActionProcessed InvoiceDocAction = "IP"
	InvoiceActionComplete  InvoiceDocAction = "CO"
	InvoiceActionVoid      InvoiceDocAction = "VO"
)
