package model

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

type Partner struct {
	ID        int       `gorm:"primaryKey;autoIncrement" json:"-"`
	Name      string    `gorm:"column:name;not null;index:idx_partner_name" json:"name" `
	CreatedAt time.Time `gorm:"column:created_at;default:current_timestamp"`
	CreatedBy string    `gorm:"column:created_by;index:idx_partner_created_by" json:"created_by"`
	User      User      `gorm:"foreignKey:created_by"`
	DNAmount  float64   `gorm:"column:dn_amount" json:"dn_amount"`
	CNAmount  float64   `gorm:"column:cn_amount" json:"cn_amount"`
	Isactive  bool      `gorm:"column:isactive;index:idx_partner_isactive;default: true" json:"isactive"`
	Code      string    `gorm:"column:bp_code;unique;not null;index:idx_partner_code" json:"bpcode"`
	UUID      uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();index:idx_partner_uuid"`
}
type PartnerRespon struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Code      string    `json:"bp_code"`
	CreatedAt time.Time `json:"created_at"`
	DNAmount  float64   `json:"dn_amount"`
	CNAmount  float64   `json:"cn_amount"`
	Isactive  bool      `json:"isactive"`
	CreatedBy string    `json:"created_by"`
}

type PartnerPartialRespon struct {
	UUID uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

func GetSeatchParamPartner() []string {
	searchParam := []string{"name", "bp_code"}
	return searchParam
}

func GetSearchParamPartnerV2(q string) string {
	q = strings.ToLower(q)
	q = "'%" + q + "%'"
	value := " lower(name)  LIKE " + q + " OR lower(bp_code) LIKE " + q

	return value
}
