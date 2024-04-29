package fileservice

import "bemyfaktur/internal/model"

type Repository interface {
	SaveFile64([]model.FileServiceRequest) ([]model.FileServiceRespont, error)
	GetFileList64(model.FileServiceRequest) ([]model.FileServiceRespont, error)

	//file format
	SaveFile([]model.FileServiceRequest) ([]model.FileServiceRespont, error)
	GetFileList(model.FileServiceRequest) ([]model.FileServiceRespont, error)

	//drop the meesage
	DeleteFile([]model.FileServiceRequest) ([]model.FileServiceRespont, error)
}