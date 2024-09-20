package invoice

import (
	"bemyfaktur/internal/model"
	"bemyfaktur/internal/model/constant"
	"bemyfaktur/internal/repository/invoice"
	"bemyfaktur/internal/repository/partner"
	"bemyfaktur/internal/repository/product"
	"bemyfaktur/internal/usecase/fileservice"
	"errors"
	"mime/multipart"
	"time"

	"github.com/google/uuid"
)

type invoiceUsecase struct {
	invoiceRepo invoice.InvoiceRepositoryInterface
	partnerRepo partner.Repository
	productRepo product.Repository
	fileService fileservice.Usecase
}

func GetUsecase(invoiceRepo invoice.InvoiceRepositoryInterface, partnerRepo partner.Repository, productRepo product.Repository, fileService fileservice.Usecase) InvoiceUsecaseInterface {
	return &invoiceUsecase{
		invoiceRepo: invoiceRepo,
		partnerRepo: partnerRepo,
		productRepo: productRepo,
		fileService: fileService,
	}
}

// CreateInvoice implements InvoiceUsecaseInterface.
func (iu *invoiceUsecase) CreateInvoice(request model.InvoiceRequest, userID string) (model.InvoiceRespont, error) {
	//# Get all data First
	data := model.InvoiceRespont{}

	//set Partner
	partnerData, err := iu.partnerRepo.ShowInternal(request.PartnerUUID)
	if err != nil || !partnerData.Isactive {
		return data, errors.New("partner not exist")
	}
	request.PartnerId = partnerData.ID

	//setCreatedBy && updateBy
	request.CreatedById = userID
	request.UpdatedById = userID

	preloadData, err := iu.invoiceRepo.Create(request, partnerData)
	if err != nil {
		return model.InvoiceRespont{}, err
	}

	return preloadData, nil
}

// DeleteInvoice implements InvoiceUsecaseInterface.
func (iu *invoiceUsecase) DeleteInvoice(id uuid.UUID) (string, error) {
	return iu.invoiceRepo.Delete(id)
}

// GetInvoice implements InvoiceUsecaseInterface.
func (iu *invoiceUsecase) GetInvoice(id uuid.UUID) (model.InvoiceRespont, error) {

	preloadData, err := iu.invoiceRepo.Show(id)

	//file service
	fileParam := model.FileServiceRequest{
		UuidDoc: preloadData.ID,
	}
	fileUrl, err := iu.fileService.GetFileUrl(fileParam)
	if err != nil {
		return model.InvoiceRespont{}, err
	}
	preloadData.File = fileUrl

	return preloadData, nil
}

// IndexInvoice implements InvoiceUsecaseInterface.
func (iu *invoiceUsecase) IndexInvoice(limit int, offset int, q string, order []string, dateFrom string, dateTo string) ([]model.InvoiceRespont, error) {
	return iu.invoiceRepo.Index(limit, offset, q, order, dateFrom, dateTo)
}

// UpdatedInvoice implements InvoiceUsecaseInterface.
/**
DOD (Definition Of Done)
1. 	validated partner
1.a validated for partner first
1.b if partner exist and isactive == true then passed
2.	insert data of invoice with struct invoice
3. 	fill invoice respon for return
*/
func (iu *invoiceUsecase) UpdatedInvoice(id uuid.UUID, request model.InvoiceRequest, userId string) (model.InvoiceRespont, error) {

	//set Partner
	partnerData, err := iu.partnerRepo.ShowInternal(request.PartnerUUID)
	if err != nil || !partnerData.Isactive {
		return model.InvoiceRespont{}, errors.New("partner not exist or inactive")
	}
	request.PartnerId = partnerData.ID

	//set updated
	request.UpdatedById = userId
	return iu.invoiceRepo.Update(id, request)
}

//-- invoice line part

// CreateInvoiceLine implements InvoiceUsecaseInterface.
func (iu *invoiceUsecase) CreateInvoiceLine(request model.InvoiceLineRequest, userId string) (model.InvoiceLineRespont, error) {
	//set createdby
	request.CreatedById = userId

	//validated the product
	product, err := iu.validateProduct(request.ProductUUID)
	if err != nil {
		return model.InvoiceLineRespont{}, err
	}
	request.ProductID = product.ID //declare the product id of int

	//validated header not void or not complete
	invoice, err := iu.validateInvoice(request.InvoiceUUId)
	if err != nil {
		return model.InvoiceLineRespont{}, err
	}
	request.InvoiceId = invoice.ID

	return iu.invoiceRepo.CreateLine(request)
}

