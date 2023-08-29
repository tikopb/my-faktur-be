package rest

import (
	"bemyfaktur/internal/usecase/invoice"
	"bemyfaktur/internal/usecase/partner"
	"bemyfaktur/internal/usecase/product"
	"strconv"

	"github.com/labstack/echo/v4"
)

type handler struct {
	partnerUsecase partner.Usecase
	productUsecase product.ProductUsecaseInterface
	invoiceUsecase invoice.InvoiceUsecaseInterface
}

func NewHandler(partnerUsecase partner.Usecase, productUsecase product.ProductUsecaseInterface, invoiceUsecase invoice.InvoiceUsecaseInterface) *handler {

	return &handler{
		partnerUsecase: partnerUsecase,
		productUsecase: productUsecase,
		invoiceUsecase: invoiceUsecase,
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
	PartnerId := c.Param("id")
	Id, err := strconv.Atoi(PartnerId)
	if err != nil {
		panic(err)
	}
	return Id
}

func handlingLimitAndOffset(c echo.Context) (int, int) {
	// Get query parameters with default values
	limitStr := c.QueryParam("limit")
	offsetStr := c.QueryParam("offset")

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
