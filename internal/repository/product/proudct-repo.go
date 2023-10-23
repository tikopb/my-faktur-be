package product

import (
	"bemyfaktur/internal/model"
	pgUtil "bemyfaktur/internal/model/paginationUtil"
	"errors"

	"gorm.io/gorm"
)

type productRepo struct {
	db         *gorm.DB
	pgUtilRepo pgUtil.Repository
}

func GetRepository(db *gorm.DB, pgRepo pgUtil.Repository) Repository {
	return &productRepo{
		db:         db,
		pgUtilRepo: pgRepo,
	}
}

// Create implements Repository.
func (pr *productRepo) Create(product model.Product) (model.ProductRespon, error) {
	data := model.ProductRespon{}
	if err := pr.db.Create(&product).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return data, errors.New("duplicate data")
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
func (pr *productRepo) Index(limit int, offset int, q string) ([]model.Product, error) {
	data := []model.Product{}

	if q != "" {
		query := " select * from products " + pr.pgUtilRepo.HandlingPaginationWhere(model.GetSeatchParamProduct(), q, "", "")
		// pr.GetSearchParam(q, limit, offset)
		if err := pr.db.Raw(query).Scan(&data).Error; err != nil {
			return data, err
		}

	} else {
		if err := pr.db.Order("name").Limit(limit).Offset(offset).Find(&data).Error; err != nil {
			return data, err
		}
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
	if err := pr.db.Updates(&data).Error; err != nil {
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
