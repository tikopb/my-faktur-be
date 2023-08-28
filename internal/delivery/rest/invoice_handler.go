package rest

import (
	"bemyfaktur/internal/model"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func (h *handler) IndexInvoice(c echo.Context) error {
	limitStr := c.QueryParam("limit")
	if limitStr == "" {
		limitStr = "15" // Default value
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		return handleError(c, http.StatusInternalServerError, err)
	}

	offsetStr := c.QueryParam("offset")
	if offsetStr == "" {
		offsetStr = "0"
	}
	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		return handleError(c, http.StatusInternalServerError, err)
	}

	data, err := h.invoiceUsecase.IndexInvoice(limit, offset)
	if err != nil {
		return handleError(c, http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": data,
	})
}

func (h *handler) GetInvoice(c echo.Context) error {
	id := transformIdToInt(c)

	data, err := h.productUsecase.GetProduct(id)
	if err != nil {
		return handleError(c, http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": data,
	})
}

func (h *handler) CreateInvoice(c echo.Context) error {
	var request model.InvoiceRequest
	err := json.NewDecoder(c.Request().Body).Decode(&request)
	if err != nil {
		return handleError(c, http.StatusInternalServerError, err)
	}

	data, err := h.invoiceUsecase.CreateInvoice(request)
	if err != nil {
		return handleError(c, http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": data,
	})
}

func (h *handler) UpdateInvoice(c echo.Context) error {
	var request model.Invoice
	//get param
	id := transformIdToInt(c)

	//run function
	err := json.NewDecoder(c.Request().Body).Decode(&request)
	if err != nil {
		return handleError(c, http.StatusInternalServerError, err)
	}

	data, err := h.invoiceUsecase.UpdatedInvoice(id, request)
	if err != nil {
		return handleError(c, http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": data,
		"msg":  "data updated",
	})
}

func (h *handler) DeleteInvoice(c echo.Context) error {
	//get param
	id := transformIdToInt(c)

	data, err := h.invoiceUsecase.DeleteInvoice(id)
	if err != nil {
		return handleError(c, http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": data,
		"msg":  "data deleted",
	})
}
