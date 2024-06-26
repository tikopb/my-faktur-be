package user

import (
	"bemyfaktur/internal/model"
	"bemyfaktur/internal/usecase/organization"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rsa"
	"errors"
	"time"

	"gorm.io/gorm"
)

type userRepo struct {
	db             *gorm.DB
	gcm            cipher.AEAD
	time           uint32
	memory         uint32
	parallelism    uint8
	keyLen         uint32
	secret         string
	signKey        *rsa.PrivateKey
	accessExp      time.Duration
	refreshTimeout time.Duration
	OrgUsecase     organization.Usecase
}

func GetRepository(
	db *gorm.DB,
	secret string,
	time, memory,
	keyLen uint32,
	parallelism uint8,
	signKey *rsa.PrivateKey,
	accessExp time.Duration,
	refreshTimeout time.Duration,
	OrgUsecase organization.Usecase,
) (Repository, error) {
	block, err := aes.NewCipher([]byte(secret))
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	return &userRepo{
		db:             db,
		gcm:            gcm,
		time:           time,
		memory:         memory,
		parallelism:    parallelism,
		keyLen:         keyLen,
		secret:         secret,
		signKey:        signKey,
		accessExp:      accessExp,
		refreshTimeout: refreshTimeout,
	}, nil
}

// RegisterUser implements Repository.
func (ur *userRepo) RegisterUser(userData model.User) (model.User, error) {
	if err := ur.db.Create(&userData).Error; err != nil {
		return model.User{}, err
	}
	return userData, nil
}

// CheckRegistered implements Repository.
var userData model.User

func (ur *userRepo) CheckRegistered(username string) (bool, error) {
	if err := ur.db.Where(model.User{Username: username}).First(&userData).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		} else {
			return false, err
		}
	}

	return userData.ID != "", nil
}

// GetUserData implements Repository.
func (ur *userRepo) GetUserData(username string) (model.User, error) {
	var userData model.User

	if err := ur.db.Where(model.User{Username: username, IsActive: true}).First(&userData).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return userData, errors.New("username not exist")
		}
		return userData, err
	}

	return userData, nil
}

// GetUserData By Id implements Repository.
func (ur *userRepo) GetUserDatById(id string) (model.User, error) {
	var userData model.User

	if err := ur.db.Where(model.User{ID: id, IsActive: true}).First(&userData).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return userData, errors.New("username not exist")
		}
		return userData, err
	}

	return userData, nil
}

// VerifyLogin implements Repository.
func (ur *userRepo) VerifyLogin(username string, password string, userData model.User) (bool, error) {
	if username != userData.Username {
		return false, nil
	}

	verified, err := ur.comparePassword(password, userData.Hash)
	if err != nil {
		return false, err
	}

	return verified, nil
}
