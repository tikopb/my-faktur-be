package rest

import (
	"bemyfaktur/internal/model"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *handler) IndexPaymentLine(c echo.Context) error {
	//set param
	limit, offset := HandlingLimitAndOffset(c)

	//id payment get
	paymentUUID, err := h.parsingId(c)
	if err != nil {
		WriteLogErorr("[delivery][rest][paymentline_handler][GetPaymentLine] ", err)
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	q := c.QueryParam("q")

	data, paymentId, err := h.paymentUsecase.IndexLine(limit, offset, paymentUUID, q)
	if err != nil {
		WriteLogErorr("[delivery][rest][paymentline_handler][IndexPaymentLine] ", err)

		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	//meta data field
	count, err := h.invoiceUsecase.HandlingPaginationLine(q, limit, offset, paymentId)
	if err != nil {
		WriteLogErorr("[delivery][rest][paymentline_handler][IndexPaymentLine] ", err)
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	meta, err = h.pgUtilRepo.PaginationUtilWithJoinTable(int64(count), limit, offset)
	if err != nil {
		WriteLogErorr("[delivery][rest][paymentline_handler][IndexPaymentLine] ", err)
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	return handleError(c, http.StatusOK, errors.New("GET SUCCESS"), meta, data)
}

func (h *handler) GetPaymentLine(c echo.Context) error {
	id, err := h.parsingId(c)
	if err != nil {
		WriteLogErorr("[delivery][rest][paymentline_handler][GetPaymentLine] ", err)
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	data, err := h.paymentUsecase.GetPaymentLine(id)
	if err != nil {
		WriteLogErorr("[delivery][rest][paymentline_handler][GetPaymentLine] ", err)
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	return handleError(c, http.StatusOK, errors.New("GET SUCCESS"), meta, data)
}

func (h *handler) CreatePaymentLine(c echo.Context) error {
	//set param
	var request model.PaymentLineRequest
	err := json.NewDecoder(c.Request().Body).Decode(&request)
	if err != nil {
		WriteLogErorr("[delivery][rest][paymentline_handler][CreatePaymentLine] ", err)
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	//getUserId
	userId, err := h.middleware.GetuserId(c.Request())
	if err != nil {
		WriteLogErorr("[delivery][rest][paymentline_handler][CreatePaymentLine] ", err)
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	//run function
	data, err := h.paymentUsecase.CreatePaymentLine(request, userId)
	if err != nil {
		WriteLogErorr("[delivery][rest][paymentline_handler][CreatePaymentLine] ", err)
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	return handleError(c, http.StatusOK, errors.New("CREATE SUCCESS"), meta, data)
}

func (h *handler) UpdatePaymentLine(c echo.Context) error {
	//set param
	id, err := h.parsingId(c)
	if err != nil {
		WriteLogErorr("[delivery][rest][paymentline_handler][UpdatePaymentLine] ", err)
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	var request model.PaymentLineRequest
	err = json.NewDecoder(c.Request().Body).Decode(&request)
	if err != nil {
		WriteLogErorr("[delivery][rest][paymentline_handler][UpdatePaymentLine] ", err)
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	//run function
	data, err := h.paymentUsecase.UpdatedPaymentLine(id, request)
	if err != nil {
		WriteLogErorr("[delivery][rest][paymentline_handler][UpdatePaymentLine] ", err)
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	return handleError(c, http.StatusOK, errors.New("UPDATE SUCCESS"), meta, data)
}

func (h *handler) DeletePaymentLine(c echo.Context) error {
	//set param
	id, err := h.parsingId(c)
	if err != nil {
		WriteLogErorr("[delivery][rest][paymentline_handler][DeletePaymentLine] ", err)
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	data, err := h.paymentUsecase.DeletePaymentLine(id)
	if err != nil {
		WriteLogErorr("[delivery][rest][paymentline_handler][DeletePaymentLine] ", err)
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	return handleError(c, http.StatusOK, errors.New("DELETE SUCCESS"), meta, data)
}
