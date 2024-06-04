package rest

import (
	"bemyfaktur/internal/model"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *handler) CreateOrganization(c echo.Context) error {
	var request model.OrganizationRequest
	err := json.NewDecoder(c.Request().Body).Decode(&request)
	if err != nil {
		WriteLogErorr("[delivery][rest][organization_handler][CreateOrganization] ", err)
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	//getuserID
	userId, err := h.middleware.GetuserId(c.Request())
	if err != nil {
		WriteLogErorr("[delivery][rest][organization_handler][CreateOrganization] ", err)
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}
	request.UserId = userId

	data, err := h.organizationUsecase.Create(request)
	if err != nil {
		WriteLogErorr("[delivery][rest][organization_handler][CreateOrganization] ", err)
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	return handleError(c, http.StatusOK, errors.New("create success"), nil, data)
}

func (h *handler) GetOrganization(c echo.Context) error {
	//getuserInformation
	user, err := h.middleware.GetUserInformation(c.Request())
	if err != nil {
		WriteLogErorr("[delivery][rest][organization_handler][GetOrganization] ", err)
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	//set the value and get request
	request := model.OrganizationRequest{
		OrgCode: user.OrganizationUUID,
	}
	data, err := h.organizationUsecase.Show(request)
	if err != nil {
		WriteLogErorr("[delivery][rest][organization_handler][GetOrganization] ", err)
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	return handleError(c, http.StatusOK, errors.New("get success"), meta, data)
}

func (h *handler) DeleteOrganization(c echo.Context) error {

	//getuserInformation
	user, err := h.middleware.GetUserInformation(c.Request())
	if err != nil {
		WriteLogErorr("[delivery][rest][organization_handler][GetOrganization] ", err)
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}
	//set the value and get request
	request := model.OrganizationRequest{
		OrgCode: user.OrganizationUUID,
	}
	data, err := h.organizationUsecase.Delete(request)
	if err != nil {
		WriteLogErorr("[delivery][rest][organization_handler][GetOrganization] ", err)
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	return handleError(c, http.StatusOK, errors.New("get success"), meta, data)

}
