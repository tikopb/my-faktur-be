package constant

type PaymentStatus string

const (
	PaymentStatusComplete  PaymentStatus = "CO"
	PaymentStatusProcessed PaymentStatus = "processed"
	PaymentStatusVoid      PaymentStatus = "vo"
)
