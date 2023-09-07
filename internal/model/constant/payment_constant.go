package constant

type PaymentStatus string

const (
	PaymentStatusDraft     PaymentStatus = "DR"
	PaymentStatusComplete  PaymentStatus = "CO"
	PaymentStatusProcessed PaymentStatus = "IP"
	PaymentStatusVoid      PaymentStatus = "VO"
)
