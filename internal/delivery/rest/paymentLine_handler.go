package rest

import (
	"bemyfaktur/internal/model"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func (h *handler) IndexPaymentLine(c echo.Context) error {
	//set param
	limit, offset := HandlingLimitAndOffset(c)
	paymentIdParam := c.QueryParam("paymentid")
	paymentId, err := strconv.Atoi(paymentIdParam)
	q := c.QueryParam("q")
	if err != nil {
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	data, err := h.paymentUsecase.IndexLine(limit, offset, paymentId, q)
	if err != nil {
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	//meta data field
	count, err := h.invoiceUsecase.HandlingPaginationLine(q, limit, offset, paymentId)
	if err != nil {
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	meta, err = h.pgUtilRepo.PaginationUtilWithJoinTable(int64(count), limit, offset)
	if err != nil {
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	return handleError(c, http.StatusOK, errors.New("GET SUCCESS"), meta, data)
}

func (h *handler) GetPaymentLine(c echo.Context) error {
	id := transformIdToInt(c)

	data, err := h.paymentUsecase.GetPaymentLine(id)
	if err != nil {
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	return handleError(c, http.StatusOK, errors.New("GET SUCCESS"), meta, data)
}

func (h *handler) CreatePaymentLine(c echo.Context) error {
	//set param
	var request model.PaymentLineRequest
	err := json.NewDecoder(c.Request().Body).Decode(&request)
	if err != nil {
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	//getUserId
	userId, err := h.middleware.GetuserId(c.Request())
	if err != nil {
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	//run function
	data, err := h.paymentUsecase.CreatePaymentLine(request, userId)
	if err != nil {
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	return handleError(c, http.StatusOK, errors.New("CREATE SUCCESS"), meta, data)
}

func (h *handler) UpdatePaymentLine(c echo.Context) error {
	//set param
	id := transformIdToInt(c)
	var request model.PaymentLineRequest
	err := json.NewDecoder(c.Request().Body).Decode(&request)
	if err != nil {
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	//run function
	data, err := h.paymentUsecase.UpdatedPaymentLine(id, request)
	if err != nil {
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	return handleError(c, http.StatusOK, errors.New("UPDATE SUCCESS"), meta, data)
}

func (h *handler) DeletePaymentLine(c echo.Context) error {
	//set param
	id := transformIdToInt(c)

	data, err := h.paymentUsecase.DeletePaymentLine(id)
	if err != nil {
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	return handleError(c, http.StatusOK, errors.New("DELETE SUCCESS"), meta, data)
}
