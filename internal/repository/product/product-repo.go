package product

import (
	"bemyfaktur/internal/model"
	pgUtil "bemyfaktur/internal/model/paginationUtil"
	"errors"

	"github.com/google/uuid"
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

	//parsing the data to product respont
	data = pr.parsingProductToProductRespon(product)
	return data, nil
}

// Index implements Repository.
func (pr *productRepo) Index(limit int, offset int, q string) ([]model.ProductRespon, error) {
	data := []model.Product{}
	var dataReturn []model.ProductRespon

	if q != "" {
		query := " select * from products " + pr.pgUtilRepo.HandlingPaginationWhere(model.GetSeatchParamProduct(), q, "", "")
		if err := pr.db.Preload("User").Raw(query).Scan(&data).Error; err != nil {
			return dataReturn, err
		}

	} else {
		if err := pr.db.Preload("User").Order("name").Limit(limit).Offset(offset).Find(&data).Error; err != nil {
			return dataReturn, err
		}
	}

	//parsing to responFormat
	for _, product := range data {
		dataReturn = append(dataReturn, pr.parsingProductToProductRespon(product))
	}

	return dataReturn, nil
}

func (pr *productRepo) Partial(q string) ([]model.ProductPartialRespon, error) {
	var dataReturn []model.ProductPartialRespon

	if q != "" {
		stringParam := model.GetSeatchParamProductV2(q) + " AND isactive = true "
		if err := pr.db.Model(&model.Product{}).Select("CONCAT(value, ' - ', name) as name, uuid").Where(stringParam).Find(&dataReturn).Error; err != nil {
			return dataReturn, err
		}
	} else {
		if err := pr.db.Model(&model.Product{}).Select("CONCAT(value, ' - ', name) as name, uuid").Where(&model.Product{IsActive: true}).Find(&dataReturn).Error; err != nil {
			return dataReturn, err
		}
	}

	return dataReturn, nil

}

// Show implements Repository.
func (pr *productRepo) Show(id uuid.UUID) (model.ProductRespon, error) {
	var data model.Product
	dataReturn := model.ProductRespon{}

	if err := pr.db.Preload("User").Where("uuid", id).First(&data).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dataReturn, errors.New("data not found")
		}
	}

	dataReturn = pr.parsingProductToProductRespon(data)

	return dataReturn, nil
}

// Show implements Repository.
func (pr *productRepo) ShowInternal(id uuid.UUID) (model.Product, error) {
	var data model.Product

	if err := pr.db.Preload("User").Where("uuid", id).First(&data).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return data, errors.New("data not found")
		}
	}

	return data, nil
}

// Update implements Repository.
func (pr *productRepo) Update(id uuid.UUID, updatedProduct model.Product) (model.ProductRespon, error) {
	dataUpdated := model.ProductRespon{}
	data, err := pr.ShowInternal(id)

	if err != nil {
		return dataUpdated, err
	}

	//slicing data update
	data.Name = updatedProduct.Name
	data.Description = updatedProduct.Description
	data.IsActive = updatedProduct.User.IsActive
	data.Value = updatedProduct.Value
	data.Upc = updatedProduct.Upc

	//inisiate data udpated system
	if err := pr.db.Updates(&data).Error; err != nil {
		return dataUpdated, err
	}

	//inisiate data update version
	dataUpdated = pr.parsingProductToProductRespon(data)

	return dataUpdated, nil
}

// Delete implements Repository.
func (pr *productRepo) Delete(id uuid.UUID) (string, error) {
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

func (pr *productRepo) parsingProductToProductRespon(product model.Product) model.ProductRespon {
	data := model.ProductRespon{
		UUID:        product.UUID,
		Name:        product.Name,
		Description: product.Description,
		IsActive:    product.IsActive,
		CreatedAt:   product.CreatedAt,
		UpdateAt:    product.UpdateAt,
		CreatedBy:   product.User.Username,
		Value:       product.Value,
		Upc:         product.Upc,
	}

	return data
}
