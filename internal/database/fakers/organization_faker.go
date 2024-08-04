package fakers

import (
	"bemyfaktur/internal/model"
	"gorm.io/gorm"
)

func OrgFaker(db *gorm.DB, org model.Organization) *model.Organization {
	return &model.Organization{
		ID:          org.ID,
		UUID:        org.UUID,
		CreatedAt:   org.CreatedAt,
		UpdateAt:    org.UpdateAt,
		CreatedBy:   org.CreatedBy,
		UpdatedBy:   org.UpdatedBy,
		OrgCode:     org.OrgCode,
		Name:        org.Name,
		Description: org.Description,
		IsActive:    org.IsActive,
	}
}
