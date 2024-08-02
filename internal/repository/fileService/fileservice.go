package fileService

import (
	"bemyfaktur/internal/model"
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/labstack/gommon/log"
	"io"
	"io/fs"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type FileserviceRepo struct {
	db *gorm.DB
}

func GetRepository(db *gorm.DB) Repository {
	return &FileserviceRepo{
		db: db,
	}
}

/*
File  format
*/

// GetFileList implements Repository.
func (f *FileserviceRepo) GetFileList(request model.FileServiceRequest) ([]model.FileServiceRespont, error) {
	//Get File Path
	file_path, err := f.GetFilePathFromEnv()
	if err != nil {
		return []model.FileServiceRespont{}, err
	}

	//query searching for data list
	directorys, err := f.GetFromDb(request.UuidDoc, request.DocType)
	if err != nil {
		return []model.FileServiceRespont{}, err
	}

	datas := []model.FileServiceRespont{}
	for _, directory := range directorys {
		// Open file
		file, err := os.Open(fmt.Sprintf("%s%s", file_path, directory.FileName))
		if err != nil {
			return []model.FileServiceRespont{}, err
		}
		defer file.Close()

		fileInfo, err := f.GetFileStat(file)
		if err != nil {
			return []model.FileServiceRespont{}, err
		}

		data := model.FileServiceRespont{
			File:     file,
			FileName: fileInfo.Name(),
		}

		datas = append(datas, data)
	}

	return datas, nil
}

// SaveFile implements Repository.
func (f *FileserviceRepo) SaveFile(request model.FileServiceRequest, form *multipart.Form) (model.FileServiceRespont, error) {

	//Get File Path
	file_path, err := f.GetFilePathFromEnv()
	if err != nil {
		return model.FileServiceRespont{}, err
	}

	// Get the count of files
	fileCount := len(form.File["files"])

	//count validation
	if f.GetCountFromDb(request.UuidDoc, request.DocType)+int64(fileCount) > 5 {
		return model.FileServiceRespont{}, errors.New("max file is 5 file per document")
	}

	//loop the file
	for _, file := range form.File["files"] {
		// // Check extension
		ext := filepath.Ext(file.Filename)
		if !(strings.ToLower(ext) == ".jpeg" || strings.ToLower(ext) == ".jpg" || strings.ToLower(ext) == ".pdf" || strings.ToLower(ext) == ".png") {
			return model.FileServiceRespont{}, errors.New("just file format with .jpg, .jpeg, .png and .pdf is supported")
		}

		// Get rename file
		newFileName := f.GetRenameFile(file.Filename, request.DocType)

		// Source
		src, err := file.Open()
		if err != nil {
			return model.FileServiceRespont{}, err
		}
		defer src.Close()

		// Create file in assets directory
		dst, err := os.Create(fmt.Sprintf("%s%s", file_path, newFileName))
		if err != nil {
			return model.FileServiceRespont{}, err
		}
		defer dst.Close()

		// Copy
		if _, err = io.Copy(dst, src); err != nil {
			return model.FileServiceRespont{}, err
		}

		//ad to db
		err = f.AddToDb(request, newFileName)
		if err != nil {
			return model.FileServiceRespont{}, err
		}

	}

	datas, err := f.GetUrlFile(request)
	if err != nil {
		return model.FileServiceRespont{}, err
	}

	//return msg
	return datas[0], nil
}

/*
File 64 format
*/
// GetFile implements Repository.
func (f *FileserviceRepo) GetFileList64(request model.FileServiceRequest) ([]model.FileServiceRespont, error) {
	returnValuelist := []model.FileServiceRespont{}

	//query searching for data list
	data, err := f.GetFromDb(request.UuidDoc, request.DocType)
	if err != nil {
		return []model.FileServiceRespont{}, err

	}

	//decode
	for _, fileData := range data {
		fileEncoded, err := f.EncodeToFile64(fileData.FileName)
		if err != nil {
			return []model.FileServiceRespont{}, err
		}
		response := model.FileServiceRespont{
			File64:   fileEncoded,
			FileName: fileData.FileName,
		}

		returnValuelist = append(returnValuelist, response)
	}

	return returnValuelist, nil
}

// SaveFile implements Repository.
func (f *FileserviceRepo) SaveFile64(request model.FileServiceRequest) (model.FileServiceRespont, error) {
	//Get File Path
	file_path, err := f.GetFilePathFromEnv()
	if err != nil {
		return model.FileServiceRespont{}, err
	}

	//check the file extension
	validation, err := f.IsValidFile64(request.File64)
	if !validation {
		return model.FileServiceRespont{}, err
	}

	//get rename file and saved to file
	newFileName := f.GetRenameFile(request.FileName, request.DocType)

	//save file
	// Decode the base64 encoded file
	decodedFile, err := f.DecodedFromFile64(string(request.File64))
	if err != nil {
		return model.FileServiceRespont{}, err
	}

	// Construct the full path including the directory
	fullPath := fmt.Sprintf("%s%s", file_path, newFileName)

	// Write the file content to disk with the new filename
	err = os.WriteFile(fullPath, decodedFile, 0644)
	if err != nil {
		return model.FileServiceRespont{}, err
	}

	//save directory to db
	err = f.AddToDb(request, newFileName)
	if err != nil {
		return model.FileServiceRespont{}, err
	}

	//prepare return value
	returnDataList := model.FileServiceRespont{
		FileName: newFileName,
	}

	return returnDataList, nil
}

// DeleteFile implements Repository.
func (f *FileserviceRepo) DeleteFile(request model.FileServiceRequest) (model.FileServiceRespont, error) {
	//Get File Path
	file_path, err := f.GetFilePathFromEnv()
	if err != nil {
		return model.FileServiceRespont{}, err
	}

	returnDataLists := model.FileServiceRespont{}

	//getFilename
	filename := request.FileName
	fullPath := fmt.Sprintf("%s%s", file_path, filename)

	//delete from os env
	err = os.Remove(fullPath)
	if err != nil {
		return model.FileServiceRespont{}, err
	}

	//delete from db
	err = f.DeleteFromDb(filename)
	if err != nil {
		return model.FileServiceRespont{}, err
	}

	//prepare return value
	returnDataLists.FileName = filename + " deleted "

	return returnDataLists, nil
}

// GetUrlFile implements Repository. used for v1 version
func (f *FileserviceRepo) GetUrlFile(request model.FileServiceRequest) ([]model.FileServiceRespont, error) {
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()
	if err != nil {
		panic("env of " + "url_FileService" + "not found")
	}

	file_service_port := viper.GetString("file_service_port")
	url_FileService := viper.GetString("url_FileService")

	dataList, err := f.GetFromDb(request.UuidDoc, request.DocType)
	if err != nil {
		return []model.FileServiceRespont{}, err
	}

	returnData := []model.FileServiceRespont{}
	for _, data := range dataList {
		returnData = append(returnData, model.FileServiceRespont{
			FileName: data.FileName,
			FileUrl:  file_service_port + url_FileService + data.FileName,
		})
	}

	return returnData, nil
}

func (f *FileserviceRepo) DecodedFromFile64(fileBytes string) ([]byte, error) {
	// Trim leading and trailing whitespace
	fileBytes = strings.TrimSpace(fileBytes)

	// Check for empty input
	if len(fileBytes) == 0 {
		return nil, errors.New("input is empty")
	}

	// Decode the base64 encoded file
	decodedFile, err := base64.StdEncoding.DecodeString(fileBytes)
	if err != nil {
		return nil, fmt.Errorf("error decoding base64: %v", err)
	}

	return decodedFile, nil
}

func (f *FileserviceRepo) EncodeToFile64(filePath string) (string, error) {
	// Read the entire file content into a byte slice
	fileBytes, err := os.ReadFile(filePath)
	if err != nil {
		return "", errors.New("failed to read file: " + err.Error())
	}

	// Encode the byte slice to base64 string
	file64 := base64.StdEncoding.EncodeToString(fileBytes)

	return file64, nil
}

func (f *FileserviceRepo) GetFileStat(file *os.File) (fs.FileInfo, error) {
	// Get file extension
	fileInfo, err := file.Stat()
	if err != nil {
		return nil, err
	}

	return fileInfo, nil
}

/*
*
file with validation of jpg, jpeg, pdf and png that can be only save to system
*
*/
func (f *FileserviceRepo) IsValidFile64(fileBytes []byte) (bool, error) {
	// Define valid magic bytes for supported formats
	validMimeTypes := map[string][][]byte{
		"image/jpeg":      {{0xff, 0xd8, 0xff, 0xe0}},                                                                         // JPG start marker
		"image/png":       {{0x89, 0x50, 0x4E, 0x47}},                                                                         // PNG signature
		"application/pdf": {{0x25, 0x50, 0x44, 0x46}, {0x25, 0x25, 0x45, 0x4F, 0x46}, {0x25, 0x21}, {0x50, 0x4B, 0x03, 0x04}}, // PDF headers and additional bytes
	}

	// Get the first few bytes of the file
	if len(fileBytes) < 4 {
		return false, errors.New("file too small to determine type")
	}
	mimeType := http.DetectContentType(fileBytes[:4])

	// Check if mimeType is valid
	for validType, magicBytesList := range validMimeTypes {
		if mimeType == validType {
			for _, magicBytes := range magicBytesList {
				if len(fileBytes) >= len(magicBytes) && bytes.Equal(fileBytes[:len(magicBytes)], magicBytes) {
					return true, nil
				}
			}
			break // If mimetype matches but magic bytes don't, it's invalid
		}
	}

	return false, errors.New("invalid file format")
}

func (f *FileserviceRepo) GetRenameFile(originalFilename string, docType string) string {
	// Generate new filename with format: yyyymmdd-originalFilename
	newFilename := fmt.Sprintf("%s-%s-%s", time.Now().Format("060102150405"), docType, strings.TrimSpace(originalFilename))
	return newFilename
}

func (f *FileserviceRepo) SaveToFile(fileBytes []byte, filename string) error {
	//Get File Path
	file_path, err := f.GetFilePathFromEnv()
	if err != nil {
		return err
	}

	// Decode the base64 encoded file
	decodedFile, err := f.DecodedFromFile64(string(fileBytes))
	if err != nil {
		return err
	}

	// Construct the full path including the directory
	fullPath := fmt.Sprintf("%s%s", file_path, filename)

	// Write the file content to disk with the new filename
	err = os.WriteFile(fullPath, decodedFile, 0644)
	if err != nil {
		return err
	}

	return nil
}

/* get data infomration from db
 */
func (f *FileserviceRepo) GetFromDb(uuid uuid.UUID, docType string) ([]model.FileService, error) {
	// Get data from db using UUID
	var data []model.FileService
	if err := f.db.Where(&model.FileService{UuidDoc: uuid, DocType: docType}).Find(&data).Error; err != nil {
		return nil, err
	}

	return data, nil
}

/* get count infomration from db
 */
func (f *FileserviceRepo) GetCountFromDb(uuid uuid.UUID, docType string) int64 {
	var count int64
	if err := f.db.Model(&model.FileService{}).Where(&model.FileService{UuidDoc: uuid, DocType: docType}).Count(&count).Error; err != nil {
		return 0
	}

	return count
}

/* savind data infomration to db for access file in future
 */
func (f *FileserviceRepo) AddToDb(request model.FileServiceRequest, newFileName string) error {
	//save directory to db
	directoryFileList := model.FileService{
		CreatedBy: request.CreatedBy,
		UpdatedBy: request.CreatedBy,
		FileName:  newFileName,
		UuidDoc:   request.UuidDoc,
		DocType:   request.DocType,
	}

	if err := f.db.Create(&directoryFileList).Error; err != nil {
		return err
	}

	return nil
}

/* savind data infomration to db for access file in future
 */
func (f *FileserviceRepo) DeleteFromDb(fileName string) error {
	//delete from db with where from filename
	if err := f.db.Where(&model.FileService{FileName: fileName}).Delete(&model.FileService{}).Error; err != nil {
		return err
	}

	return nil
}

// get basepath file from .env
func (f *FileserviceRepo) GetFilePathFromEnv() (string, error) {
	basePath := viper.GetString("assets_path")
	if basePath == "" {
		return "", errors.New("assets_path not found in .env file")
	}

	return basePath, nil
}

// update and delete file use for V1 fileservice update
func (f *FileserviceRepo) DeleteAndUpdateV1(request model.FileServiceRequest, requestDeleted model.FileServiceRespont, form *multipart.Form) (model.FileServiceRespont, error) {
	//start to delete
	deletedRequest := model.FileServiceRequest{
		FileName: requestDeleted.FileName,
	}

	fileDeleted, err := f.DeleteFile(deletedRequest)
	if err != nil {
		return model.FileServiceRespont{}, err
	}

	log.Info(fileDeleted, "deleted file", fileDeleted.FileName)

	//start uploading new file
	fileUplouaded, err := f.SaveFile(request, form)
	if err != nil {
		return model.FileServiceRespont{}, err
	}

	return fileUplouaded, nil
}
