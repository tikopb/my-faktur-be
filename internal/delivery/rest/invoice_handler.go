package rest

import (
	"bemyfaktur/internal/model"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

func (h *handler) IndexInvoice(c echo.Context) error {
	//set param
	limit, offset := HandlingLimitAndOffset(c)

	//get parameter
	q := c.QueryParam("q")
	q = strings.ToLower(q)

	//setOrderData
	order, err := h.GetOrderClauses(c)
	if err != nil {
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	data, err := h.invoiceUsecase.IndexInvoice(limit, offset, q, order)
	if err != nil {
		WriteLogErorr("[delivery][rest][invoice_handler][IndexInvoice] ", err)
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	//meta data field
	count, err := h.invoiceUsecase.HandlingPagination(q, limit, offset)
	if err != nil {
		WriteLogErorr("[delivery][rest][invoice_handler][IndexInvoice] ", err)
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	meta, err = h.pgUtilRepo.PaginationUtilWithJoinTable(int64(count), limit, offset)
	if err != nil {
		WriteLogErorr("[delivery][rest][invoice_handler][IndexInvoice] ", err)
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	return handleError(c, http.StatusOK, errors.New("GET SUCCESS"), meta, data)
}

func (h *handler) GetInvoice(c echo.Context) error {
	id, err := h.parsingId(c)
	if err != nil {
		WriteLogErorr("[delivery][rest][invoice_handler][GetInvoice] ", err)
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	data, err := h.invoiceUsecase.GetInvoice(id)
	if err != nil {
		WriteLogErorr("[delivery][rest][invoice_handler][GetInvoice] ", err)
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	return handleError(c, http.StatusOK, errors.New("GET SUCCESS"), meta, data)
}

func (h *handler) CreateInvoice(c echo.Context) error {
	var request model.InvoiceRequest
	err := json.NewDecoder(c.Request().Body).Decode(&request)
	if err != nil {
		WriteLogErorr("[delivery][rest][invoice_handler][CreateInvoice] ", err)
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	//getUserId
	userId, err := h.middleware.GetuserId(c.Request())
	if err != nil {
		WriteLogErorr("[delivery][rest][invoice_handler][CreateInvoice] ", err)
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	data, err := h.invoiceUsecase.CreateInvoice(request, userId)
	if err != nil {
		WriteLogErorr("[delivery][rest][invoice_handler][CreateInvoice] ", err)
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}
	return handleError(c, http.StatusOK, errors.New("CREATE "+data.BatchNo+" SUCCESS"), meta, data)
}

func (h *handler) UpdateInvoice(c echo.Context) error {
	var request model.InvoiceRequest
	//get param
	id, err := h.parsingId(c)
	if err != nil {
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	//run function
	err = json.NewDecoder(c.Request().Body).Decode(&request)
	if err != nil {
		WriteLogErorr("[delivery][rest][invoice_handler][UpdateInvoice] ", err)
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	//getUpdateByUserId
	userId, err := h.middleware.GetuserId(c.Request())
	if err != nil {
		WriteLogErorr("[delivery][rest][invoice_handler][CreateInvoice] ", err)
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	data, err := h.invoiceUsecase.UpdatedInvoice(id, request, userId)
	if err != nil {
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	return handleError(c, http.StatusOK, errors.New("UPDATE "+data.BatchNo+" SUCCESS"), meta, data)

}

func (h *handler) DeleteInvoice(c echo.Context) error {
	//get param
	id, err := h.parsingId(c)
	if err != nil {
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	data, err := h.invoiceUsecase.DeleteInvoice(id)
	if err != nil {
		WriteLogErorr("[delivery][rest][invoice_handler][DeleteInvoice] ", err)
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	return handleError(c, http.StatusOK, errors.New("DELETE "+data+"SUCCESS"), meta, data)

}
