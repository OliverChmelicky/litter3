package middleware

import (
	"context"
	"firebase.google.com/go/auth"
	"github.com/labstack/echo"
	"net/http"
)

type MiddlewareService struct {
	Connection *auth.Client
}

func NewMiddlewareService(authConn *auth.Client) (*MiddlewareService, error) {
	return &MiddlewareService{Connection: authConn}, nil

}

func (s *MiddlewareService) AuthorizeUser(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Request().Header.Get("Authorization")
		if token == "" {
			return c.JSON(http.StatusUnauthorized, "error missing token")
		}

		firebaseToken, err := s.Connection.VerifyIDToken(context.Background(), token)
		if err != nil {
			return c.String(http.StatusUnauthorized, "Invalid authorization")
		}

		userId := firebaseToken.Claims["userId"]
		c.Set("userId", userId)

		return next(c)
	}
}

func (s *MiddlewareService) FillUserContext(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Request().Header.Get("Authorization")
		if token == "" {
			return next(c)
		}

		firebaseToken, err := s.Connection.VerifyIDToken(context.Background(), token)
		if err != nil {
			return c.String(http.StatusUnauthorized, "Invalid authorization")
		}

		userId := firebaseToken.Claims["userId"]
		c.Set("userId", userId)

		return next(c)
	}
}
