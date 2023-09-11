package documentutil

import (
	"bemyfaktur/internal/model"
	"errors"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
)

type Repository interface {
	GetDocumentNo(tableName string) (string, error)
}

type documentUtilRepo struct {
	db *gorm.DB
}

func GetRepository(db *gorm.DB) Repository {
	return &documentUtilRepo{
		db: db,
	}
}

// GetDocumentNo implements Repository.
func (dr *documentUtilRepo) GetDocumentNo(tableName string) (string, error) {
	count, err := dr.GetCount(tableName)
	if err != nil {
		return "", err
	}
	documentNoTemp, err := dr.GetDocumentNoTemp(tableName)
	if err != nil {
		return "", err
	}

	//INV-
	documentNo := documentNoTemp.Prefix + "-" + count + "-" + documentNoTemp.Suffix

	return documentNo, nil
}

func (dr *documentUtilRepo) GetDocumentNoTemp(tableName string) (model.DocumentNoTemp, error) {
	data := model.DocumentNoTemp{}
	now := time.Now()

	if err := dr.db.Where("table_name = ? AND start_date <= ? AND end_date >= ?", tableName, now, now).First(&data).Error; err != nil {
		return data, err
	}

	return data, nil
}

/*
	get counting of last number table

if not get the data created first and create agait
*/
func (dr *documentUtilRepo) GetCount(tableName string) (string, error) {
	var count string

	query := `
	SELECT counting 
	FROM document_no_temps dnt
	WHERE table_name = ?
	AND start_date <= NOW()
	AND end_date >= NOW();
	`
	if err := dr.db.Raw(query, tableName).Scan(&count).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err := dr.createDocumentNoTemp(tableName)
			if err != nil {
				return "", err
			}
			dr.GetCount(tableName) //looping state
		}
		return count, err
	}

	formattedCounting := strings.Repeat("0", 4-len(fmt.Sprint(count))) + fmt.Sprint(count)

	return formattedCounting, nil

}

/*
creating data documentno temp
- created data base on time.now.month first date and times.now.last date
- insert tableanme and prefix suffix from last data and created
*/
func (dr *documentUtilRepo) createDocumentNoTemp(tableName string) error {
	firstDate, lastDate := dr.getFirstAndLastDateOfMonth()
	data := model.DocumentNoTemp{
		Prefix:    "INV",
		Suffix:    "JS",
		TableName: tableName, //##@ fix this!
		StartDate: firstDate,
		EndDate:   lastDate,
	}

	if err := dr.db.Create(&data).Error; err != nil {
		return err
	}
	return nil
}

// get firtst and last date of current run month
func (dr *documentUtilRepo) getFirstAndLastDateOfMonth() (firstDate, lastDate time.Time) {
	now := time.Now()
	year, month, _ := now.Date()

	// Calculate the first day of the current month
	firstDate = time.Date(year, month, 1, 0, 0, 0, 0, now.Location())

	// Calculate the last day of the current month
	nextMonth := month + 1
	if nextMonth > 12 {
		nextMonth = 1
		year++
	}
	lastDate = time.Date(year, nextMonth, 1, 0, 0, 0, 0, now.Location()).Add(-time.Second)

	return firstDate, lastDate
}
