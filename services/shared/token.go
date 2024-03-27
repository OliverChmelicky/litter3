package shared

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"firebase.google.com/go/auth"
	"github.com/labstack/echo"
)

type MiddlewareService struct {
	Connection *auth.Client
}

func NewMiddlewareService(authConn *auth.Client) (*MiddlewareService, error) {
	return &MiddlewareService{Connection: authConn}, nil

}

func (s *MiddlewareService) AuthorizeUser(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return c.JSON(http.StatusUnauthorized, WrapError(ErrNoToken, fmt.Errorf("No token was provided ")))
		}

		str := strings.Fields(authHeader)
		var token string
		if len(str) > 0 {
			token = str[len(str)-1]
		} else {
			return c.JSON(http.StatusUnauthorized, WrapError(ErrNoToken, fmt.Errorf("Bad Authorization header ")))
		}

		firebaseToken, err := s.Connection.VerifyIDToken(context.Background(), token)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, WrapError(ErrUnauthenticated, err))
		}

		userId, ok := firebaseToken.Claims["userId"].(string)
		if !ok {
			return c.JSON(http.StatusUnauthorized, WrapError(ErrUnauthenticated, err))
		}
		c.Set("userId", userId)

		return next(c)
	}
}

func (s *MiddlewareService) FillUserContext(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "Bearer null" {
			c.Set("userId", "")
			return next(c)
		}

		str := strings.Fields(authHeader)
		var token string
		if len(str) > 0 {
			token = str[len(str)-1]
		} else {
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

func (s *MiddlewareService) CorsHeadder(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set(echo.HeaderAccessControlAllowOrigin, "*")
		return next(c)
	}
}
