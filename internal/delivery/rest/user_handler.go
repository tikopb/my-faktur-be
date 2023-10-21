package rest

import (
	"bemyfaktur/internal/model"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *handler) Getuser(c echo.Context) error {
	Id := transformIdToInt(c)

	meta = nil
	data, err := h.productUsecase.GetProduct(Id)
	if err != nil {
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	return handleError(c, http.StatusOK, errors.New("GET SUCCESS"), meta, data)
}

func (h *handler) RegisterUser(c echo.Context) error {
	var request model.RegisterRequest
	err := json.NewDecoder(c.Request().Body).Decode(&request)
	if err != nil {
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	meta = nil
	userData, err := h.authUsecase.RegisterUser(request)
	if err != nil {
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	return handleError(c, 200, errors.New("REGISTER SUCCESS"), meta, userData)
}

func (h *handler) Login(c echo.Context) error {
	var request model.LoginRequest
	err := json.NewDecoder(c.Request().Body).Decode(&request)
	if err != nil {
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	meta = nil
	sessionData, err := h.authUsecase.Login(request)
	if err != nil {
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	return handleError(c, http.StatusOK, errors.New("AUTHORIZED"), meta, sessionData)
}

func (h *handler) RefreshSession(c echo.Context) error {
	var request model.UserSession
	err := json.NewDecoder(c.Request().Body).Decode(&request)
	if err != nil {
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	meta = nil
	sessionData, err := h.authUsecase.RefreshToken(request)
	if err != nil {
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	return handleError(c, http.StatusOK, errors.New("AUTHORIZED"), meta, sessionData)
}

func (h *handler) LogOut(c echo.Context) error {
	var request model.UserSession
	err := json.NewDecoder(c.Request().Body).Decode(&request)
	if err != nil {
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	meta = nil
	h.authUsecase.LogOutUser(request)

	return handleError(c, http.StatusOK, errors.New("LOG OUT SUCCESS"), meta, request)
}
