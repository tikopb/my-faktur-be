package user

import (
	"bemyfaktur/internal/model"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/sirupsen/logrus"
)

type Claims struct {
	jwt.StandardClaims
}

type RefreshClaims struct {
	jwt.StandardClaims
}

type LogOutSession struct {
	AccessToken  string
	RefreshToken string
	Created      time.Time
}

var LogOutSessionArrrayMemory []LogOutSession
var RoutineRunning bool = false

// CreateUserSession implements Repository.
func (ur *userRepo) CreateUserSession(userID string) (model.UserSession, error) {
	//generate access token first and saving into claims
	accessToken, err := ur.generateAccessToken(userID)
	if err != nil {
		return model.UserSession{}, err
	}

	//generate refresh token first and saving into refresh claims
	refreshToken, err := ur.generateRefreshToken(userID)
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

	//cek logout token session first!
	if ur.checkLogOutSession(data.AccessToken) {
		return "", errors.New("access token expired/invalid")
	}

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

func (ur *userRepo) CheckRefreshToken(RefreshToken string) (userID string, err error) {

	//cek logout token session first!
	if ur.checkLogOutSession(RefreshToken) {
		return "", errors.New("access token expired/invalid")
	}

	refreshToken, err := jwt.ParseWithClaims(RefreshToken, &RefreshClaims{}, func(token *jwt.Token) (interface{}, error) {
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

func (ur *userRepo) generateRefreshToken(userID string) (string, error) {
	refreshTokenExp := time.Now().Add(ur.refreshTimeout).Unix()
	refreshClaims := RefreshClaims{
		jwt.StandardClaims{
			ExpiresAt: refreshTokenExp,
			Subject:   userID,
		},
	}

	refreshJwt := jwt.NewWithClaims(jwt.GetSigningMethod("RS256"), refreshClaims)

	return refreshJwt.SignedString(ur.signKey)
}

func (ur *userRepo) LogOut(data model.UserSession) {
	logOutData := LogOutSession{
		AccessToken:  data.AccessToken,
		RefreshToken: data.RefreshToken,
		Created:      time.Now(),
	}

	// Append the LogOutSession data to the LogOutSessionArrrayMemory slice
	LogOutSessionArrrayMemory = append(LogOutSessionArrrayMemory, logOutData)
}

func (ur *userRepo) checkLogOutSession(token string) bool {
	for _, session := range LogOutSessionArrrayMemory {
		if session.AccessToken == token || session.RefreshToken == token {
			return true
		}
	}

	// The UserSession doesn't exist in LogOutSessionArrrayMemory
	if !RoutineRunning {
		go ur.CleanupSessions()
		RoutineRunning = true
	}
	return false
}

// func to clean the loutseesion after 2 hour of running
func (ur *userRepo) CleanupSessions() {
	var sessions []LogOutSession

	for {
		// Lock sessions for writing
		var updatedSessions []LogOutSession

		// Check each session and remove sessions older than 2 hours
		for _, session := range sessions {
			if time.Since(session.Created) <= 2*time.Hour {
				updatedSessions = append(updatedSessions, session)
			}
		}

		// Update the sessions
		sessions = updatedSessions

		// Unlock sessions
		logrus.Info("Cleanup completed. Current sessions")
		for _, session := range sessions {
			fmt.Println(session)
		}

		// Sleep for 1 hour before the next cleanup
		time.Sleep(time.Hour)
	}
}

func (ur *userRepo) GetuserIdFromClaims(accesstoken string) (string, error) {

	accessToken, err := jwt.ParseWithClaims(accesstoken, &Claims{}, func(token *jwt.Token) (interface{}, error) {
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
