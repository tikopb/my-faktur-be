package user

import (
	"bemyfaktur/internal/model"
	"errors"

	"gorm.io/gorm"
)

type userRepo struct {
	db *gorm.DB
}

func GetRepository(db *gorm.DB) Repository {
	return &userRepo{
		db: db,
	}
}

// Create implements Repository.
func (*userRepo) Create(user model.User) (model.User, error) {
	panic("unimplemented")
}

// Delete implements Repository.
func (*userRepo) Delete(id string) (string, error) {
	panic("unimplemented")
}

// Index implements Repository.
func (*userRepo) Index(limit int, offset int) ([]model.User, error) {
	panic("unimplemented")
}

// Show implements Repository.
func (ur *userRepo) Show(id string) (model.User, error) {
	data := model.User{}

	if err := ur.db.Where(model.User{ID: id}).Preload("User").First(&data).Error; err != nil {
		return data, errors.New("data not found")
	}

	return data, nil
}

// Update implements Repository.
func (*userRepo) Update(id string, updateduser model.User) (model.User, error) {
	panic("unimplemented")
}
