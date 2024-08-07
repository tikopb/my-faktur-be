package fileService

import (
	"bemyfaktur/internal/model"
	"mime/multipart"
)

type Repository interface {
	//file64 format
	GetFileList64(model.FileServiceRequest) ([]model.FileServiceRespont, error)
	SaveFile64(model.FileServiceRequest) (model.FileServiceRespont, error)

	//file format
	GetFileList(model.FileServiceRequest) ([]model.FileServiceRespont, error)
	SaveFile(model.FileServiceRequest, *multipart.Form) (model.FileServiceRespont, error)

	//drop file
	DeleteFile(model.FileServiceRequest) (model.FileServiceRespont, error)

	//Get File Url
	GetUrlFile(model.FileServiceRequest) ([]model.FileServiceRespont, error)

	//update and delete file use for V1 fileservice update
	DeleteAndUpdateV1(request model.FileServiceRequest, requestDeleted model.FileServiceRespont, form *multipart.Form) (model.FileServiceRespont, error)
}
