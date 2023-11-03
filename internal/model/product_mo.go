package model

import "time"

type Product struct {
	ID          int       `json:"id" gorm:"primaryKey;autoIncrement"`
	Name        string    `json:"name" gorm:"column:name;unique;not null"`
	Description string    `json:"description" gorm:"column:description"`
	CreatedBy   string    `gorm:"column:created_by" json:"created_by"`
	CreatedAt   time.Time `gorm:"column:created_at" json:"created_at"`
	User        User      `gorm:"foreignKey:created_by"`
	IsActive    bool      `gorm:"column:isactive"`
}
type ProductRespon struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	IsActive    bool      `json:"isactive"`
	CreatedBy   string    `json:"createdby"`
	CreatedAt   time.Time `json:"created_at"`
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
