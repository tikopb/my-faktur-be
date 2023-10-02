package rest

import (
	"bemyfaktur/internal/usecase/invoice"
	"bemyfaktur/internal/usecase/partner"
	"bemyfaktur/internal/usecase/payment"
	"bemyfaktur/internal/usecase/product"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// Declare meta and data at the package level
var meta interface{}
var data interface{}

type handler struct {
	partnerUsecase partner.Usecase
	productUsecase product.ProductUsecaseInterface
	invoiceUsecase invoice.InvoiceUsecaseInterface
	paymentUsecase payment.PaymentUsecaseInterface
	db             *gorm.DB
}

type pagination struct {
	Current_page int   `json:"current_page"`
	Total_page   int   `json:"total_page"`
	Per_page     int   `json:"per_page"`
	Total_data   int64 `json:"total_data"`
}

type handlerRespont struct {
	Status  int
	Message string
	Meta    interface{} `json:"meta"`
	Data    interface{} `json:"data"`
}

func NewHandler(partnerUsecase partner.Usecase, productUsecase product.ProductUsecaseInterface, invoiceUsecase invoice.InvoiceUsecaseInterface, paymentUsecase payment.PaymentUsecaseInterface, db *gorm.DB) *handler {

	return &handler{
		partnerUsecase: partnerUsecase,
		productUsecase: productUsecase,
		invoiceUsecase: invoiceUsecase,
		paymentUsecase: paymentUsecase,
		db:             db,
	}
}

func handleError(c echo.Context, statusCode int, err error, meta interface{}, data interface{}) error {
	var response handlerRespont
	if statusCode != 200 {
		response = handlerRespont{
			Status:  statusCode,
			Message: "internal error: " + err.Error(),
			Meta:    meta,
			Data:    data,
		}
	} else {
		response = handlerRespont{
			Status:  statusCode,
			Message: "PROCESS SUCCESS: " + err.Error(),
			Meta:    meta,
			Data:    data,
		}
	}

	return c.JSON(statusCode, response)
}

func transformIdToInt(c echo.Context) int {
	// get param
	ID := c.Param("id")
	Id, err := strconv.Atoi(ID)
	if err != nil {
		panic(err)
	}
	return Id
}

func HandlingLimitAndOffset(c echo.Context) (int, int) {
	// Get query parameters with default values
	limitStr := c.QueryParam("limit")
	if limitStr == "" {
		limitStr = "15" // Default value
	}
	offsetStr := c.QueryParam("offset")
	if offsetStr == "" {
		offsetStr = "0"
	}

	// Convert to integers with error handling
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		panic("error converting 'limit' to integer")
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		panic("error converting 'offset' to integer")
	}

	// Return the values
	return limit, offset
}

// ===== PAGINATION TOOLS =====
// master class pagination
func (h *handler) PaginationUtil(tabelName string, searchParam []string, limit int, offset int, q string) (pagination, error) {
	meta := pagination{}
	param := ""
	q = strings.ToLower(q)

	//if searching where is not null then execute for where prosees
	if q != "" && len(searchParam) > 0 {
		param = h.HandlingPaginationWhere(searchParam, q)
	}

	//get count data total with where variabel
	query := ` select count(id) as count from ` + tabelName

	var count int64
	if q != "" {
		query = query + param
		if err := h.db.Raw(query).Scan(&count).Error; err != nil {
			return meta, err
		}
	} else {
		if err := h.db.Raw(query).Scan(&count).Error; err != nil {
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

func (h *handler) HandlingPaginationWhere(searchParam []string, q string) string {
	//execute looping param
	var param string
	for i, searchparam := range searchParam {
		if i == len(searchParam)-1 {
			param += "lower(" + searchparam + ") like '%" + q + "%'"
		} else {
			param += "lower(" + searchparam + ") like '%" + q + "%' OR "
		}
	}

	param = " where " + param
	return param
}

func (p *handler) PaginationUtilWithJoinTable(count int64, limit int, offset int) (pagination, error) {
	meta := pagination{}

	totalPage := math.Ceil(float64(count) / float64(limit))
	//set meta data
	meta.Current_page = offset + 1
	meta.Total_page = int(totalPage)
	meta.Per_page = limit
	meta.Total_data = count

	return meta, nil
}
