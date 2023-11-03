package rest

import (
	"bemyfaktur/internal/usecase/auth"
	"bemyfaktur/internal/usecase/invoice"
	"bemyfaktur/internal/usecase/partner"
	"bemyfaktur/internal/usecase/payment"
	"bemyfaktur/internal/usecase/product"
	"fmt"
	"net/http"
	"strings"

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
	partnerUsecase partner.Usecase
	productUsecase product.ProductUsecaseInterface
	invoiceUsecase invoice.InvoiceUsecaseInterface
	paymentUsecase payment.PaymentUsecaseInterface
	authUsecase    auth.Usecase
	db             *gorm.DB
	pgUtilRepo     pgUtil.Repository
	middleware     midUtil.MidlewareInterface
}

type handlerRespont struct {
	Status  int
	Message string
	Meta    interface{} `json:"meta"`
	Data    interface{} `json:"data"`
}

func NewHandler(partnerUsecase partner.Usecase, productUsecase product.ProductUsecaseInterface, invoiceUsecase invoice.InvoiceUsecaseInterface, paymentUsecase payment.PaymentUsecaseInterface, pgRepo pgUtil.Repository, authUsecase auth.Usecase, middleware midUtil.MidlewareInterface, db *gorm.DB) *handler {

	return &handler{
		partnerUsecase: partnerUsecase,
		productUsecase: productUsecase,
		invoiceUsecase: invoiceUsecase,
		paymentUsecase: paymentUsecase,
		authUsecase:    authUsecase,
		pgUtilRepo:     pgRepo,
		db:             db,
		middleware:     middleware,
	}
}

func handleError(c echo.Context, statusCode int, err error, meta interface{}, data interface{}) error {
	var response handlerRespont

	if statusCode != http.StatusOK {
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

// func (h *handler) GetOrderClauses(c echo.Context) ([]string, error) {
// 	urlString := c.Request().URL.String()
// 	u, err := url.Parse(urlString)
// 	if err != nil {
// 		return nil, err
// 	}

// 	query := u.Query()
// 	var orderClauses []string

// 	// Declare the orderKey variable outside of the loop.
// 	var orderKey string

// 	for key, values := range query {
// 		if strings.HasPrefix(key, "sort") {
// 			for _, value := range values {
// 				orderKey = strings.TrimPrefix(value, "sort=")
// 				orderValue := "asc" // Default to ascending order
// 				orderKey = strings.TrimPrefix(orderKey, "order=")
// 				if orderValueParam := query.Get("order" + orderKey); orderValueParam != "" {
// 					orderValue = orderValueParam
// 				}
// 				orderClauses = append(orderClauses, fmt.Sprintf("%s %s", orderKey, orderValue))
// 			}

// 		}
// 	}

// 	fmt.Println(orderClauses)

// 	return orderClauses, nil
// }

func (h *handler) GetOrderClauses(c echo.Context) ([]string, error) {
	// Get the URL parameters
	sort := c.QueryParam("sort")
	order := c.QueryParam("order")

	if sort == "" && order == "" {
		return []string{}, nil
	}

	// Create the order clauses
	orderClauses := []string{}
	for _, field := range strings.Split(sort, ",") {
		orderClauses = append(orderClauses, fmt.Sprintf("%s %s", field, order))
	}

	fmt.Println(orderClauses)
	return orderClauses, nil
}
