package user

import "bemyfaktur/internal/model"

type Repository interface {
	RegisterUser(userData model.User) (model.User, error)
	CheckRegistered(username string) (bool, error)
	GenerateUserHash(password string) (hash string, err error)
	VerifyLogin(username, password string, userData model.User) (bool, error)
	GetUserData(username string) (model.User, error)
	CreateUserSession(userID string) (model.UserSession, error)
	CheckSession(data model.UserSession) (string, error)
	CheckRefreshToken(RefreshToken string) (string, error)
	LogOut(data model.UserSession)
}
