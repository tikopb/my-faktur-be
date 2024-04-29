package rest

import "github.com/labstack/echo/v4"

func (h *handler) UploadFile(c echo.Context) error {

	panic("")
}

func (h *handler) GetTheFileBaseUrl(c echo.Context) error {
	filename := c.QueryParam("name")
	return c.File("./assets/" + filename)
}
