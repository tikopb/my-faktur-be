package paginationutil

import (
	"fmt"
	"math"
	"strings"

	"gorm.io/gorm"
)

type Pagination struct {
	Current_page int   `json:"current_page"`
	Total_page   int   `json:"total_page"`
	Per_page     int   `json:"per_page"`
	Total_data   int64 `json:"total_data"`
}

type Repository interface {
	PaginationUtil(tabelName string, searchParam []string, limit int, offset int, q string, startDate string, endDate string) (Pagination, error)
	PaginationUtilWithJoinTable(count int64, limit int, offset int) (Pagination, error)
	HandlingPaginationWhere(searchParam []string, q string, startDate string, endDate string) string
}

type paginationUtilRepo struct {
	db *gorm.DB
}

func GetRepository(db *gorm.DB) *paginationUtilRepo {
	return &paginationUtilRepo{
		db: db,
	}
}

// master class pagination
func (p *paginationUtilRepo) PaginationUtil(tabelName string, searchParam []string, limit int, offset int, q string, startDate string, endDate string) (Pagination, error) {
	meta := Pagination{}
	param := ""
	q = strings.ToLower(q)

	//if searching where is not null then execute for where prosees
	if q != "" && len(searchParam) > 0 {
		param = p.HandlingPaginationWhere(searchParam, q, startDate, endDate)
	}

	//get count data total with where variabel
	query := ` select count(id) as count from ` + tabelName

	var count int64
	if q != "" {
		query = query + param
		if err := p.db.Raw(query).Scan(&count).Error; err != nil {
			return meta, err
		}
	} else {
		if err := p.db.Raw(query).Scan(&count).Error; err != nil {
			fmt.Println(query)
			return meta, err
		}
	}

	totalPage := math.Ceil(float64(count) / float64(limit))

	//set meta data
	meta.Current_page = offset + 1
	meta.Total_page = int(totalPage)
	meta.Per_page = limit
	meta.Total_data = count

	return meta, nil
}

func (p *paginationUtilRepo) HandlingPaginationWhere(searchParam []string, q string, startDate string, endDate string) string {

	//lower case the q parameter
	q = strings.ToLower(q)
	// execute looping param) string {
	// execute looping param
	var param string
	for i, searchparam := range searchParam {
		if i == len(searchParam)-1 {
			param += "lower(" + searchparam + ") like '%" + q + "%'"
		} else {
			param += "lower(" + searchparam + ") like '%" + q + "%' OR "
		}
	}

	// Add the condition for createdDate using BETWEEN
	if startDate != "" && endDate != "" {
		if param != "" {
			param += " AND "
		}
		param += "createdDate BETWEEN '" + startDate + "' AND '" + endDate + "'"
	}

	// Create the WHERE clause
	if param != "" {
		param = " WHERE " + param
	}

	return param
}

func (p *paginationUtilRepo) PaginationUtilWithJoinTable(count int64, limit int, offset int) (Pagination, error) {
	meta := Pagination{}

	totalPage := math.Ceil(float64(count) / float64(limit))
	//set meta data
	meta.Current_page = offset + 1
	meta.Total_page = int(totalPage)
	meta.Per_page = limit
	meta.Total_data = count

	return meta, nil
}
