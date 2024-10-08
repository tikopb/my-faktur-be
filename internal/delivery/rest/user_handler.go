package rest

import (
	"bemyfaktur/internal/model"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *handler) Getuser(c echo.Context) error {
	param := c.QueryParam("username")
	meta = nil
	data, err := h.authUsecase.GetUser(param)
	if err != nil {
		WriteLogErorr("[delivery][rest][user_handler][Getuser] Get user failed", err)
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	return handleError(c, http.StatusOK, errors.New("GET SUCCESS"), meta, data)
}

func (h *handler) RegisterUser(c echo.Context) error {
	var request model.RegisterRequest
	err := json.NewDecoder(c.Request().Body).Decode(&request)
	if err != nil {
		WriteLogErorr("[delivery][rest][user_handler][RegisterUser] register failed ", err)
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	meta = nil
	userData, err := h.authUsecase.RegisterUser(request)
	if err != nil {
		WriteLogErorr("[delivery][rest][user_handler][RegisterUser] register failed ", err)
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	return handleError(c, 200, errors.New("REGISTER SUCCESS"), meta, userData)
}

func (h *handler) Login(c echo.Context) error {
	var request model.LoginRequest
	err := json.NewDecoder(c.Request().Body).Decode(&request)
	if err != nil {
		WriteLogErorr("[delivery][rest][user_handler][Login] ", err)
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	meta = nil
	sessionData, err := h.authUsecase.Login(request)
	if err != nil {
		WriteLogErorr("[delivery][rest][user_handler][Login] ", err)
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	msg := "[delivery][rest][user_handler][Login] LOGIN SUCCESS " + sessionData.UserInformation.Username
	WriteLogInfo(msg)

	errMsg := errors.New("AUTHORIZED")
	//add msg to FE
	if sessionData.UserSession.OrganizationID == 0 {
		errMsg = errors.New("Organization Not Exist journey, continue to start creating organization")
	}

	return handleError(c, http.StatusOK, errMsg, meta, sessionData)
}

func (h *handler) RefreshSession(c echo.Context) error {
	var request model.UserSession

	refreshToken, err := h.middleware.GetValueParamHeader(c.Request(), "Refresh-token")
	if err != nil {
		WriteLogErorr("[delivery][rest][user_handler][RefreshSession] ", err)
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	//decalre refresh token on sessioin is
	request.RefreshToken = refreshToken

	//clear the meta
	meta = nil
	sessionData, err := h.authUsecase.RefreshToken(refreshToken)
	if err != nil {
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	msg := "[delivery][rest][user_handler][RefreshSession] Refresh SUCCESS "
	WriteLogInfo(msg)
	return handleError(c, http.StatusOK, errors.New("AUTHORIZED"), meta, sessionData)
}

func (h *handler) LogOut(c echo.Context) error {

	//get header information
	refreshToken, err := h.middleware.GetValueParamHeader(c.Request(), "Refresh-token")
	if err != nil {
		WriteLogErorr("[delivery][rest][user_handler][LogOut] ", err)
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}
	accessToken, err := h.middleware.GetSessionData(c.Request())
	if err != nil {
		WriteLogErorr("[delivery][rest][user_handler][LogOut] ", err)
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	request := model.UserSession{
		AccessToken:  accessToken.AccessToken,
		RefreshToken: refreshToken,
	}

	meta = nil
	h.authUsecase.LogOutUser(request)

	msg := "[delivery][rest][user_handler][RefreshSession] log out success"
	WriteLogInfo(msg)
	return handleError(c, http.StatusOK, errors.New("LOG OUT SUCCESS"), meta, nil)
}

func (h *handler) CorsTest(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": "success",
	})
}
