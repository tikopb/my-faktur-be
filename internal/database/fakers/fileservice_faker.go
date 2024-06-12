package fakers

import (
	"bemyfaktur/internal/model"

	"gorm.io/gorm"
)

func FileServiceFakers(db *gorm.DB, fileService model.FileService) *model.FileService {
	return &model.FileService{
		UUID:      fileService.UUID,
		CreatedAt: fileService.CreatedAt,
		UpdateAt:  fileService.UpdateAt,
		CreatedBy: fileService.UpdatedBy,
		UpdatedBy: fileService.UpdatedBy,
		FileName:  fileService.FileName,
		UuidDoc:   fileService.UuidDoc,
		DocType:   fileService.DocType,
	}
}
