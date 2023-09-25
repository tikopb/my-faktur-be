package rest

import (
	"bemyfaktur/internal/model"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func (h *handler) IndexPartner(c echo.Context) error {

	data, err := h.partnerUsecase.IndexPartner()
	if err != nil {
		fmt.Printf("got error %s\n", err.Error())

		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	return handleError(c, http.StatusOK, errors.New("GET SUCCESS"), meta, data)

}

func (h *handler) GetPartner(c echo.Context) error {
	PartnerId := c.Param("id")
	Id, err := strconv.Atoi(PartnerId)
	if err != nil {
		panic(err)
	}

	data, err := h.partnerUsecase.GetPartner(Id)
	if err != nil {
		fmt.Printf("got error %s\n", err.Error())

		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	return handleError(c, http.StatusOK, errors.New("GET SUCCESS"), meta, data)

}

func (h *handler) CreatePartner(c echo.Context) error {
	var request model.Partner
	err := json.NewDecoder(c.Request().Body).Decode(&request)
	if err != nil {
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}
	data, err := h.partnerUsecase.CreatePartner(request)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"erorr": err.Error(),
		})
	}

	return handleError(c, http.StatusOK, errors.New("CREATE "+data.Name+" SUCCESS"), meta, data)

}

func (h *handler) UpdatedPartner(c echo.Context) error {
	//get param
	PartnerId := c.Param("id")
	Id, err := strconv.Atoi(PartnerId)
	if err != nil {
		panic(err)
	}

	//search Data and validate
	var request model.Partner
	err = json.NewDecoder(c.Request().Body).Decode(&request)
	if err != nil {
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	data, err := h.partnerUsecase.UpdatedPartner(Id, request)
	if err != nil {
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	//return value
	return handleError(c, http.StatusOK, errors.New("Update "+data.Name+" SUCCESS"), meta, data)
}

func (h *handler) DeletePartner(c echo.Context) error {
	//get param
	PartnerId := c.Param("id")
	Id, err := strconv.Atoi(PartnerId)
	if err != nil {
		panic(err)
	}

	data, err := h.partnerUsecase.Deletepartner(Id)
	if err != nil {
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	//return value
	return handleError(c, http.StatusOK, errors.New("DELETE SUCCESS"), meta, data)

}
