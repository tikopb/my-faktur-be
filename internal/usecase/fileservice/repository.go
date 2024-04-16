package fileservice

import "bemyfaktur/internal/model"

type Repository interface {
	SaveFile([]model.FileServiceRequest) ([]model.FileServiceRespont, error)
	GetFileList(model.FileServiceRequest) ([]model.FileServiceRespont, error)
	DeleteFile([]model.FileServiceRequest) ([]model.FileServiceRespont, error)
}
