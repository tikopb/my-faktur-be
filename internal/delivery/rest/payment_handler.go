package rest

import (
	"bemyfaktur/internal/model"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

func (h *handler) IndexPayment(c echo.Context) error {

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
		WriteLogErorr("[delivery][rest][payment_handler][IndexPayment] ", err)
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	data, err := h.paymentUsecase.Indexpayment(limit, offset, q, order, dateFrom, dateTo)
	if err != nil {
		WriteLogErorr("[delivery][rest][payment_handler][IndexPayment] ", err)
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	//meta data field
	count, err := h.paymentUsecase.HandlingPagination(q, limit, offset, dateFrom, dateTo)
	if err != nil {
		WriteLogErorr("[delivery][rest][payment_handler][IndexPayment] ", err)
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	meta, err = h.pgUtilRepo.PaginationUtilWithJoinTable(int64(count), limit, offset)
	if err != nil {
		WriteLogErorr("[delivery][rest][payment_handler][IndexPayment] ", err)
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	return handleError(c, http.StatusOK, errors.New("GET SUCCESS"), meta, data)
}

func (h *handler) Getpayment(c echo.Context) error {
	id, err := h.parsingId(c)
	if err != nil {
		WriteLogErorr("[delivery][rest][payment_handler][Getpayment] ", err)
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	data, err := h.paymentUsecase.Getpayment(id)
	if err != nil {
		WriteLogErorr("[delivery][rest][payment_handler][Getpayment] ", err)
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	return handleError(c, http.StatusOK, errors.New("GET SUCCESS"), meta, data)
}

func (h *handler) CreatePayment(c echo.Context) error {
	var request model.PaymentRequest
	err := json.NewDecoder(c.Request().Body).Decode(&request)
	if err != nil {
		WriteLogErorr("[delivery][rest][payment_handler][CreatePayment] ", err)
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	//getUserId
	userId, err := h.middleware.GetuserId(c.Request())
	if err != nil {
		WriteLogErorr("[delivery][rest][payment_handler][CreatePayment] ", err)
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	data, err := h.paymentUsecase.Createpayment(request, userId)
	if err != nil {
		WriteLogErorr("[delivery][rest][payment_handler][CreatePayment] ", err)
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	return handleError(c, http.StatusOK, errors.New("CREATE "+data.BatchNo+" SUCCESS"), meta, data)

}

func (h *handler) UpdatePayment(c echo.Context) error {
	var request model.PaymentRequest
	id, err := h.parsingId(c)
	if err != nil {
		WriteLogErorr("[delivery][rest][payment_handler][UpdatePayment] ", err)
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	err = json.NewDecoder(c.Request().Body).Decode(&request)
	if err != nil {
		WriteLogErorr("[delivery][rest][payment_handler][UpdatePayment] ", err)
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	//getUserId
	userId, err := h.middleware.GetuserId(c.Request())
	if err != nil {
		WriteLogErorr("[delivery][rest][payment_handler][UpdatePayment] ", err)
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	request.UpdatedBy = userId
	request.CreatedBy = userId

	data, err := h.paymentUsecase.Updatedpayment(id, request)
	if err != nil {
		WriteLogErorr("[delivery][rest][payment_handler][UpdatePayment] ", err)
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	return handleError(c, http.StatusOK, errors.New("UPDATE "+data.BatchNo+" SUCCESS"), meta, data)

}

func (h *handler) DeletePayment(c echo.Context) error {
	id, err := h.parsingId(c)
	if err != nil {
		WriteLogErorr("[delivery][rest][payment_handler][DeletePayment] ", err)
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	data, err := h.paymentUsecase.Deletepayment(id)
	if err != nil {
		WriteLogErorr("[delivery][rest][payment_handler][DeletePayment] ", err)
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	return handleError(c, http.StatusOK, errors.New("DELETE "+data+" SUCCESS"), meta, data)
}

func (h *handler) CreatePaymentV2(c echo.Context) error {
	var request model.PaymentRequestV2
	err := json.NewDecoder(c.Request().Body).Decode(&request)
	if err != nil {
		WriteLogErorr("[delivery][rest][payment_handler][CreatePaymentV2]", err)
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	//getUserId
	userId, err := h.middleware.GetuserId(c.Request())
	if err != nil {
		WriteLogErorr("[delivery][rest][payment_handler][CreatePaymentV2] ", err)
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	data, err := h.paymentUsecase.CreatePaymentV2(request, userId)
	if err != nil {
		WriteLogErorr("[delivery][rest][payment_handler][CreatePaymentV2] ", err)
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	return handleError(c, http.StatusOK, errors.New("Create "+data.DocumentNo+" SUCCESS"), meta, data)
}

func (h *handler) CreatePaymentV3(c echo.Context) error {
	//file validation
	form, err := c.MultipartForm()
	if err != nil {
		return err
	}

	var request model.PaymentRequestV2
	err = json.Unmarshal([]byte(c.Request().FormValue("data")), &request)
	if err != nil {
		WriteLogErorr("[delivery][rest][invoice_handler][CreatePaymentV3] ", err)
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	//getUserId
	userInf, err := h.middleware.GetUserInformation(c.Request())
	if err != nil {
		WriteLogErorr("[delivery][rest][payment_handler][CreatePaymentV3] ", err)
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}
	request.Header.OrganizationId = userInf.OrganizationID

	//start the usercase process
	data, err := h.paymentUsecase.PostPaymentV3(request, userInf.UserId, form)
	if err != nil {
		WriteLogErorr("[delivery][rest][payment_handler][CreatePaymentV3] ", err)
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	return handleError(c, http.StatusOK, errors.New("Created Successfull "+data.Data.BatchNo+"with documentno"+data.Data.DocumentNo), meta, data)
}

func (h *handler) UpdatePaymentV3(c echo.Context) error {
	//file validation
	form, err := c.MultipartForm()
	if err != nil {
		return err
	}

	//get param
	id, err := h.parsingId(c)
	if err != nil {
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	var request model.PaymentRequest

	//getUserId
	userId, err := h.middleware.GetuserId(c.Request())
	if err != nil {
		WriteLogErorr("[delivery][rest][payment_handler][UpdatePayment] ", err)
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	request.UpdatedBy = userId
	request.CreatedBy = userId
	//run the function unmarshal from form to struct
	err = json.Unmarshal([]byte(c.Request().FormValue("data")), &request)
	if err != nil {
		WriteLogErorr("[delivery][rest][payment_handler][UpdatePaymentV3] ", err)
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	data, err := h.paymentUsecase.UpdatePaymentV3(id, request, form)
	if err != nil {
		WriteLogErorr("[delivery][rest][payment_handler][UpdatePaymentV3] ", err)
		return handleError(c, http.StatusInternalServerError, err, meta, data)
	}

	return handleError(c, http.StatusOK, errors.New("UPDATE "+data.Data.BatchNo+" SUCCESS"), meta, data)
}
