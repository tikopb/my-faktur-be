package rest

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *handler) Getuser(c echo.Context) error {
	Id := transformIdToInt(c)

	data, err := h.productUsecase.GetProduct(Id)
	if err != nil {
		return handleError(c, http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": data,
	})
}
