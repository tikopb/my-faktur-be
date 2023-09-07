package rest

import (
	"bemyfaktur/internal/model"
	"encoding/json"
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
		return handleError(c, http.StatusInternalServerError, err)
	}

	data, err := h.paymentUsecase.IndexLine(limit, offset, paymentId, q)
	if err != nil {
		return handleError(c, http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": data,
		"msg":  "get succsess",
	})
}

func (h *handler) GetPaymentLine(c echo.Context) error {
	id := transformIdToInt(c)

	data, err := h.paymentUsecase.GetPaymentLine(id)
	if err != nil {
		return handleError(c, http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": data,
		"msg":  "get siccsess",
	})
}

func (h *handler) CreatePaymentLine(c echo.Context) error {
	//set param
	var request model.PaymentLineRequest
	err := json.NewDecoder(c.Request().Body).Decode(&request)
	if err != nil {
		return handleError(c, http.StatusInternalServerError, err)
	}

	//run function
	data, err := h.paymentUsecase.CreatePaymentLine(request)
	if err != nil {
		return handleError(c, http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": data,
		"msg":  "create succsess",
	})
}

func (h *handler) UpdatePaymentLine(c echo.Context) error {
	//set param
	id := transformIdToInt(c)
	var request model.PaymentLineRequest
	err := json.NewDecoder(c.Request().Body).Decode(&request)
	if err != nil {
		return handleError(c, http.StatusInternalServerError, err)
	}

	//run function
	data, err := h.paymentUsecase.UpdatedPaymentLine(id, request)
	if err != nil {
		return handleError(c, http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": data,
		"msg":  "update succsess",
	})
}

func (h *handler) DeletePaymentLine(c echo.Context) error {
	//set param
	id := transformIdToInt(c)

	data, err := h.paymentUsecase.DeletePaymentLine(id)
	if err != nil {
		return handleError(c, http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": "data deleted" + data,
		"msg":  "delete succsess",
	})
}
