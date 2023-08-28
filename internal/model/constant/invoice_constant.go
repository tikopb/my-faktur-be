package constant

type InvoiceStatus string

const (
	InvoiceStatusDraft     InvoiceStatus = "DR"
	InvoiceStatusProcessed InvoiceStatus = "IP"
	InvoiceStatusComplete  InvoiceStatus = "CO"
	InvoiceStatusVoid      InvoiceStatus = "VO"
)
