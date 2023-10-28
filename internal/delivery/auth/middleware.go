package auth

import (
	"bemyfaktur/internal/model"
	"bemyfaktur/internal/model/constant"
	"bemyfaktur/internal/usecase/auth"
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
)

type authMiddleware struct {
	authUsecase auth.Usecase
}

type MidlewareInterface interface {
	CheckAuth(next echo.HandlerFunc) echo.HandlerFunc
	GetSessionData(r *http.Request) (model.UserSession, error)
	GetuserId(r *http.Request) (string, error)
	GetValueParamHeader(r *http.Request, param string) (string, error)
}

func GetAuthMiddleware(authusecase auth.Usecase) MidlewareInterface {
	return &authMiddleware{
		authUsecase: authusecase,
	}
}

func (am *authMiddleware) CheckAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		sessionData, err := am.GetSessionData(c.Request())
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
