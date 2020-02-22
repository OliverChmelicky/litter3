package main

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	middlewareService "github.com/olo/litter3/middleware"
	log "github.com/sirupsen/logrus"
)

func main() {
	e := echo.New()
	tokenMiddleware, err := middlewareService.NewMiddlewareService()
	if err != nil {
		log.Fatal(err)
	}

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", tokenMiddleware.AllOk, tokenMiddleware.AuthorizeUser)

	e.Logger.Fatal(e.Start(":1323"))
}
