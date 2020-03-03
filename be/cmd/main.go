package main

import (
	"github.com/go-pg/pg/v9"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	middlewareService "github.com/olo/litter3/middleware"
	log "github.com/sirupsen/logrus"
)

type user struct {
	Id        string
	FirstName string
	LastName  string
	Email     string
}

func main() {
	db := pg.Connect(&pg.Options{
		User:     "goo",
		Password: "goo",
		Database: "goo",
		Addr:     "localhost:5432",
	})
	defer db.Close()
	db.AddQueryHook(middlewareService.DbMiddleware{})

	usr := &user{Id: "9", FirstName: "Oliver"}
	//res, err := db.Model(&usr).Insert(&usr)
	err := db.Insert(usr)
	if err != nil {
		log.Error(err)
		return
	}

	//log.Info(res)

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
