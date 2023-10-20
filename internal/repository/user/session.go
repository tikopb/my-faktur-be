package user

import (
	"bemyfaktur/internal/model"
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
)

type Claims struct {
	jwt.StandardClaims
}

type RefreshClaims struct {
	jwt.StandardClaims
}

// CreateUserSession implements Repository.
func (ur *userRepo) CreateUserSession(userID string) (model.UserSession, error) {
	//generate access token first and saving into claims
	accessToken, err := ur.generateAccessToken(userID)
	if err != nil {
		return model.UserSession{}, err
	}

	//generate refresh token first and saving into refresh claims
	refreshToken, err := ur.GenerateRefreshToken(userID)
	if err != nil {
		return model.UserSession{}, err
	}

	return model.UserSession{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// CheckSession implements Repository.
func (ur *userRepo) CheckSession(data model.UserSession) (userID string, err error) {
	accessToken, err := jwt.ParseWithClaims(data.AccessToken, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return &ur.signKey.PublicKey, nil
	})

	if err != nil {
		return "", errors.New("access token expired/invalid")
	}

	accessTokenClaims, ok := accessToken.Claims.(*Claims)
	if !ok {
		return "", errors.New("unauthorized")
	}

	if accessToken.Valid {
		return accessTokenClaims.Subject, nil
	}

	return "", errors.New("unauthorized")
}

func (ur *userRepo) CheckRefreshToken(data model.UserSession) (userID string, err error) {
	refreshToken, err := jwt.ParseWithClaims(data.RefreshToken, &RefreshClaims{}, func(token *jwt.Token) (interface{}, error) {
		return &ur.signKey.PublicKey, nil
	})
	if err != nil {
		return "", errors.New("refresh token expired/invalid")
	}

	refreshTokenClaims, ok := refreshToken.Claims.(*RefreshClaims)
	if !ok {
		return "", errors.New("unauthorized")
	}

	if refreshToken.Valid {
		return refreshTokenClaims.Subject, nil
	}

	return "", errors.New("unauthorized")
}

func (ur *userRepo) generateAccessToken(userID string) (string, error) {
	accessTokenExp := time.Now().Add(ur.accessExp).Unix()
	accessClaims := Claims{
		jwt.StandardClaims{
			ExpiresAt: accessTokenExp,
			Subject:   userID,
		},
	}

	accessJwt := jwt.NewWithClaims(jwt.GetSigningMethod("RS256"), accessClaims)

	return accessJwt.SignedString(ur.signKey)
}

func (ur *userRepo) GenerateRefreshToken(userID string) (string, error) {
	refreshTokenExp := time.Now().Add(48 * time.Hour).Unix()
	refreshClaims := RefreshClaims{
		jwt.StandardClaims{
			ExpiresAt: refreshTokenExp,
			Subject:   userID,
		},
	}

	refreshJwt := jwt.NewWithClaims(jwt.GetSigningMethod("RS256"), refreshClaims)

	return refreshJwt.SignedString(ur.signKey)
}
