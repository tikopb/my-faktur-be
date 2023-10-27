package auth

import (
	"bemyfaktur/internal/model"
	"bemyfaktur/internal/repository/user"
	"strings"

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
func (au *authStruct) Login(request model.LoginRequest) (model.UserSessionRespond, error) {
	userSessionRespond := model.UserSessionRespond{}

	userData, err := au.userRepo.GetUserData(request.Username)
	if err != nil {
		return userSessionRespond, err
	}

	verified, err := au.userRepo.VerifyLogin(request.Username, request.Password, userData)
	if err != nil {
		return userSessionRespond, err
	}

	if !verified {
		return userSessionRespond, errors.New("can't verifit user login")
	}

	userSession, err := au.userRepo.CreateUserSession(userData.ID)
	if err != nil {
		return userSessionRespond, err
	}

	userSessionRespond = model.UserSessionRespond{
		UserSession:     userSession,
		UserInformation: userData,
	}
	return userSessionRespond, nil
}

// Refresh Token implements Usecase.
func (au *authStruct) RefreshToken(refreshToken string) (model.UserSession, error) {
	userSession := model.UserSession{}

	userID, err := au.userRepo.CheckRefreshToken(refreshToken)
	if err != nil {
		return model.UserSession{}, err
	}

	//generate of token access and refresh token with refres token as variabel
	userSession, err = au.userRepo.CreateUserSession(userID)
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

	if strings.Contains(request.Username, " ") {
		return model.User{}, errors.New("username can't having space")
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
		FullName: request.FullName,
		IsActive: true,
	})
	if err != nil {
		return model.User{}, err
	}

	return userData, nil
}

func (au *authStruct) LogOutUser(request model.UserSession) {
	au.userRepo.LogOut(request)
}
