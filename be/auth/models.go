package auth

import (
	"firebase.google.com/go/auth"
	"github.com/dgrijalva/jwt-go"
)

type JwtCustomClaims struct {
	Admin bool `json:"admin"`
	jwt.StandardClaims
}

type AuthService struct {
	Connection *auth.Client
}