// DeleteInvoiceLine implements InvoiceUsecaseInterface.
func (iu *invoiceUsecase) DeleteInvoiceLine(id uuid.UUID) (string, error) {
	return iu.invoiceRepo.DeleteLine(id)
}

// GetInvoiceLine implements InvoiceUsecaseInterface.
func (iu *invoiceUsecase) GetInvoiceLine(id uuid.UUID) (model.InvoiceLineRespont, error) {
	return iu.invoiceRepo.ShowLine(id)
}

// UpdatedInvoiceLine implements InvoiceUsecaseInterface.
/*
	Definition Of Done
	- Validate Product:
		- Before using a product, ensure that its 'isActive' flag is set to true.
		- If the product is not valid, an expected error function should be triggered.
	- validate invoice:
		- invoice must be active and in draft posisition
	- Calculate Line Amount:
		- This function calculates the line amount based on the quantity and discounts applied.
		- The calculation considers whether the discount is a percentage or a fixed amount.
		- The result is the product of (quantity * unit price) minus the discount.
*/
func (iu *invoiceUsecase) UpdatedInvoiceLine(id uuid.UUID, request model.InvoiceLineRequest) (model.InvoiceLineRespont, error) {
	//validated the product
	//validated the product
	product, err := iu.validateProduct(request.ProductUUID)
	if err != nil {
		return model.InvoiceLineRespont{}, err
	}
	request.ProductID = product.ID //declare the product id of int

	invoice, err := iu.validateInvoice(request.InvoiceUUId)
	if err != nil {
		return model.InvoiceLineRespont{}, err
	}
	request.InvoiceId = invoice.ID

	return iu.invoiceRepo.UpdateLine(id, request)
}

// IndexLine implements InvoiceUsecaseInterface.
func (iu *invoiceUsecase) IndexLine(limit int, offset int, invoiceId uuid.UUID, q string, order []string) ([]model.InvoiceLineRespont, error) {
	//validated header not void or not complete
	invoice, err := iu.invoiceRepo.ShowInternal(invoiceId)
	if err != nil {
		return []model.InvoiceLineRespont{}, err
	}
	return iu.invoiceRepo.IndexLine(limit, offset, invoice.ID, q, order)
}

// validated product is activated on production
func (iu *invoiceUsecase) validateProduct(id uuid.UUID) (model.Product, error) {
	product, err := iu.productRepo.ShowInternal(id)
	if err != nil {
		return model.Product{}, err
	}
	if !product.IsActive {
		return model.Product{}, errors.New("product is not activated, please review the data")
	}
	return product, nil
}

// validatet the invoiced header data
func (iu *invoiceUsecase) validateInvoice(id uuid.UUID) (model.Invoice, error) {
	invoice, err := iu.invoiceRepo.ShowInternal(id)
	if err != nil {
		return model.Invoice{}, err
	}
	if invoice.Status == constant.InvoiceStatusVoid {
		return model.Invoice{}, errors.New("invoice status is void, please review the data")
	} else if invoice.Status == constant.InvoiceStatusProcessed {
		return model.Invoice{}, errors.New("invoice status is in progress, please review the data")
	} else if invoice.Status == constant.InvoiceStatusComplete {
		return model.Invoice{}, errors.New("invoice status is complete, please review the data")
	}

	return invoice, nil
}

