package rest

import (
	"bemyfaktur/internal/usecase/auth"
	"bemyfaktur/internal/usecase/fileservice"
	"bemyfaktur/internal/usecase/invoice"
	"bemyfaktur/internal/usecase/organization"
	"bemyfaktur/internal/usecase/partner"
	"bemyfaktur/internal/usecase/payment"
	"bemyfaktur/internal/usecase/product"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"strconv"

	midUtil "bemyfaktur/internal/delivery/auth"

	pgUtil "bemyfaktur/internal/model/paginationUtil"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// Declare meta and data at the package level
var meta interface{}
var data interface{}

type handler struct {
	partnerUsecase      partner.Usecase
	productUsecase      product.ProductUsecaseInterface
	invoiceUsecase      invoice.InvoiceUsecaseInterface
	paymentUsecase      payment.PaymentUsecaseInterface
	fileserviceUsecase  fileservice.Usecase
	authUsecase         auth.Usecase
	db                  *gorm.DB
	pgUtilRepo          pgUtil.Repository
	middleware          midUtil.MidlewareInterface
	organizationUsecase organization.Usecase
}

type handlerRespont struct {
	Status  int
	Message string
	Meta    interface{} `json:"meta"`
	Data    interface{} `json:"data"`
}

func NewHandler(partnerUsecase partner.Usecase, productUsecase product.ProductUsecaseInterface, invoiceUsecase invoice.InvoiceUsecaseInterface, paymentUsecase payment.PaymentUsecaseInterface, fileserviceUsecase fileservice.Usecase, pgRepo pgUtil.Repository, authUsecase auth.Usecase, organizationUsecase organization.Usecase, middleware midUtil.MidlewareInterface, db *gorm.DB) *handler {

	return &handler{
		partnerUsecase:      partnerUsecase,
		productUsecase:      productUsecase,
		invoiceUsecase:      invoiceUsecase,
		paymentUsecase:      paymentUsecase,
		fileserviceUsecase:  fileserviceUsecase,
		authUsecase:         authUsecase,
		pgUtilRepo:          pgRepo,
		db:                  db,
		middleware:          middleware,
		organizationUsecase: organizationUsecase,
	}
}

func handleError(c echo.Context, statusCode int, err error, meta interface{}, data interface{}) error {
	var response handlerRespont

	if strings.Contains(err.Error(), "data not found") {
		statusCode = http.StatusNotFound
	}

	if statusCode != http.StatusOK && statusCode != http.StatusCreated {
		response = handlerRespont{
			Status:  statusCode,
			Message: "internal error: " + err.Error(),
			Meta:    nil,
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
	if limitStr == "" { //convertion default value if offset not found
		limitStr = "15" // Default value
	}
	offsetStr := c.QueryParam("offset")
	if offsetStr == "" { //convertion default value if offset not found
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

func (h *handler) HandlingDateFromAndDateTo(c echo.Context) (string, string, error) {
	dateFromParam := c.QueryParam("date_from")
	dateToParam := c.QueryParam("date_to")
	if dateFromParam == "" || dateToParam == "" {
		return "", "", errors.New("date from and date To can't be null")
	}
	layout := "2006-01-02"

	// parsing date From
	dateFrom, err := time.Parse(layout, dateFromParam)
	if err != nil {
		return "", "", err
	}
	dateTo, err := time.Parse(layout, dateToParam)
	if err != nil {
		return "", "", err
	}

	// Formatting dateFrom and dateTo as "YYYY-MM-DD" strings
	formattedDateFrom := dateFrom.Format("2006-01-02")
	formattedDateTo := dateTo.Format("2006-01-02")

	return formattedDateFrom, formattedDateTo, nil
}

func WriteLogErorr(msg string, err error) {
	formatter := &logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	}
	logrus.SetFormatter(formatter)
	logrus.WithFields(logrus.Fields{
		"err": err,
	}).Error(msg, err.Error())
}

func WriteLogInfo(msg string) {
	formatter := &logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	}
	logrus.SetFormatter(formatter)
	logrus.Info(msg)
}

func (h *handler) ParsingUUID(value string) (uuid.UUID, error) {
	// Parse the string into a UUID
	uuid, err := uuid.Parse(value)
	if err != nil {
		// Handle the error, e.g., return an error response
		return uuid, err
	}

	return uuid, nil
}

func (h *handler) GetOrderClauses(c echo.Context) ([]string, error) {
	// Get the URL parameters
	sort := c.QueryParam("sort")
	order := c.QueryParam("order")
	orderClauses := []string{}

	if strings.TrimSpace(sort) == "" && strings.TrimSpace(order) == "" {
		return orderClauses, nil
	}
	// Create the order clauses
	for _, field := range strings.Split(sort, ",") {
		orderClauses = append(orderClauses, fmt.Sprintf("%s %s", field, order))
	}

	return orderClauses, nil
}
