package main

import (
	"github.com/go-pg/pg/v9"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	middlewareService "github.com/olo/litter3/middleware"
	log "github.com/sirupsen/logrus"
)

type user struct {
	id        string
	firstName string
	lastName  string
	email     string
}

func main() {
	db := pg.Connect(&pg.Options{
		User:     "goo",
		Password: "goo",
		Database: "goo",
	})
	defer db.Close()
	db.AddQueryHook(dbMiddleware{})

	usr := user{id: "9"}
	res, err := db.Model(&usr).Insert(&usr)
	if err != nil {
		log.Error(err)
		return
	}

	log.Info(res)

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
