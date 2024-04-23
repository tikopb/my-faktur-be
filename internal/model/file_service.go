package model

import (
	"mime/multipart"
	"time"

	"github.com/google/uuid"
)

type FileService struct {
	ID          int       `gorm:"primaryKey;autoIncrement"`
	UUID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();index:idx_docService_uuid"`
	CreatedAt   time.Time `gorm:"column:created_at;default:current_timestamp"`
	UpdateAt    time.Time `gorm:"column:updated_at;default:current_timestamp"`
	CreatedBy   string    `gorm:"column:created_by;" json:"created_by"`
	User        User      `gorm:"foreignKey:created_by"`
	UpdatedBy   string    `gorm:"column:updated_by" json:"updated_by"`
	UserUpdated User      `gorm:"foreignKey:updated_by"`
	FileName    string    `gorm:"column:file_name;unique;index:idx_docService_filename"`
	UuidDoc     uuid.UUID `gorm:"column:uuid_doc;index:idx_docService_uuidDoc"`
	DocType     string    `gorm:"column:doctype;index:idx_docService_docType"`
}

type FileServiceRequest struct {
	File      *multipart.FileHeader `form:"file"`
	File64    []byte                `json:"file64"`
	UuidDoc   uuid.UUID             `json:"uuid_doc"`
	DocType   string                `json:"doctype"`
	FileName  string                `json:"filename"`
	CreatedBy string                `json:"-"`
}
type FileServiceRespont struct {
	File64   string                `json:"file64"`
	File     *multipart.FileHeader `form:"file"`
	FileName string                `json:"filename"`
}
