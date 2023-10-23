package rest

import (
	"bemyfaktur/internal/model"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *handler) IndexProduct(c echo.Context) error {
	//set param
	limit, offset := HandlingLimitAndOffset(c)

	//get parameter
	q := c.QueryParam("q")

	data, err := h.productUsecase.IndexProduct(limit, offset, q)
	if err != nil {
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	//meta data field
	searchParams := model.GetSeatchParamProduct()

	meta, err = h.pgUtilRepo.PaginationUtil("products", searchParams, limit, offset, q, "", "")
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

	//getUserId
	userId, err := h.middleware.GetuserId(c.Request())
	if err != nil {
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	err = json.NewDecoder(c.Request().Body).Decode(&request)
	if err != nil {
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	data, err := h.productUsecase.CreateProduct(request, userId)
	if err != nil {
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	meta = nil
	return handleError(c, http.StatusCreated, errors.New("CREATE "+data.Name+" SUCCESS"), meta, data)
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
