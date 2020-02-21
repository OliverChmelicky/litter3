package auth

import (
    "context"
    "firebase.google.com/go/auth"
    "fmt"
    "log"
    "net/http"
    "time"

    "github.com/dgrijalva/jwt-go"
    "github.com/labstack/echo"
    "github.com/labstack/echo/middleware"
    firebase "firebase.google.com/go"
    "google.golang.org/api/option"
)

func createAuthService() *AuthService{
    opt := option.WithCredentialsFile("../secrets/litter3-olo-gcp-firebase-adminsdk-6ar5p-9f1130c1cc.json")
    app, err := firebase.NewApp(context.Background(), nil, opt)
    if err != nil {
        log.Fatalf("error initializing app: %v", err)
    }

    authConn, err := app.Auth(context.Background())
    if err != nil {
        log.Fatalf("je to v rici")
    }


    tokenZNetu := ""

    token, err := authConn.VerifyIDToken(context.Background(), tokenZNetu)
    if err != nil {
        log.Fatalf("Neoveril som")
    }

    fmt.Println("Mam")
    fmt.Println(token.Claims)

}

func (s *AuthService) AuthorizeUser(c echo.Context) *auth.Token{
    authToken :=

    token, err := s.Connection.VerifyIDToken(context.Background(), c.)
    if err != nil {
        log.Fatalf("Neoveril som")
    }

    return token
}
