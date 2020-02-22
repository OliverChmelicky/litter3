package middleware

import "firebase.google.com/go/auth"

type MiddlewareService struct {
	Connection *auth.Client
}
