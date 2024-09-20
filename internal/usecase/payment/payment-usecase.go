package payment

import (
	"bemyfaktur/internal/model"
	"bemyfaktur/internal/model/constant"
	"bemyfaktur/internal/repository/invoice"
	"bemyfaktur/internal/repository/partner"
	"bemyfaktur/internal/repository/payment"
	"bemyfaktur/internal/usecase/fileservice"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"mime/multipart"
	"time"
)

type paymentUsecase struct {
	paymentRepo payment.PaymentRepositoryinterface
	invoiceRepo invoice.InvoiceRepositoryInterface
	partnerRepo partner.Repository
	fileService fileservice.Usecase
}

func GetUsecase(paymentRepo payment.PaymentRepositoryinterface, invoiceRepo invoice.InvoiceRepositoryInterface, partnerRepo partner.Repository, fileService fileservice.Usecase) PaymentUsecaseInterface {
	return &paymentUsecase{
		paymentRepo: paymentRepo,
		invoiceRepo: invoiceRepo,
		partnerRepo: partnerRepo,
		fileService: fileService,
	}
}

// Createpayment header payment part
// Createpayment implements PaymentUsecaseInterface.
func (pu *paymentUsecase) Createpayment(request model.PaymentRequest, userId string) (model.PaymentRespont, error) {
	request.CreatedBy = userId
	request.UpdatedBy = userId

	//validate the partner
	partner, err := pu.partnerRepo.ShowInternal(request.PartnerUUID)
	if err != nil || !partner.Isactive {
		return model.PaymentRespont{}, errors.New("partner not exist or inactived")
	}

	request.PartnerID = partner.ID

	return pu.paymentRepo.Create(request)
}

// Deletepayment implements PaymentUsecaseInterface.
func (pu *paymentUsecase) Deletepayment(id uuid.UUID) (string, error) {
	return pu.paymentRepo.Delete(id)
}

// Getpayment implements PaymentUsecaseInterface.
func (pu *paymentUsecase) Getpayment(id uuid.UUID) (model.PaymentRespont, error) {
	return pu.paymentRepo.Show(id)
}

// Indexpayment implements PaymentUsecaseInterface.
func (pu *paymentUsecase) Indexpayment(limit int, offset int, q string, order []string, dateFrom string, dateTo string) ([]model.PaymentRespont, error) {
	return pu.paymentRepo.Index(limit, offset, q, order, dateFrom, dateTo)
}

// Updatedpayment implements PaymentUsecaseInterface.
func (pu *paymentUsecase) Updatedpayment(id uuid.UUID, request model.PaymentRequest) (model.PaymentRespont, error) {

	return pu.paymentRepo.Update(id, request)
}

// invoice line part
// CreateInvoiceLine implements PaymentUsecaseInterface.
func (pu *paymentUsecase) CreatePaymentLine(request model.PaymentLineRequest, userId string) (model.PaymentLineRespont, error) {
	data := model.PaymentLineRespont{}
	invoice, err := pu.invoiceRepo.ShowInternal(request.Invoice_uuid)
	if err != nil {
		return data, err
	} else if invoice.Status != constant.InvoiceStatusComplete {
		return data, errors.New("invoice not in completed")
	}
	//set the value to invoice_id because relation key used with id int not uuid
	request.Invoice_id = invoice.ID

	header, err := pu.paymentRepo.ShowInternal(request.PaymentUUID)
	if err != nil {
		return data, err
	}
	request.PaymentID = header.ID

	//set created by and updated by
	request.CreatedBy = userId
	request.UpdatedBy = userId

	//return value
	return pu.paymentRepo.CreateLine(request)
}

// GetInvoiceLine implements PaymentUsecaseInterface.
func (pu *paymentUsecase) GetPaymentLine(id uuid.UUID) (model.PaymentLineRespont, error) {
	return pu.paymentRepo.ShowLine(id)
}

// IndexLine implements PaymentUsecaseInterface.
func (pu *paymentUsecase) IndexLine(limit int, offset int, paymentUUID uuid.UUID, q string) ([]model.PaymentLineRespont, int, error) {
	header, err := pu.paymentRepo.ShowInternal(paymentUUID)
	if err != nil {
		return []model.PaymentLineRespont{}, 0, err
	}
	//getting the payment_id in integer
	paymentId := header.ID

	data, err := pu.paymentRepo.IndexLine(limit, offset, paymentId)
	return data, paymentId, err
}

// UpdatedInvoiceLine implements PaymentUsecaseInterface.
func (pu *paymentUsecase) UpdatedPaymentLine(id uuid.UUID, request model.PaymentLineRequest) (model.PaymentLineRespont, error) {
	data := model.PaymentLineRespont{}
	invoice, err := pu.invoiceRepo.ShowInternal(request.Invoice_uuid)
	if err != nil {
		return data, err
	} else if invoice.Status != constant.InvoiceStatusComplete {
		return data, errors.New("invoice not in completed")
	}
	//set invoice_id
	request.Invoice_id = invoice.ID

	//return value
	return pu.paymentRepo.UpdateLine(id, request)
}

// DeleteInvoiceLine implements PaymentUsecaseInterface.
func (pu *paymentUsecase) DeletePaymentLine(id uuid.UUID) (string, error) {
	return pu.paymentRepo.DeleteLine(id)
}

