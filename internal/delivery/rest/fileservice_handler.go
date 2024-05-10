package rest

import (
	"bemyfaktur/internal/model"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func (h *handler) UploadFile(c echo.Context) error {

	//getUserId
	userId, err := h.middleware.GetuserId(c.Request())
	if err != nil {
		WriteLogErorr("[delivery][rest][fileservice_handler][UploadFile] ", err)
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	// Get document data request field
	uuidString := c.FormValue("uuid_doc")
	parsedUUID, err := uuid.Parse(uuidString)
	if err != nil {
		WriteLogErorr("[delivery][rest][fileservice_handler][UploadFile] ", err)
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	//field the request
	request := model.FileServiceRequest{
		File:      nil,
		File64:    nil,
		UuidDoc:   parsedUUID,
		DocType:   c.FormValue("doctype"),
		FileName:  "",
		CreatedBy: userId,
	}

	// File validation
	form, err := c.MultipartForm()
	if err != nil {
		return err
	}

	//send to usecase of file service
	data, err = h.fileserviceUsecase.SaveFile(request, form)
	if err != nil {
		return handleError(c, http.StatusInternalServerError, err, nil, nil)
	}
	return handleError(c, http.StatusOK, errors.New("file success uploaded"), nil, data)
}

func (h *handler) GetTheFileBaseUrl(c echo.Context) error {
	filename := c.QueryParam("name")
	return c.File("./assets/" + filename)
}
