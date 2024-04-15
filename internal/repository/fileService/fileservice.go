package fileservice

import (
	"bemyfaktur/internal/model"
	"encoding/base64"
	"path/filepath"
	"strings"

	"gorm.io/gorm"
)

type fileserviceRepo struct {
	db *gorm.DB
}

func GetRepository(db *gorm.DB) Repository {
	return &fileserviceRepo{
		db: db,
	}
}

// GetFile implements Repository.
func (f *fileserviceRepo) GetFile(model.FileServiceRequest) ([]model.FileServiceRespont, error) {

	panic("unimplemented")
}

// SaveFile implements Repository.
func (f *fileserviceRepo) SaveFile([]model.FileServiceRequest) ([]model.FileServiceRespont, error) {
	panic("unimplemented")
}

// DeleteFile implements Repository.
func (f *fileserviceRepo) DeleteFile([]model.FileServiceRequest) error {
	panic("unimplemented")
}

func (f *fileserviceRepo) DecodedFile(data string) ([]byte, error) {
	decodedBytes, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		panic("")
	}
	return decodedBytes, nil
}

func (f *fileserviceRepo) EncodedFile(data string) (string, error) {
	encodedBytes := base64.StdEncoding.EncodeToString([]byte(data))
	return encodedBytes, nil
}

func (f *fileserviceRepo) IsValidFileExtension(filename string) bool {
	validExtensions := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".pdf":  true,
		".png":  true,
	}
	ext := strings.ToLower(filepath.Ext(filename))
	_, isValid := validExtensions[ext]
	return isValid
}
