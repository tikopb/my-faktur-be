package rest

import (
	"bemyfaktur/internal/model"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

func (h *handler) IndexPartner(c echo.Context) error {
	//set param
	limit, offset := HandlingLimitAndOffset(c)

	//get parameter
	q := c.QueryParam("q")

	//setOrderData
	order, err := h.GetOrderClauses(c)
	if err != nil {
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	data, err := h.partnerUsecase.IndexPartner(limit, offset, q, order)
	if err != nil {
		//fmt.Printf("got error %s\n", err.Error())
		WriteLogErorr("[delivery][rest][IndexPartner] ", err)
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	//meta data field
	searchParams := model.GetSeatchParamPartner()
	meta, err = h.pgUtilRepo.PaginationUtil("partners", searchParams, limit, offset, q, "", "")
	if err != nil {
		WriteLogErorr("[delivery][rest][IndexPartner] ", err)
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}
	//meta data field
	msg := "GET SUCCESS"

	//add validation if partner == null then set roolback
	if len(data) < 1 {
		data = []model.PartnerRespon{}
	}

	return handleError(c, http.StatusOK, errors.New(msg), meta, data)
}

func (h *handler) GetPartner(c echo.Context) error {
	// Get the "id" parameter from the request
	meta = nil
	idStr := c.Param("id")

	// Parse the string into a UUID
	partnerID, err := h.ParsingUUID(idStr)
	if err != nil {
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	data, err := h.partnerUsecase.GetPartner(partnerID)
	if err != nil {
		WriteLogErorr("[delivery][rest][GetPartner] ", err)
		return handleError(c, http.StatusNotFound, err, meta, data)
	}

	return handleError(c, http.StatusOK, errors.New("GET SUCCESS"), meta, data)

}

func (h *handler) CreatePartner(c echo.Context) error {
	var request model.Partner

	//getUserId
	userId, err := h.middleware.GetuserId(c.Request())
	if err != nil {
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	err = json.NewDecoder(c.Request().Body).Decode(&request)
	if err != nil {
		WriteLogErorr("[delivery][rest][CreatePartner] ", err)
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}
	data, err := h.partnerUsecase.CreatePartner(request, userId)
	if err != nil {
		WriteLogErorr("[delivery][rest][CreatePartner] ", err)
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	meta = nil
	return handleError(c, http.StatusOK, errors.New("CREATE "+data.Name+" SUCCESS"), meta, data)

}

func (h *handler) UpdatedPartner(c echo.Context) error {
	//get param
	idStr := c.Param("id")

	// Parse the string into a UUID
	partnerID, err := h.ParsingUUID(idStr)
	if err != nil {
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	//search Data and validate
	var request model.Partner
	err = json.NewDecoder(c.Request().Body).Decode(&request)
	if err != nil {
		WriteLogErorr("[delivery][rest][UpdatedPartner] ", err)
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	data, err := h.partnerUsecase.UpdatedPartner(partnerID, request)
	if err != nil {
		WriteLogErorr("[delivery][rest][UpdatedPartner] ", err)
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	//return value
	meta = nil
	return handleError(c, http.StatusOK, errors.New("Update "+data.Name+" SUCCESS"), meta, data)
}

func (h *handler) DeletePartner(c echo.Context) error {
	//get param
	idStr := c.Param("id")

	// Parse the string into a UUID
	partnerID, err := h.ParsingUUID(idStr)
	if err != nil {
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	data, err := h.partnerUsecase.Deletepartner(partnerID)
	if err != nil {
		WriteLogErorr("[delivery][rest][DeletePartner] ", err)
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	//return value
	meta = nil
	return handleError(c, http.StatusOK, errors.New("DELETE SUCCESS"), meta, data)

}

func (h *handler) PartialPartner(c echo.Context) error {
	//get parameter
	q := c.QueryParam("q")
	q = strings.ToLower(q)

	data, err := h.partnerUsecase.PartialGet(q)
	if err != nil {
		WriteLogErorr("[delivery][rest][PartialPartner] ", err)

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
