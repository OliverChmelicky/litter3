package main

import (
	"context"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"github.com/go-pg/pg/v9"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	middlewareService "github.com/olo/litter3/middleware"
	"github.com/olo/litter3/trash"
	"github.com/olo/litter3/user"
	log "github.com/sirupsen/logrus"
	"google.golang.org/api/option"
)

func main() {
	db := pg.Connect(&pg.Options{
		User:     "goo",
		Password: "goo",
		Database: "goo",
		Addr:     "localhost:5432",
	})
	defer db.Close()
	_, err := db.Exec("SELECT 1")
	if err != nil {
		log.Error("PostgreSQL is down")
	}
	db.AddQueryHook(middlewareService.DbMiddleware{})

	firebaseAuth, err := getFirebaseAuth()
	if err != nil {
		log.Fatal(err)
	}

	e := echo.New()
	tokenMiddleware, err := middlewareService.NewMiddlewareService(firebaseAuth)
	if err != nil {
		log.Fatal(err)
	}
	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	userService := user.CreateService(db)
	e.POST("/users/new", userService.CreateUser)
	e.GET("users/:id", userService.GetUser, tokenMiddleware.FillUserContext)
	e.GET("users/me", userService.GetCurrentUser, tokenMiddleware.AuthorizeUser)
	e.PUT("users/me", userService.UpdateUser, tokenMiddleware.AuthorizeUser)

	e.POST("/societies/new", userService.CreateUser, tokenMiddleware.AuthorizeUser)

	e.POST("/societies/new", userService.CreateSociety, tokenMiddleware.AuthorizeUser)

	trashService := trash.CreateService(db)
	e.GET("/trash/:id", trashService.GetTrashById)
	e.POST("/trash", trashService.CreateTrash) //point makes troubles while inserting

	e.Logger.Fatal(e.Start(":1323"))
}

func getFirebaseAuth() (*auth.Client, error) {
	opt := option.WithCredentialsFile("secrets/litter3-olo-gcp-firebase-adminsdk-6ar5p-9f1130c1cc.json")
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Errorf("error initializing app: %s\n", err.Error())
		return &auth.Client{}, err
	}

	authConn, err := app.Auth(context.Background())
	if err != nil {
		return &auth.Client{}, err
	}

	return authConn, nil
}
