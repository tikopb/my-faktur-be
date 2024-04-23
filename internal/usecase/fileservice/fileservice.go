package fileservice

import (
	"bemyfaktur/internal/model"
	fileservice "bemyfaktur/internal/repository/fileService"
	"errors"
)

type fileServiceUsecase struct {
	fileServiceRepo fileservice.Repository
}

func GetRepository(fileServiceRepo fileservice.Repository) Repository {
	return &fileServiceUsecase{
		fileServiceRepo: fileServiceRepo,
	}
}

// GetFileList implements Repository.
func (f *fileServiceUsecase) GetFileList(request model.FileServiceRequest) ([]model.FileServiceRespont, error) {
	return f.fileServiceRepo.GetFileList(request)
}

// SaveFile implements Repository.
func (f *fileServiceUsecase) SaveFile([]model.FileServiceRequest) ([]model.FileServiceRespont, error) {
	panic("unimplemented")
}

/** Byte 64 format
*
**/
// GetFileList implements Repository.
func (f *fileServiceUsecase) GetFileList64(request model.FileServiceRequest) ([]model.FileServiceRespont, error) {
	//get the repository of getFile
	return f.fileServiceRepo.GetFileList64(request)
}

// SaveFile implements Repository.
func (f *fileServiceUsecase) SaveFile64(requests []model.FileServiceRequest) ([]model.FileServiceRespont, error) {
	if len(requests) > 5 {
		return []model.FileServiceRespont{}, errors.New("file maximal that can be save just 5 document")
	}

	//prepare return data list
	dataReturnList := []model.FileServiceRespont{}
	for _, request := range requests {
		data, err := f.fileServiceRepo.SaveFile(request)
		if err != nil {
			return []model.FileServiceRespont{}, err
		}

		datas := model.FileServiceRespont{
			File64:   data.File64,
			FileName: data.FileName,
		}

		dataReturnList = append(dataReturnList, datas)
	}

	//return method
	return dataReturnList, nil
}

// DeleteFile implements Repository.
func (f *fileServiceUsecase) DeleteFile(requests []model.FileServiceRequest) ([]model.FileServiceRespont, error) {
	returnPartsings := []model.FileServiceRespont{}
	for _, request := range requests {
		data, err := f.fileServiceRepo.DeleteFile(request)
		if err != nil {
			return []model.FileServiceRespont{}, err
		}

		parsingData := model.FileServiceRespont{
			FileName: data.FileName,
		}

		returnPartsings = append(returnPartsings, parsingData)
	}

	return returnPartsings, nil
}
