package auth

import (
	"bemyfaktur/internal/model/constant"
	"bemyfaktur/internal/usecase/auth"
	"context"

	"github.com/labstack/echo/v4"
)

type authMiddleware struct {
	authUsecase auth.Usecase
}

func GetAuthMiddleware(authusecase auth.Usecase) *authMiddleware {
	return &authMiddleware{
		authUsecase: authusecase,
	}
}

func (am *authMiddleware) CheckAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		sessionData, err := GetSessionData(c.Request())
		if err != nil {
			return &echo.HTTPError{
				Code:     401,
				Message:  err.Error(),
				Internal: err,
			}
		}

		userID, err := am.authUsecase.CheckSession(sessionData)
		if err != nil {
			return &echo.HTTPError{
				Code:     401,
				Message:  err.Error(),
				Internal: err,
			}
		}

		authContext := context.WithValue(c.Request().Context(), constant.AuthContextKey, userID)
		c.SetRequest(c.Request().WithContext(authContext))

		if err := next(c); err != nil {
			return err
		}

		return nil
	}
}
