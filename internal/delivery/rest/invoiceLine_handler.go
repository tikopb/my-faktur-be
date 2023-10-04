package rest

import (
	"bemyfaktur/internal/model"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func (h *handler) IndexInvoiceLine(c echo.Context) error {
	//set limit and offset
	limit, offset := HandlingLimitAndOffset(c)
	invoiceIdParam := c.QueryParam("invoiceId")
	invoiceId, err := strconv.Atoi(invoiceIdParam)
	q := c.QueryParam("q")

	if err != nil {
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	data, err := h.invoiceUsecase.IndexLine(limit, offset, invoiceId, q)
	if err != nil {
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	//meta data field
	count, err := h.invoiceUsecase.HandlingPaginationLine(q, limit, offset, invoiceId)
	if err != nil {
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	meta, err = h.pgUtilRepo.PaginationUtilWithJoinTable(int64(count), limit, offset)
	if err != nil {
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	return handleError(c, http.StatusOK, errors.New("GET SUCCESS"), meta, data)

}

func (h *handler) GetInvoiceLine(c echo.Context) error {
	id := transformIdToInt(c)

	data, err := h.invoiceUsecase.GetInvoiceLine(id)
	if err != nil {
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	return handleError(c, http.StatusOK, errors.New("GET SUCCESS"), meta, data)
}

func (h *handler) CreateInvoiceLine(c echo.Context) error {
	var request model.InvoiceLine
	err := json.NewDecoder(c.Request().Body).Decode(&request)
	if err != nil {
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	data, err := h.invoiceUsecase.CreateInvoiceLine(request)
	if err != nil {
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	return handleError(c, http.StatusOK, errors.New("CREATE SUCCESS"), meta, data)
}

func (h *handler) UpdatedInvoiceLine(c echo.Context) error {
	//get param
	id := transformIdToInt(c)
	var request model.InvoiceLine
	err := json.NewDecoder(c.Request().Body).Decode(&request)
	if err != nil {
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	//run function

	data, err := h.invoiceUsecase.UpdatedInvoiceLine(id, request, request.ProductID)
	if err != nil {
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	return handleError(c, http.StatusOK, errors.New("UPDATE SUCCESS"), meta, data)
}

func (h *handler) DeleteInvoiceLine(c echo.Context) error {
	//get param
	id := transformIdToInt(c)

	data, err := h.invoiceUsecase.DeleteInvoiceLine(id)
	if err != nil {
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	return handleError(c, http.StatusOK, errors.New("DELETE SUCCESS"), meta, data)
}
