package rest

import (
	"bemyfaktur/internal/usecase/invoice"
	"bemyfaktur/internal/usecase/partner"
	"bemyfaktur/internal/usecase/payment"
	"bemyfaktur/internal/usecase/product"
	"math"
	"strconv"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type handler struct {
	partnerUsecase partner.Usecase
	productUsecase product.ProductUsecaseInterface
	invoiceUsecase invoice.InvoiceUsecaseInterface
	paymentUsecase payment.PaymentUsecaseInterface
}

type pagination struct {
	Page       int
	Limit      int
	Total_page int
	Offset     int
}

func NewHandler(partnerUsecase partner.Usecase, productUsecase product.ProductUsecaseInterface, invoiceUsecase invoice.InvoiceUsecaseInterface, paymentUsecase payment.PaymentUsecaseInterface) *handler {

	return &handler{
		partnerUsecase: partnerUsecase,
		productUsecase: productUsecase,
		invoiceUsecase: invoiceUsecase,
		paymentUsecase: paymentUsecase,
	}
}

func handleError(c echo.Context, statusCode int, err error) error {
	return c.JSON(statusCode, map[string]interface{}{
		"error": err.Error(),
		"msg":   "internal error" + err.Error(),
	})
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

func PaginationUtil(tabelName string, searchParam []string, db *gorm.DB, limit int, offset int) pagination {
	meta := pagination{}

	//#run query for pagination
	//execute looping param
	var param string
	if len(searchParam) > 0 {
		for _, searchparam := range searchParam {
			param = searchparam + ","
		}
	}

	//get count data total with where variabel
	query := `
		select count(id) as countData from invoices i 
	`
	if len(searchParam) > 0 {
		query += `where = ?,`
	}
	var count int
	db.Raw(query, tabelName, param).Scan(&count)

	totalPage := math.Ceil(float64(count) / float64(limit))

	//set meta data
	meta.Limit = limit
	meta.Offset = offset
	meta.Total_page = int(totalPage)
	meta.Page = offset

	return meta
}
