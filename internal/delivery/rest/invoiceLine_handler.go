package rest

import (
	"bemyfaktur/internal/model"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *handler) IndexInvoiceLine(c echo.Context) error {
	//set limit and offset
	limit, offset := HandlingLimitAndOffset(c)
	q := c.QueryParam("q")

	//get invoice ID
	invoiceId, err := h.parsingId(c)
	if err != nil {
		WriteLogErorr("[delivery][rest][invoiceLine_handler][GetInvoiceLine] ", err)

		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	if err != nil {
		WriteLogErorr("[delivery][rest][invoiceLine_handler][IndexInvoiceLine] ", err)
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	//setOrderData
	order, err := h.GetOrderClauses(c)
	if err != nil {
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	data, err := h.invoiceUsecase.IndexLine(limit, offset, invoiceId, q, order)
	if err != nil {
		WriteLogErorr("[delivery][rest][invoiceLine_handler][IndexInvoiceLine] ", err)
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	//meta data field
	if len(data) > 0 {
		invoiceIntId := data[0].Invoice.Id
		count, err := h.invoiceUsecase.HandlingPaginationLine(q, limit, offset, invoiceIntId)
		if err != nil {
			WriteLogErorr("[delivery][rest][invoiceLine_handler][IndexInvoiceLine] ", err)
			return handleError(c, http.StatusInternalServerError, err, meta, data)
		}
		meta, err = h.pgUtilRepo.PaginationUtilWithJoinTable(int64(count), limit, offset)
		if err != nil {
			WriteLogErorr("[delivery][rest][invoiceLine_handler][IndexInvoiceLine] ", err)
			return handleError(c, http.StatusInternalServerError, err, meta, data)
		}
	} else {
		meta = nil
	}

	return handleError(c, http.StatusOK, errors.New("GET SUCCESS"), meta, data)

}

func (h *handler) GetInvoiceLine(c echo.Context) error {
	id, err := h.parsingId(c)
	if err != nil {
		WriteLogErorr("[delivery][rest][invoiceLine_handler][GetInvoiceLine] ", err)

		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	data, err := h.invoiceUsecase.GetInvoiceLine(id)
	if err != nil {
		WriteLogErorr("[delivery][rest][invoiceLine_handler][GetInvoiceLine] ", err)
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	return handleError(c, http.StatusOK, errors.New("GET SUCCESS"), meta, data)
}

func (h *handler) CreateInvoiceLine(c echo.Context) error {
	var request model.InvoiceLineRequest
	err := json.NewDecoder(c.Request().Body).Decode(&request)
	if err != nil {
		WriteLogErorr("[delivery][rest][invoiceLine_handler][CreateInvoiceLine] ", err)
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	//getUserId
	userId, err := h.middleware.GetuserId(c.Request())
	if err != nil {
		WriteLogErorr("[delivery][rest][invoiceLine_handler][CreateInvoiceLine] ", err)
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	data, err := h.invoiceUsecase.CreateInvoiceLine(request, userId)
	if err != nil {
		WriteLogErorr("[delivery][rest][invoiceLine_handler][CreateInvoiceLine] ", err)
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	return handleError(c, http.StatusOK, errors.New("CREATE SUCCESS"), meta, data)
}

func (h *handler) UpdatedInvoiceLine(c echo.Context) error {
	//get param
	id, err := h.parsingId(c)
	if err != nil {
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	var request model.InvoiceLineRequest
	err = json.NewDecoder(c.Request().Body).Decode(&request)
	if err != nil {
		WriteLogErorr("[delivery][rest][invoiceLine_handler][UpdatedInvoiceLine] ", err)
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	//getUserId
	userId, err := h.middleware.GetuserId(c.Request())
	if err != nil {
		WriteLogErorr("[delivery][rest][invoiceLine_handler][CreateInvoiceLine] ", err)
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}
	request.UpdatedById = userId

	//run function
	data, err := h.invoiceUsecase.UpdatedInvoiceLine(id, request)
	if err != nil {
		WriteLogErorr("[delivery][rest][invoiceLine_handler][UpdatedInvoiceLine] ", err)
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	return handleError(c, http.StatusOK, errors.New("UPDATE SUCCESS"), meta, data)
}

func (h *handler) DeleteInvoiceLine(c echo.Context) error {
	//get param
	id, err := h.parsingId(c)
	if err != nil {
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	data, err := h.invoiceUsecase.DeleteInvoiceLine(id)
	if err != nil {
		WriteLogErorr("[delivery][rest][invoiceLine_handler][DeleteInvoiceLine] ", err)
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	return handleError(c, http.StatusOK, errors.New("DELETE SUCCESS"), meta, data)
}