// CreateV2 implements PaymentUsecaseInterface.
func (pu *paymentUsecase) CreatePaymentV2(request model.PaymentRequestV2, userId string) (model.PaymentRespont, error) {
	//header validation
	request.Header.CreatedBy = userId
	request.Header.UpdatedBy = userId
	// validate the partner
	partner, err := pu.partnerRepo.ShowInternal(request.Header.PartnerUUID)
	if err != nil || !partner.Isactive {
		return model.PaymentRespont{}, errors.New("partner not exist or inactived")
	}

	request.Header.PartnerID = partner.ID

	//line validation
	linesPost := []model.PaymentLineRequest{}
	for _, line := range request.Line {
		//validate invoice
		invoice, err := pu.invoiceRepo.ShowInternal(line.Invoice_uuid)
		fmt.Println(err)
		if err == nil {
			//set the value to invoice_id because relation key used with id int not uuid
			line.Invoice_id = invoice.ID

			//set created by
			line.CreatedBy = userId
			line.UpdatedBy = userId
			line.OrganizationId = invoice.OrganizationId
			linesPost = append(linesPost, line)
		} else if err.Error() == "data not found" {
			//if err !nil then return erorr
			return model.PaymentRespont{}, errors.New("data of invoice not found")
		} else if invoice.Status != constant.InvoiceStatusComplete {
			return model.PaymentRespont{}, errors.New("invoice not in completed")
		} else if err != nil {
			//if err !nil then return erorr
			return model.PaymentRespont{}, err
		}
	}

	//define repeat line and change it into linespost it will change the data with data that already validate
	request.Line = linesPost

	//hit the request
	data, err := pu.paymentRepo.CreateV2(request)
	if err != nil {
		return model.PaymentRespont{}, err
	}

	//set return value of data
	return data, nil
}

func (pu *paymentUsecase) HandlingPagination(q string, limit int, offset int, dateFrom string, dateTo string) (int64, error) {
	count, err := pu.paymentRepo.HandlingPagination(q, limit, offset, dateFrom, dateTo)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (pu *paymentUsecase) HandlingPaginationLine(q string, limit int, offset int, paymentID int) (int64, error) {
	count, err := pu.paymentRepo.HandlingPaginationLine(q, limit, offset, paymentID)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// post data for payment with v3 api versioning
func (pu *paymentUsecase) PostPaymentV3(request model.PaymentRequestV2, userID string, form *multipart.Form) (model.PaymentRespontV3, error) {
	//run createpaymentV2
	data, err := pu.CreatePaymentV2(request, userID)
	if err != nil {
		return model.PaymentRespontV3{}, err
	}

	//file service start
	fileRequest := model.FileServiceRequest{
		File:       nil,
		File64:     nil,
		UuidDoc:    data.UUID,
		DocType:    "PAY",
		FileName:   "",
		CreatedBy:  data.CreatedBy.UserId,
		FileAction: "",
	}

	//set to usecase of file service formation
	dataFile, err := pu.fileService.SaveFile(fileRequest, form)
	if err != nil {
		return model.PaymentRespontV3{}, err
	}

	return model.PaymentRespontV3{data, dataFile}, nil
}

// set the update for payment with v3 api versioning
func (pu *paymentUsecase) UpdatePaymentV3(id uuid.UUID, request model.PaymentRequest, form *multipart.Form) (model.PaymentRespontV3, error) {
	//run on process update
	data, err := pu.paymentRepo.Update(id, request)
	if err != nil {
		return model.PaymentRespontV3{}, err
	}

	//start update the document
	//just update when form is updated
	returnData := model.PaymentRespontV3{}
	imageActionStr := form.Value["image_action"][0]
	if imageActionStr == string(constant.FileActionUpdate) {
		if len(form.File["files"]) > 0 {
			fileRequest := model.FileServiceRequest{
				UuidDoc:   data.UUID,
				DocType:   "INV",
				CreatedBy: data.CreatedBy.UserId,
			}
			files, err := pu.fileService.DeleteAndUpdateV1(fileRequest, form)
			if err != nil {
				return model.PaymentRespontV3{}, err
			}

			returnData = model.PaymentRespontV3{
				Data: data,
				File: files,
			}
		}
	}

	//set the return when no erorr
	return returnData, nil
}

/*
* update for docaction function
* run on getting data internal
* set value of data
* run the function of update
 */
func (pu *paymentUsecase) StatusUpdateV3(id uuid.UUID, userId string, docAction constant.PaymentDocAction) (model.PaymentRespont, error) {
	//get the data of paymentData internal
	paymentData, err := pu.paymentRepo.ShowInternal(id)
	if err != nil {
		return model.PaymentRespont{}, err
	}

	//set value of docaction
	paymentData.DocAction = docAction
	paymentData.UpdatedBy = userId

	if paymentData.PayDate.IsZero() {
		paymentData.PayDate = time.Now()
	}

	requestData, err := pu.paymentRepo.ParsingPaymentToPaymentRequest(paymentData)
	if err != nil {
		return model.PaymentRespont{}, err
	}

	//run the function for update
	data, err := pu.paymentRepo.Update(id, requestData)
	if err != nil {
		return model.PaymentRespont{}, err
	}

	return data, nil
}
