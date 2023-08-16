package model

type Product struct {
	ID          int    `json:"id" gorm:"primaryKey;autoIncrement"`
	Name        string `json:"name" gorm:"column:name"`
	Description string `json:"description" gorm:"column:description"`
	CreatedBy   string `gorm:"column:created_by" json:"created_by"`
	User        User   `gorm:"foreignKey:created_by"`
	IsActive    bool   `gorm:"column:isactive"`
}
