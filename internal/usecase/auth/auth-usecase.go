package auth

import (
	"bemyfaktur/internal/model"
	"bemyfaktur/internal/repository/user"

	"errors"

	"github.com/google/uuid"
)

type authStruct struct {
	userRepo user.Repository
}

func GetUsecase(userRepo user.Repository) Usecase {
	return &authStruct{
		userRepo: userRepo,
	}
}

// CheckSession implements Usecase.
func (au *authStruct) CheckSession(sessionData model.UserSession) (userID string, err error) {
	userID, err = au.userRepo.CheckSession(sessionData)
	if err != nil {
		return "", err
	}
	return userID, nil
}

// Login implements Usecase.
func (au *authStruct) Login(request model.LoginRequest) (model.UserSession, error) {
	userData, err := au.userRepo.GetUserData(request.Username)
	if err != nil {
		return model.UserSession{}, err
	}

	verified, err := au.userRepo.VerifyLogin(request.Username, request.Password, userData)
	if err != nil {
		return model.UserSession{}, err
	}

	if !verified {
		return model.UserSession{}, errors.New("can't verifit user login")
	}

	userSession, err := au.userRepo.CreateUserSession(userData.ID)
	if err != nil {
		return model.UserSession{}, err
	}

	return userSession, nil
}

// RegisterUser implements Usecase.
func (au *authStruct) RegisterUser(request model.RegisterRequest) (model.User, error) {
	userRegistered, err := au.userRepo.CheckRegistered(request.Username)
	if err != nil {
		return model.User{}, err
	}
	if userRegistered {
		return model.User{}, errors.New("user already registered")
	}

	userHash, err := au.userRepo.GenerateUserHash(request.Password)
	if err != nil {
		return model.User{}, err
	}

	userData, err := au.userRepo.RegisterUser(model.User{
		ID:       uuid.New().String(),
		Username: request.Username,
		Hash:     userHash,
	})
	if err != nil {
		return model.User{}, err
	}

	return userData, nil
}
