package middleware

import (
	"context"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"fmt"
	"github.com/labstack/echo"
	log "github.com/sirupsen/logrus"
	"google.golang.org/api/option"
	"net/http"
)

type MiddlewareService struct {
	Connection *auth.Client
}

func NewMiddlewareService() (*MiddlewareService, error) {
	opt := option.WithCredentialsFile("secrets/litter3-olo-gcp-firebase-adminsdk-6ar5p-9f1130c1cc.json")
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Error("error initializing app: %v", err)
		return nil, err
	}

	authConn, err := app.Auth(context.Background())
	if err != nil {
		log.Fatalf("je to v rici")
	}

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
			log.Error(err)
			return c.String(http.StatusUnauthorized, "Invalid authorization")
		}

		userId := firebaseToken.Claims["userId"]
		fmt.Println(userId)
		//c.Set("user_id", )

		return next(c)
	}
}

func (s *MiddlewareService) AllOk(c echo.Context) error {
	fmt.Println("Dnu")
	return c.String(http.StatusOK, "Ahoj")
}
