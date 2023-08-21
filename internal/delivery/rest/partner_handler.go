package rest

import (
	"bemyfaktur/internal/model"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func (h *handler) IndexPartner(c echo.Context) error {

	data, err := h.partnerUsecase.IndexPartner()
	if err != nil {
		fmt.Printf("got error %s\n", err.Error())

		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": data,
	})
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

		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": data,
	})
}

func (h *handler) CreatePartner(c echo.Context) error {
	var request model.Partner
	err := json.NewDecoder(c.Request().Body).Decode(&request)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
	}
	data, err := h.partnerUsecase.CreatePartner(request)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"erorr": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": data,
	})
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
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
	}

	data, err := h.partnerUsecase.UpdatedPartner(Id, request)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
	}

	//return value
	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": data,
		"msg":  "dataUpdated",
	})
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
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
	}

	//return value
	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": data,
	})
}