func (iu *invoiceUsecase) HandlingPagination(q string, limit int, offset int, dateFrom string, dateTo string) (int64, error) {
	count, err := iu.invoiceRepo.HandlingPagination(q, limit, offset, dateFrom, dateTo)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (iu *invoiceUsecase) HandlingPaginationLine(q string, limit int, offset int, invoiceId int) (int64, error) {
	count, err := iu.invoiceRepo.HandlingPaginationLine(q, limit, offset, invoiceId)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// CreateInvoiceV2 implements InvoiceUsecaseInterface.
func (iu *invoiceUsecase) CreateInvoiceV2(request model.InvoiceRequestV2, userId string) (model.InvoiceRespontV2, error) {
	//set Partner
	partnerData, err := iu.partnerRepo.ShowInternal(request.Header.PartnerUUID)
	if err != nil || !partnerData.Isactive {
		return model.InvoiceRespontV2{}, errors.New("partner not exist")
	}
	request.Header.PartnerId = partnerData.ID

	//setCreatedBy && updateBy
	request.Header.CreatedById = userId
	request.Header.UpdatedById = userId

	//validate product
	//validated the product
	lines := request.Line
	linesPost := []model.InvoiceLineRequest{}
	for _, line := range lines {
		//validate product
		product, err := iu.validateProduct(line.ProductUUID)
		if err != nil {
			return model.InvoiceRespontV2{}, err
		}
		line.ProductID = product.ID //declare the product id of int

		//set createdby
		line.CreatedById = userId
		line.UpdatedById = userId
		line.OrganizationId = request.Header.OrganizationId
		linesPost = append(linesPost, line)
	}

	header, line, err := iu.invoiceRepo.CreateInvoiceV2(request.Header, linesPost, partnerData)
	if err != nil {
		return model.InvoiceRespontV2{}, err
	}

	data := model.InvoiceRespontV2{
		Header: header,
		Line:   line,
	}

	return data, nil
}

// CreateInvoiceV2 implements InvoiceUsecaseInterface.
func (iu *invoiceUsecase) CreateInvoiceV3(request model.InvoiceRequestV2, userId string) (model.InvoiceRespontV2, error) {
	//set Partner
	partnerData, err := iu.partnerRepo.ShowInternal(request.Header.PartnerUUID)
	if err != nil || !partnerData.Isactive {
		return model.InvoiceRespontV2{}, errors.New("partner not exist")
	}
	request.Header.PartnerId = partnerData.ID

	//setCreatedBy && updateBy
	request.Header.CreatedById = userId
	request.Header.UpdatedById = userId

	//validate product
	//validated the product
	lines := request.Line
	linesPost := []model.InvoiceLineRequest{}
	for _, line := range lines {
		//validate product
		product, err := iu.validateProduct(line.ProductUUID)
		if err != nil {
			return model.InvoiceRespontV2{}, err
		}
		line.ProductID = product.ID //declare the product id of int

		//set createdby
		line.CreatedById = userId
		line.UpdatedById = userId
		linesPost = append(linesPost, line)
	}

	header, line, err := iu.invoiceRepo.CreateInvoiceV2(request.Header, linesPost, partnerData)
	if err != nil {
		return model.InvoiceRespontV2{}, err
	}

	data := model.InvoiceRespontV2{
		Header: header,
		Line:   line,
	}

	return data, nil
}

// Update Invoicev3 implements InvoiceUsecaseInterface.
func (iu *invoiceUsecase) UpdateInvoiceV3(id uuid.UUID, request model.InvoiceRequest, form multipart.Form) (model.InvoiceRespontV3, error) {

	header, err := iu.invoiceRepo.Update(id, request)
	if err != nil {
		return model.InvoiceRespontV3{}, err
	}

	data := model.InvoiceRespontV2{
		Header: header,
		Line:   header.Line,
	}

	//start update the document
	//just update when form is updated
	returnData := model.InvoiceRespontV3{}
	imageActionStr := form.Value["image_action"][0]
	if imageActionStr == string(constant.FileActionUpdate) {
		if len(form.File["files"]) > 0 {
			fileRequest := model.FileServiceRequest{
				UuidDoc:   data.Header.ID,
				DocType:   "INV",
				CreatedBy: data.Header.CreatedBy.UserId,
			}
			files, err := iu.fileService.DeleteAndUpdateV1(fileRequest, &form)
			if err != nil {
				return model.InvoiceRespontV3{}, err
			}

			returnData = model.InvoiceRespontV3{
				Data: data,
				File: files,
			}
		}
	}

	return returnData, nil
}

// Partial implements InvoiceUsecaseInterface.
func (iu *invoiceUsecase) Partial(partner_id uuid.UUID, q string) ([]model.InvoicePartialRespont, error) {
	//get partner id
	partner, err := iu.partnerRepo.ShowInternal(partner_id)
	if err != nil {
		return []model.InvoicePartialRespont{}, err
	}

	return iu.invoiceRepo.Partial(partner.ID, q)
}

func (iu *invoiceUsecase) StatusUpdateV3(id uuid.UUID, userId string, docAction constant.InvoiceDocAction) (model.InvoiceRespont, error) {
	//Get the data of invoiceData internal
	invoiceData, err := iu.invoiceRepo.ShowInternal(id)
	if err != nil {
		return model.InvoiceRespont{}, err
	}

	//set value of doaction
	invoiceData.DocAction = docAction
	invoiceData.UpdatedBy = userId

	if invoiceData.PayDate.IsZero() {
		invoiceData.PayDate = time.Now()
	}

	requestData, err := iu.invoiceRepo.ParsingInvoiceToInvoiceRequest(invoiceData)
	if err != nil {
		return model.InvoiceRespont{}, err
	}

	//run the function for update
	data, err := iu.invoiceRepo.Update(id, requestData)
	if err != nil {
		return model.InvoiceRespont{}, err
	}

	return data, nil
}
