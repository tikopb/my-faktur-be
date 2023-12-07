package rest

import (
	"bemyfaktur/internal/model"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

var idStr string

func (h *handler) IndexProduct(c echo.Context) error {
	//set param
	limit, offset := HandlingLimitAndOffset(c)

	//get parameter
	q := c.QueryParam("q")

	//setOrderData
	order, err := h.GetOrderClauses(c)
	if err != nil {
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	data, err := h.productUsecase.IndexProduct(limit, offset, q, order)
	if err != nil {
		WriteLogErorr("[delivery][rest][product_handler][IndexProduct] ", err)
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	//meta data field
	searchParams := model.GetSeatchParamProduct()

	meta, err = h.pgUtilRepo.PaginationUtil("products", searchParams, limit, offset, q, "", "")
	if err != nil {
		WriteLogErorr("[delivery][rest][product_handler][IndexProduct] ", err)
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	//add validation if partner == null then set roolback
	if len(data) < 1 {
		data = []model.ProductRespon{}
	}

	return handleError(c, http.StatusOK, errors.New("GET SUCCESS"), meta, data)
}

func (h *handler) GetProduct(c echo.Context) error {
	// Parse the string into a UUID
	productId, err := h.parsingId(c)
	if err != nil {
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	data, err := h.productUsecase.GetProduct(productId)
	if err != nil {
		WriteLogErorr("[delivery][rest][product_handler][IndexProduct] ", err)
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	return handleError(c, http.StatusOK, errors.New("GET SUCCESS"), meta, data)
}

func (h *handler) CreateProduct(c echo.Context) error {
	var request model.Product

	//getUserId
	userId, err := h.middleware.GetuserId(c.Request())
	if err != nil {
		WriteLogErorr("[delivery][rest][product_handler][CreateProduct] ", err)
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	err = json.NewDecoder(c.Request().Body).Decode(&request)
	if err != nil {
		WriteLogErorr("[delivery][rest][product_handler][CreateProduct] ", err)
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	data, err := h.productUsecase.CreateProduct(request, userId)
	if err != nil {
		WriteLogErorr("[delivery][rest][product_handler][CreateProduct] ", err)
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	meta = nil
	return handleError(c, http.StatusCreated, errors.New("CREATE "+data.Name+" SUCCESS"), meta, data)
}

func (h *handler) UpdatedProduct(c echo.Context) error {
	var request model.Product
	// get param
	Id, err := h.parsingId(c)
	if err != nil {
		WriteLogErorr("[delivery][rest][product_handler][UpdatedProduct] ", err)
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	//run function
	err = json.NewDecoder(c.Request().Body).Decode(&request)
	if err != nil {
		WriteLogErorr("[delivery][rest][product_handler][UpdatedProduct] ", err)
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	data, err := h.productUsecase.UpdatedProduct(Id, request)
	if err != nil {
		WriteLogErorr("[delivery][rest][product_handler][UpdatedProduct] ", err)
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	return handleError(c, http.StatusOK, errors.New("UPDATE "+data.Name+" SUCCESS"), meta, data)
}

func (h *handler) DeleteProduct(c echo.Context) error {
	// get param
	Id, err := h.parsingId(c)
	if err != nil {
		WriteLogErorr("[delivery][rest][product_handler][UpdatedProduct] ", err)
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	data, err := h.productUsecase.DeleteProduct(Id)
	if err != nil {
		WriteLogErorr("[delivery][rest][product_handler][DeleteProduct] ", err)
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	return handleError(c, http.StatusOK, errors.New("DELETE "+data+" SUCCESS"), meta, data)
}

func (h *handler) parsingId(c echo.Context) (uuid.UUID, error) {
	idStr = c.Param("id")
	// Parse the string into a UUID
	productId, err := h.ParsingUUID(idStr)
	if err != nil {
		return uuid.UUID{}, err
	}

	return productId, nil
}

func (h *handler) PartialProduct(c echo.Context) error {
	//get parameter
	q := c.QueryParam("q")
	q = strings.ToLower(q)

	data, err := h.productUsecase.Partial(q)
	if err != nil {
		WriteLogErorr("[delivery][rest][PartialProduct] ", err)

		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	//meta data field
	meta = nil
	msg := "GET SUCCESS"
	if len(data) == 0 {
		msg = "data not found"
	}

	return handleError(c, http.StatusOK, errors.New(msg), meta, data)
}
