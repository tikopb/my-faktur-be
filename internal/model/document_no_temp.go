package model

import (
	"time"
)

type DocumentNoTemp struct {
	ID        int       `json:"id" gorm:"primaryKey;autoIncrement"`
	Prefix    string    `json:"prefix" gorm:"prefix;not null"`
	Suffix    string    `json:"suffix" gorm:"suffix;not null"`
	TableName string    `json:"tableName" gorm:"table_name;not null"`
	Counting  int       `json:"counting" gorm:"counting;default:0"`
	StartDate time.Time `json:"startDate" gorm:"start_date;not null"`
	EndDate   time.Time `json:"endDate" gorm:"end_date;not null"`
	CreatedAt time.Time `json:"createdAt" gorm:"column:created_at;default:CURRENT_TIMESTAMP"`
}
