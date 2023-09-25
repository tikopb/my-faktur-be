package rest

import (
	"bemyfaktur/internal/model"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func (h *handler) IndexProduct(c echo.Context) error {

	limitStr := c.QueryParam("limit")
	if limitStr == "" {
		limitStr = "15" // Default value
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	offsetStr := c.QueryParam("offset")
	if offsetStr == "" {
		offsetStr = "0"
	}
	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	//get parameter
	q := c.QueryParam("q")

	data, err = h.productUsecase.IndexProduct(limit, offset, q)
	if err != nil {
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	meta, err = h.PaginationUtil(h.getTableName(), h.getSeatchParam(), limit, offset, q)
	if err != nil {
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	return handleError(c, http.StatusOK, errors.New("GET SUCCESS"), meta, data)
}

func (h *handler) GetProduct(c echo.Context) error {
	Id := transformIdToInt(c)

	data, err := h.productUsecase.GetProduct(Id)
	if err != nil {
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	return handleError(c, http.StatusOK, errors.New("GET SUCCESS"), meta, data)
}

func (h *handler) CreateProduct(c echo.Context) error {
	var request model.Product
	err := json.NewDecoder(c.Request().Body).Decode(&request)
	if err != nil {
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	data, err := h.productUsecase.CreateProduct(request)
	if err != nil {
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	return handleError(c, http.StatusOK, errors.New("CREATE "+data.Name+" SUCCESS"), meta, data)
}

func (h *handler) UpdatedProduct(c echo.Context) error {
	var request model.Product
	// get param
	Id := transformIdToInt(c)

	//run function
	err := json.NewDecoder(c.Request().Body).Decode(&request)
	if err != nil {
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	data, err := h.productUsecase.UpdatedProduct(Id, request)
	if err != nil {
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	return handleError(c, http.StatusOK, errors.New("UPDATE "+data.Name+" SUCCESS"), meta, data)
}

func (h *handler) DeleteProduct(c echo.Context) error {
	// get param
	Id := transformIdToInt(c)

	data, err := h.productUsecase.DeleteProduct(Id)
	if err != nil {
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	return handleError(c, http.StatusOK, errors.New("DELETE "+data+" SUCCESS"), meta, data)
}

func (h *handler) getTableName() string {
	return "products"
}

func (h *handler) getSeatchParam() []string {
	searchParam := []string{"name", "description"}
	return searchParam
}
