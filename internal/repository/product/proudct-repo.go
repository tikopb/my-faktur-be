package product

import (
	"bemyfaktur/internal/model"
	"errors"

	"gorm.io/gorm"
)

type productRepo struct {
	db *gorm.DB
}

func GetRepository(db *gorm.DB) Repository {
	return &productRepo{
		db: db,
	}
}

// Create implements Repository.
func (pr *productRepo) Create(product model.Product) (model.ProductRespon, error) {
	data := model.ProductRespon{}
	if err := pr.db.Create(&product).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return data, errors.New("duplicatet data")
		}
		return data, err
	}

	data = model.ProductRespon{
		Name:        product.Name,
		Description: product.Description,
		IsActive:    product.IsActive,
	}

	return data, nil
}

// Index implements Repository.
func (pr *productRepo) Index(limit int, offset int) ([]model.Product, error) {
	data := []model.Product{}

	if err := pr.db.Order("name").Limit(limit).Offset(offset).Find(&data).Error; err != nil {
		return data, err
	}
	return data, nil
}

// Show implements Repository.
func (pr *productRepo) Show(id int) (model.Product, error) {
	var data model.Product

	if err := pr.db.First(&data, id).Preload("Product").Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return data, errors.New("data not found")
		}
	}
	return data, nil
}

// Update implements Repository.
func (pr *productRepo) Update(id int, updatedProduct model.Product) (model.ProductRespon, error) {
	dataUpdated := model.ProductRespon{}
	data, err := pr.Show(id)

	if err != nil {
		return dataUpdated, err
	}

	//slicing data update
	data.Name = updatedProduct.Name
	data.Description = updatedProduct.Description
	data.IsActive = updatedProduct.User.IsActive

	//inisiate data udpated system
	if err := pr.db.Save(&data).Error; err != nil {
		return dataUpdated, err
	}

	//inisiate data update version
	dataUpdated = model.ProductRespon{
		Name:        data.Name,
		Description: data.Name,
		IsActive:    data.IsActive,
	}

	return dataUpdated, nil
}

// Delete implements Repository.
func (pr *productRepo) Delete(id int) (string, error) {
	data, err := pr.Show(id)
	name := data.Name

	if err != nil {
		return "", err
	}
	if err := pr.db.Delete(&data).Error; err != nil {
		return "", err
	}
	return name, nil
}
