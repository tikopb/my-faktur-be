package rest

import (
	"bemyfaktur/internal/model"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func handleError(c echo.Context, statusCode int, err error) error {
	return c.JSON(statusCode, map[string]interface{}{
		"error": err.Error(),
		"msg":   "internal error" + err.Error(),
	})
}

func transformIdToInt(c echo.Context) int {
	// get param
	PartnerId := c.Param("id")
	Id, err := strconv.Atoi(PartnerId)
	if err != nil {
		panic(err)
	}
	return Id
}

func (h *handler) IndexProduct(c echo.Context) error {
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

	data, err := h.productUsecase.IndexPartner(limit, offset)
	if err != nil {
		return handleError(c, http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": data,
	})
}

func (h *handler) GetProduct(c echo.Context) error {
	Id := transformIdToInt(c)

	data, err := h.productUsecase.GetProduct(Id)
	if err != nil {
		return handleError(c, http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": data,
	})
}

func (h *handler) CreateProduct(c echo.Context) error {
	var request model.Product
	err := json.NewDecoder(c.Request().Body).Decode(&request)
	if err != nil {
		return handleError(c, http.StatusInternalServerError, err)
	}

	data, err := h.productUsecase.CreateProduct(request)
	if err != nil {
		return handleError(c, http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": data,
	})
}

func (h *handler) UpdatedProduct(c echo.Context) error {
	var request model.Product
	// get param
	Id := transformIdToInt(c)

	//run function
	err := json.NewDecoder(c.Request().Body).Decode(&request)
	if err != nil {
		return handleError(c, http.StatusInternalServerError, err)
	}

	data, err := h.productUsecase.UpdatedProduct(Id, request)
	if err != nil {
		return handleError(c, http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": data,
		"msg":  "data updated",
	})
}

func (h *handler) DeleteProduct(c echo.Context) error {
	// get param
	Id := transformIdToInt(c)

	data, err := h.productUsecase.DeleteProduct(Id)
	if err != nil {
		return handleError(c, http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": data,
		"msg":  "data deleted",
	})
}
