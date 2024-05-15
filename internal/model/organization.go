package model

import (
	"time"

	"github.com/google/uuid"
)

type Organization struct {
	ID          int       `json:"-" gorm:"primaryKey;autoIncrement"`
	UUID        uuid.UUID `json:"id" gorm:"type:uuid;default:uuid_generate_v4();index:idx_invoice_uuid"`
	CreatedAt   time.Time `gorm:"column:created_at;default:current_timestamp"`
	UpdateAt    time.Time `gorm:"column:updated_at;default:current_timestamp"`
	CreatedBy   string    `gorm:"column:created_by;" json:"created_by"`
	User        User      `gorm:"foreignKey:created_by"`
	UpdatedBy   string    `gorm:"column:updated_by" json:"updated_by"`
	UserUpdated User      `gorm:"foreignKey:updated_by"`
	OrgCode     string    `gorm:"column:org_code;not null"`
	Name        string    `gorm:"name;not null"`
	Description string    `gorm:"description"`
	IsActive    bool      `gorm:"isactive;default: true"`
}

type OrganizationRequest struct {
	OrgCode     string `json:"org_code"`
	Name        string `json:"name"`
	Description string `json:"description"`
	IsActive    bool   `json:"isactive"`
}

type OrganizationRespont struct {
	ID          uuid.UUID `json:"ID"`
	OrgCode     string    `json:"org_code"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	IsActive    bool      `json:"isactive"`
}
