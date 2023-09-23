package constant

type PaymentStatus string

const (
	PaymentStatusDraft     PaymentStatus = "DR"
	PaymentStatusComplete  PaymentStatus = "CO"
	PaymentStatusProcessed PaymentStatus = "IP"
	PaymentStatusVoid      PaymentStatus = "VO"
)

type PaymentDocAction string

const (
	PaymentDocActionDraft     PaymentDocAction = "DR"
	PaymentDocActionComplete  PaymentDocAction = "CO"
	PaymentDocActionProcessed PaymentDocAction = "IP"
	PaymentDocActionVoid      PaymentDocAction = "VO"
)
