package main

import (
	firebase "firebase.google.com/go"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/olo/litter3/auth"
	"google.golang.org/api/option"
)

func main() {
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Unauthenticated route
	e.GET("/", allOk)

	// Configure middleware with the custom claims type
	config := middleware.JWTConfig{
		Claims:     &auth.JwtCustomClaims{},
		SigningKey: []byte("secret"),
	}
	e.Use(auth.AuthorizeUser)

	e.Logger.Fatal(e.Start(":1323"))
}
