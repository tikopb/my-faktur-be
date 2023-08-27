package rest

import (
	"bemyfaktur/internal/model"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func (h *handler) IndexInvoiceLine(c echo.Context) error {
	//set limit and offset
	limit, offset := handlingLimitAndOffset(c)
	invoiceIdParam := c.QueryParam("invoiceId")
	invoiceId, err := strconv.Atoi(invoiceIdParam)
	if err != nil {
		return handleError(c, http.StatusInternalServerError, err)
	}

	data, err := h.invoiceUsecase.IndexLine(limit, offset, invoiceId)
	if err != nil {

	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": data,
		"msg":  "get succsess",
	})
}

func (h *handler) GetInvoiceLine(c echo.Context) error {
	id := transformIdToInt(c)

	data, err := h.invoiceUsecase.GetInvoiceLine(id)
	if err != nil {
		handleError(c, http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": data,
		"msg":  "get succsess",
	})
}

func (h *handler) CreateInvoiceLine(c echo.Context) error {
	var request model.InvoiceLine
	err := json.NewDecoder(c.Request().Body).Decode(&request)
	if err != nil {
		return handleError(c, http.StatusInternalServerError, err)
	}

	data, err := h.invoiceUsecase.CreateInvoiceLine(request)
	if err != nil {
		return handleError(c, http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": data,
		"msg":  "create succsess",
	})
}

func (h *handler) UpdatedInvoiceLine(c echo.Context) error {
	//get param
	id := transformIdToInt(c)
	var request model.InvoiceLine
	err := json.NewDecoder(c.Request().Body).Decode(&request)
	if err != nil {
		return handleError(c, http.StatusInternalServerError, err)
	}

	//run function

	data, err := h.invoiceUsecase.UpdatedInvoiceLine(id, request, request.ProductID)
	if err != nil {
		return handleError(c, http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": data,
		"msg":  "create succsess",
	})
}

func (h *handler) DeleteInvoiceLine(c echo.Context) error {
	//get param
	id := transformIdToInt(c)

	data, err := h.invoiceUsecase.DeleteInvoiceLine(id)
	if err != nil {
		return handleError(c, http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": data,
		"msg":  "data deleted",
	})
}
