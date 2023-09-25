package rest

import (
	"bemyfaktur/internal/model"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *handler) IndexPayment(c echo.Context) error {

	//set param
	limit, offset := HandlingLimitAndOffset(c)

	data, err := h.paymentUsecase.Indexpayment(limit, offset)
	if err != nil {
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	return handleError(c, http.StatusOK, errors.New("GET SUCCESS"), meta, data)
}

func (h *handler) Getpayment(c echo.Context) error {
	id := transformIdToInt(c)

	data, err := h.paymentUsecase.Getpayment(id)
	if err != nil {
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	return handleError(c, http.StatusOK, errors.New("GET SUCCESS"), meta, data)
}

func (h *handler) CreatePayment(c echo.Context) error {
	var request model.PaymentRequest
	err := json.NewDecoder(c.Request().Body).Decode(&request)
	if err != nil {
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	data, err := h.paymentUsecase.Createpayment(request)
	if err != nil {
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	return handleError(c, http.StatusOK, errors.New("CREATE "+data.BatchNo+" SUCCESS"), meta, data)

}

func (h *handler) UpdatePayment(c echo.Context) error {
	var request model.PaymentRequest
	id := transformIdToInt(c)
	err := json.NewDecoder(c.Request().Body).Decode(&request)
	if err != nil {
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	data, err := h.paymentUsecase.Updatedpayment(id, request)
	if err != nil {
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	return handleError(c, http.StatusOK, errors.New("UPDATE "+data.BatchNo+" SUCCESS"), meta, data)

}

func (h *handler) DeletePayment(c echo.Context) error {
	id := transformIdToInt(c)
	data, err := h.paymentUsecase.Deletepayment(id)
	if err != nil {
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	return handleError(c, http.StatusOK, errors.New("DELETE "+data+" SUCCESS"), meta, data)
}
