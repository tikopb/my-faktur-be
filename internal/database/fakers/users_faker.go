package fakers

import (
	"bemyfaktur/internal/model"
	"time"

	"gorm.io/gorm"
)

func UserFacker(db *gorm.DB) *model.User {
	return &model.User{
		ID:             "0bdfa34b-cdc4-4c4b-b087-f6b6d7cf81d5",
		Username:       "Sam_Will",
		Hash:           "$argon2id$v=19$m=65536,t=1,p=32$UzafyF+xamjTQCcxico7Lw$AAAAAAAAAAAAAAAA/a0578kNnOq/jwmjcyJdSm8JQTI=",
		FullName:       "Sam Will",
		CreatedAt:      time.Time{},
		IsActive:       true,
		OrganizationId: 1,
	}
}
