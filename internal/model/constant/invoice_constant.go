package constant

type InvoiceStatus string

const (
	InvoiceStatusComplete  InvoiceStatus = "CO"
	InvoiceStatusProcessed InvoiceStatus = "processed"
	InvoiceStatusVoid      InvoiceStatus = "vo"
)
