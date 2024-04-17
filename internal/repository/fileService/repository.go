package fileservice

import "bemyfaktur/internal/model"

type Repository interface {
	GetFileList(model.FileServiceRequest) ([]model.FileServiceRespont, error)
	SaveFile(model.FileServiceRequest) (model.FileServiceRespont, error)
	DeleteFile(model.FileServiceRequest) (model.FileServiceRespont, error)
}
