package fakers

import (
	"bemyfaktur/internal/model"
	"time"

	"gorm.io/gorm"
)

func OrganizationFakers(db *gorm.DB, org model.Organization) *model.Organization {
	return &model.Organization{
		ID:          org.ID,
		UUID:        org.UUID,
		CreatedAt:   time.Now(),
		UpdateAt:    time.Now(),
		CreatedBy:   org.CreatedBy,
		UpdatedBy:   org.UpdatedBy,
		OrgCode:     org.OrgCode,
		Name:        org.Name,
		Description: org.Description,
		IsActive:    org.IsActive,
	}
}
