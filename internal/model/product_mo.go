package model

import (
	"time"

	"github.com/google/uuid"
)

type Product struct {
	ID          int       `json:"-" gorm:"primaryKey;autoIncrement" `
	Name        string    `json:"name" gorm:"column:name;unique;not null;index:idx_product_name"`
	Value       string    `json:"value" gorm:"column:value;not null;index:idx_product_value"`
	Upc         string    `json:"upc" gorm:"column:upc;not null;index:idx_product_upc"`
	Description string    `json:"description" gorm:"column:description"`
	CreatedBy   string    `gorm:"column:created_by;index:idx_product_created_by" json:"created_by"`
	CreatedAt   time.Time `gorm:"column:created_at;default:current_timestamp" json:"created_at"`
	UpdateAt    time.Time `gorm:"column:updated_at;default:current_timestamp" json:"updated_at"`
	User        User      `json:"user" gorm:"foreignKey:created_by"`
	IsActive    bool      `json:"isactive" gorm:"column:isactive"`
	UUID        uuid.UUID `json:"id" gorm:"type:uuid;default:uuid_generate_v4();index:idx_product_uuid"`
}

type ProductRespon struct {
	UUID        uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	IsActive    bool      `json:"isactive"`
	CreatedBy   string    `json:"createdby"`
	CreatedAt   time.Time `json:"created_at"`
	UpdateAt    time.Time `json:"updated_at"`
	Value       string    `json:"value"`
	Upc         string    `json:"upc"`
}

func GetSeatchParamProduct() []string {
	searchParam := []string{"name", "description"}
	return searchParam
}

// searching for join table with other model
func GetSeatchParamProductV2(q string) string {
	value := " lower(name)  LIKE " + q + " OR lower(description) LIKE " + q
	return value
}
