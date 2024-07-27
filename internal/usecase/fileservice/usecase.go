package fileservice

import (
	"bemyfaktur/internal/model"
	"mime/multipart"
)

type Usecase interface {
	SaveFile64([]model.FileServiceRequest) ([]model.FileServiceRespont, error)
	GetFileList64(model.FileServiceRequest) ([]model.FileServiceRespont, error)

	//file format
	SaveFile(model.FileServiceRequest, *multipart.Form) (model.FileServiceRespont, error)
	GetFileList(model.FileServiceRequest) ([]model.FileServiceRespont, error)

	//drop the meesage
	DeleteFile([]model.FileServiceRequest) ([]model.FileServiceRespont, error)

	//getFileUrl
	GetFileUrl(model.FileServiceRequest) ([]model.FileServiceRespont, error)

	//Delete and update v1
	DeleteAndUpdateV1(request model.FileServiceRequest, requestDeleted []model.FileServiceRespont, form *multipart.Form) (model.FileServiceRespont, error)
}
