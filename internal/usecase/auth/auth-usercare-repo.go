package auth

import "bemyfaktur/internal/model"

type Usecase interface {
	RegisterUser(request model.RegisterRequest) (model.User, error)
	Login(request model.LoginRequest) (model.UserSession, error)
	CheckSession(sessionData model.UserSession) (userID string, err error)
	RefreshToken(userSession model.UserSession) (model.UserSession, error)
}
