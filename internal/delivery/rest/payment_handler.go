package rest

import (
	"bemyfaktur/internal/model"
	"encoding/json"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *handler) IndexPayment(c echo.Context) error {

	//set param
	limit, offset := HandlingLimitAndOffset(c)

	data, err := h.paymentUsecase.Indexpayment(limit, offset)
	if err != nil {
		return handleError(c, http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": data,
		"msg":  "get succsess",
	})

}

func (h *handler) Getpayment(c echo.Context) error {
	id := transformIdToInt(c)

	data, err := h.paymentUsecase.Getpayment(id)
	if err != nil {
		return handleError(c, http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": data,
		"msg":  "get succsess",
	})
}

func (h *handler) CreatePayment(c echo.Context) error {
	var request model.PaymentRequest
	err := json.NewDecoder(c.Request().Body).Decode(&request)
	if err != nil {
		return handleError(c, http.StatusInternalServerError, err)
	}

	data, err := h.paymentUsecase.Createpayment(request)
	if err != nil {
		return handleError(c, http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": data,
		"msg":  "get succsess",
	})

}

func (h *handler) UpdatePayment(c echo.Context) error {
	var request model.PaymentRequest
	id := transformIdToInt(c)
	err := json.NewDecoder(c.Request().Body).Decode(&request)
	if err != nil {
		return handleError(c, http.StatusInternalServerError, err)
	}

	data, err := h.paymentUsecase.Updatedpayment(id, request)
	if err != nil {
		return handleError(c, http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": data,
		"msg":  "get succsess",
	})
}

func (h *handler) DeletePayment(c echo.Context) error {
	id := transformIdToInt(c)
	data, err := h.paymentUsecase.Deletepayment(id)
	if err != nil {
		return handleError(c, http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": "DETELE SUCCSESS:" + data,
		"msg":  "delete succsess",
	})

}

//line
