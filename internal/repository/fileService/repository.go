package fileservice

import "bemyfaktur/internal/model"

type Repository interface {
	SaveFile([]model.FileServiceRequest) ([]model.FileServiceRespont, error)
	GetFile(model.FileServiceRequest) ([]model.FileServiceRespont, error)
	DeleteFile([]model.FileServiceRequest) error
}
