package rest

import (
	"bemyfaktur/internal/model"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func (h *handler) IndexInvoice(c echo.Context) error {
	//set param
	limit, offset := HandlingLimitAndOffset(c)

	//get parameter
	q := c.QueryParam("q")
	q = strings.ToLower(q)

	//setOrderData
	order, err := h.GetOrderClauses(c)
	if err != nil {
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	//get dateFrom and date to
	dateFrom, dateTo, err := h.HandlingDateFromAndDateTo(c)
	if err != nil {
		WriteLogErorr("[delivery][rest][invoice_handler][IndexInvoice] ", err)
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	data, err := h.invoiceUsecase.IndexInvoice(limit, offset, q, order, dateFrom, dateTo)
	if err != nil {
		WriteLogErorr("[delivery][rest][invoice_handler][IndexInvoice] ", err)
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	//meta data field
	count, err := h.invoiceUsecase.HandlingPagination(q, limit, offset, dateFrom, dateTo)
	if err != nil {
		WriteLogErorr("[delivery][rest][invoice_handler][IndexInvoice] ", err)
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	meta, err = h.pgUtilRepo.PaginationUtilWithJoinTable(int64(count), limit, offset)
	if err != nil {
		WriteLogErorr("[delivery][rest][invoice_handler][IndexInvoice] ", err)
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	return handleError(c, http.StatusOK, errors.New("GET SUCCESS"), meta, data)
}

func (h *handler) GetInvoice(c echo.Context) error {
	id, err := h.parsingId(c)
	if err != nil {
		WriteLogErorr("[delivery][rest][invoice_handler][GetInvoice] ", err)
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	data, err := h.invoiceUsecase.GetInvoice(id)
	if err != nil {
		WriteLogErorr("[delivery][rest][invoice_handler][GetInvoice] ", err)
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	return handleError(c, http.StatusOK, errors.New("GET SUCCESS"), meta, data)
}

func (h *handler) CreateInvoice(c echo.Context) error {
	var request model.InvoiceRequest
	err := json.NewDecoder(c.Request().Body).Decode(&request)
	if err != nil {
		WriteLogErorr("[delivery][rest][invoice_handler][CreateInvoice] ", err)
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	//getUserId
	userId, err := h.middleware.GetuserId(c.Request())
	if err != nil {
		WriteLogErorr("[delivery][rest][invoice_handler][CreateInvoice] ", err)
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	data, err := h.invoiceUsecase.CreateInvoice(request, userId)
	if err != nil {
		WriteLogErorr("[delivery][rest][invoice_handler][CreateInvoice] ", err)
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}
	return handleError(c, http.StatusOK, errors.New("CREATE "+data.BatchNo+" SUCCESS"), meta, data)
}

func (h *handler) UpdateInvoice(c echo.Context) error {
	var request model.InvoiceRequest
	//get param
	id, err := h.parsingId(c)
	if err != nil {
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	//run function
	err = json.NewDecoder(c.Request().Body).Decode(&request)
	if err != nil {
		WriteLogErorr("[delivery][rest][invoice_handler][UpdateInvoice] ", err)
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	//getUpdateByUserId
	userId, err := h.middleware.GetuserId(c.Request())
	if err != nil {
		WriteLogErorr("[delivery][rest][invoice_handler][CreateInvoice] ", err)
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	data, err := h.invoiceUsecase.UpdatedInvoice(id, request, userId)
	if err != nil {
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	return handleError(c, http.StatusOK, errors.New("UPDATE "+data.BatchNo+" SUCCESS"), meta, data)

}

func (h *handler) DeleteInvoice(c echo.Context) error {
	//get param
	id, err := h.parsingId(c)
	if err != nil {
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	data, err := h.invoiceUsecase.DeleteInvoice(id)
	if err != nil {
		WriteLogErorr("[delivery][rest][invoice_handler][DeleteInvoice] ", err)
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	return handleError(c, http.StatusOK, errors.New("DELETE "+data+"SUCCESS"), meta, data)

}

// generate data invoice header and line in same time
func (h *handler) CreateInvoiceV2(c echo.Context) error {
	var request model.InvoiceRequestV2
	err := json.NewDecoder(c.Request().Body).Decode(&request)
	if err != nil {
		WriteLogErorr("[delivery][rest][invoice_handler][CreateInvoiceV2] ", err)
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	//getUserId
	userId, err := h.middleware.GetuserId(c.Request())
	if err != nil {
		WriteLogErorr("[delivery][rest][invoice_handler][CreateInvoiceV2] ", err)
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	data, err := h.invoiceUsecase.CreateInvoiceV2(request, userId)
	if err != nil {
		WriteLogErorr("[delivery][rest][invoice_handler][CreateInvoiceV2] ", err)
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	return handleError(c, http.StatusOK, errors.New("CREATE "+data.Header.BatchNo+" SUCCESS"), meta, data)
}

/*
generate data invoice header and line in same time with form-data format! and handling file in one file at a time
api handler base on form-data format
*/
func (h *handler) CreateInvoiceV3(c echo.Context) error {
	// File validation
	form, err := c.MultipartForm()
	if err != nil {
		return err
	}

	var request model.InvoiceRequestV2
	err = json.Unmarshal([]byte(c.Request().FormValue("data")), &request)
	if err != nil {
		WriteLogErorr("[delivery][rest][invoice_handler][CreateInvoiceV2] ", err)
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	//getUserId
	userInf, err := h.middleware.GetUserInformation(c.Request())
	if err != nil {
		WriteLogErorr("[delivery][rest][invoice_handler][CreateInvoiceV2] ", err)
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	//set org_id
	request.Header.OrganizationId = userInf.OrganizationID
	data, err := h.invoiceUsecase.CreateInvoiceV2(request, userInf.UserId)
	if err != nil {
		WriteLogErorr("[delivery][rest][invoice_handler][CreateInvoiceV2] ", err)
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	//extract the form-data format
	fileRequest := model.FileServiceRequest{
		File:      nil,
		File64:    nil,
		UuidDoc:   data.Header.ID,
		DocType:   "INV",
		FileName:  "",
		CreatedBy: userInf.UserId,
	}

	//send to usecase of file service
	dataFile, err := h.fileserviceUsecase.SaveFile(fileRequest, form)
	if err != nil {
		return handleError(c, http.StatusInternalServerError, err, nil, nil)
	}

	returnData := model.InvoiceRespontV3{
		Data: data,
		File: dataFile,
	}

	return handleError(c, http.StatusOK, errors.New("CREATE "+data.Header.BatchNo+" SUCCESS"), meta, returnData)
}

func (h *handler) Partialnvoice(c echo.Context) error {
	//get parameter
	q := c.QueryParam("q")
	q = strings.ToLower(q)

	//get partner_id
	partner_id := c.QueryParam("partner_id")
	if partner_id == "" {
		err := errors.New("partner_id must be filled")
		WriteLogErorr("[delivery][rest][invoice_handler][InvoicePartial] ", err)
		return handleError(c, http.StatusNotFound, err, nil, nil)
	}
	//parsing to uuid, give erorr when failed.
	uuidPartnerID, err := uuid.Parse(partner_id)
	if err != nil {
		err := errors.New("partner_id must be a valid UUID")
		WriteLogErorr("[delivery][rest][invoice_handler][InvoicePartial] ", err)
		return handleError(c, http.StatusNotFound, err, nil, nil)
	}

	data, err := h.invoiceUsecase.Partial(uuidPartnerID, q)
	if err != nil {
		WriteLogErorr("[delivery][rest][invoice_handler][InvoicePartial] ", err)
		return handleError(c, http.StatusInternalServerError, err, nil, nil)
	}

	return handleError(c, http.StatusOK, errors.New("Get SUCCESS"), meta, data)
}
