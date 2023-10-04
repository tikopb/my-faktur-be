package model

import "time"

type Partner struct {
	ID        int       `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `gorm:"column:created_at;default:current_timestamp"`
	CreatedBy string    `gorm:"column:created_by" json:"created_by"`
	User      User      `gorm:"foreignKey:created_by"`
	DNAmount  float64   `gorm:"column:dn_amount" json:"dn_amount"`
	CNAmount  float64   `gorm:"column:cn_amount" json:"cn_amount"`
	Isactive  bool      `gorm:"column:isactive" json:"isactive"`
	Code      string    `gorm:"column:bp_code;unique;not null" json:"bpcode"`
	Invoice   []Invoice
}

type PartnerRespon struct {
	Name     string  `json:"name"`
	DNAmount float64 `json:"dn_amount"`
	CNAmount float64 `json:"cn_amount"`
	Isactive bool    `json:"isactive"`
}

func GetSeatchParamPartner() []string {
	searchParam := []string{"name", "bp_code"}
	return searchParam
}

func GetSeatchParamPartnerV2(q string) string {
	q = "'%" + q + "%'"
	value := " lower(name)  LIKE " + q + " OR lower(bp_code) LIKE " + q

	return value
}
