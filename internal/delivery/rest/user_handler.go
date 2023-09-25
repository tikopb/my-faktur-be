package rest

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *handler) Getuser(c echo.Context) error {
	Id := transformIdToInt(c)

	data, err := h.productUsecase.GetProduct(Id)
	if err != nil {
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	return handleError(c, http.StatusOK, errors.New("GET SUCCESS"), meta, data)
}
