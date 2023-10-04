package rest

import (
	"bemyfaktur/internal/usecase/invoice"
	"bemyfaktur/internal/usecase/partner"
	"bemyfaktur/internal/usecase/payment"
	"bemyfaktur/internal/usecase/product"

	"strconv"

	pgUtil "bemyfaktur/internal/model/paginationUtil"

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
	pgUtilRepo     pgUtil.Repository
}

type handlerRespont struct {
	Status  int
	Message string
	Meta    interface{} `json:"meta"`
	Data    interface{} `json:"data"`
}

func NewHandler(partnerUsecase partner.Usecase, productUsecase product.ProductUsecaseInterface, invoiceUsecase invoice.InvoiceUsecaseInterface, paymentUsecase payment.PaymentUsecaseInterface, pgRepo pgUtil.Repository, db *gorm.DB) *handler {

	return &handler{
		partnerUsecase: partnerUsecase,
		productUsecase: productUsecase,
		invoiceUsecase: invoiceUsecase,
		paymentUsecase: paymentUsecase,
		pgUtilRepo:     pgRepo,
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
