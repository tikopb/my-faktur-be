package auth

import (
	"bemyfaktur/internal/model"
	"errors"
	"net/http"
	"strings"
)

func (am *authMiddleware) GetSessionData(r *http.Request) (model.UserSession, error) {
	authString := r.Header.Get("Authorization")
	splitString := strings.Split(authString, " ")
	if len(splitString) != 2 {
		return model.UserSession{}, errors.New("unauthorized")
	}
	accessString := splitString[1]

	return model.UserSession{
		AccessToken: accessString,
	}, nil
}

func (am *authMiddleware) GetuserId(r *http.Request) (string, error) {
	sessionData, err := am.GetSessionData(r)
	if err != nil {
		return "", err
	}

	userID, err := am.authUsecase.CheckSession(sessionData)
	if err != nil {
		return "", err
	}

	return userID, nil
}
