package rest

import (
	"bemyfaktur/internal/model"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

func (h *handler) IndexPayment(c echo.Context) error {

	//set param
	limit, offset := HandlingLimitAndOffset(c)

	//get parameter
	q := c.QueryParam("q")
	q = strings.ToLower(q)

	data, err := h.paymentUsecase.Indexpayment(limit, offset, q)
	if err != nil {
		WriteLogErorr("[delivery][rest][payment_handler][IndexPayment] ", err)
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	//meta data field
	count, err := h.paymentUsecase.HandlingPagination(q, limit, offset)
	if err != nil {
		WriteLogErorr("[delivery][rest][payment_handler][IndexPayment] ", err)
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	meta, err = h.pgUtilRepo.PaginationUtilWithJoinTable(int64(count), limit, offset)
	if err != nil {
		WriteLogErorr("[delivery][rest][payment_handler][IndexPayment] ", err)
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	return handleError(c, http.StatusOK, errors.New("GET SUCCESS"), meta, data)
}

func (h *handler) Getpayment(c echo.Context) error {
	id := transformIdToInt(c)

	data, err := h.paymentUsecase.Getpayment(id)
	if err != nil {
		WriteLogErorr("[delivery][rest][payment_handler][Getpayment] ", err)
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	return handleError(c, http.StatusOK, errors.New("GET SUCCESS"), meta, data)
}

func (h *handler) CreatePayment(c echo.Context) error {
	var request model.PaymentRequest
	err := json.NewDecoder(c.Request().Body).Decode(&request)
	if err != nil {
		WriteLogErorr("[delivery][rest][payment_handler][CreatePayment] ", err)
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	//getUserId
	userId, err := h.middleware.GetuserId(c.Request())
	if err != nil {
		WriteLogErorr("[delivery][rest][payment_handler][CreatePayment] ", err)
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	data, err := h.paymentUsecase.Createpayment(request, userId)
	if err != nil {
		WriteLogErorr("[delivery][rest][payment_handler][CreatePayment] ", err)
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	return handleError(c, http.StatusOK, errors.New("CREATE "+data.BatchNo+" SUCCESS"), meta, data)

}

func (h *handler) UpdatePayment(c echo.Context) error {
	var request model.PaymentRequest
	id := transformIdToInt(c)
	err := json.NewDecoder(c.Request().Body).Decode(&request)
	if err != nil {
		WriteLogErorr("[delivery][rest][payment_handler][UpdatePayment] ", err)
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	data, err := h.paymentUsecase.Updatedpayment(id, request)
	if err != nil {
		WriteLogErorr("[delivery][rest][payment_handler][UpdatePayment] ", err)
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	return handleError(c, http.StatusOK, errors.New("UPDATE "+data.BatchNo+" SUCCESS"), meta, data)

}

func (h *handler) DeletePayment(c echo.Context) error {
	id := transformIdToInt(c)
	data, err := h.paymentUsecase.Deletepayment(id)
	if err != nil {
		WriteLogErorr("[delivery][rest][payment_handler][DeletePayment] ", err)
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	return handleError(c, http.StatusOK, errors.New("DELETE "+data+" SUCCESS"), meta, data)
}
