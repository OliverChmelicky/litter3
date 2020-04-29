package main

import (
	"context"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"fmt"
	"github.com/go-pg/pg/v9"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/olo/litter3/fileupload"
	middlewareService "github.com/olo/litter3/middleware"
	"github.com/olo/litter3/trash"
	"github.com/olo/litter3/user"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/api/option"
	"net/http"
)

func main() {
	viper.SetDefault("dbUsr", "goo")
	viper.SetDefault("dbPass", "goo")
	viper.SetDefault("dbName", "goo")
	viper.SetDefault("dbAddr", "localhost:5432")
	viper.SetDefault("firebaseCredentials", "../secrets/litter3-olo-gcp-firebase-adminsdk-6ar5p-9f1130c1cc.json")
	viper.SetDefault("gcpBucketName", "litter3-olo-gcp.appspot.com")

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	db := pg.Connect(&pg.Options{
		User:     viper.GetString("dbUsr"),
		Password: viper.GetString("dbPass"),
		Database: viper.GetString("dbName"),
		Addr:     viper.GetString("dbAddr"),
	})
	defer db.Close()
	_, err = db.Exec("SELECT 1")
	if err != nil {
		log.Error("PostgreSQL is down")
		return
	}
	db.AddQueryHook(middlewareService.DbMiddleware{})

	opt := option.WithCredentialsFile(viper.GetString("firebaseCredentials"))
	firebaseAuth, err := getFirebaseAuth(opt)
	if err != nil {
		log.Fatal(err)
	}

	fileuploadService := fileupload.CreateService(db, opt, viper.GetString("gcpBucketName"))

	e := echo.New()
	tokenMiddleware, err := middlewareService.NewMiddlewareService(firebaseAuth)
	if err != nil {
		log.Fatal(err)
	}
	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(tokenMiddleware.CorsHeadder)
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowMethods: []string{http.MethodOptions},
	}))

	userService := user.CreateService(db, firebaseAuth, fileuploadService)
	e.POST("/users/new", userService.CreateUser)
	e.GET("users/:id", userService.GetUser)
	e.GET("users/me", userService.GetCurrentUser, tokenMiddleware.AuthorizeUser)
	e.PUT("users/me", userService.UpdateUser, tokenMiddleware.AuthorizeUser)
	e.POST("/users/friend/add/email", userService.ApplyForFriendshipByEmail, tokenMiddleware.AuthorizeUser)
	e.POST("/users/friend/add/id", userService.ApplyForFriendshipById, tokenMiddleware.AuthorizeUser)

	e.POST("/societies/new", userService.CreateUser, tokenMiddleware.AuthorizeUser)
	e.PUT("/societies/update", userService.UpdateSociety, tokenMiddleware.AuthorizeUser)

	e.POST("membership", userService.ApplyForMembership, tokenMiddleware.AuthorizeUser)
	e.DELETE("membership/:societyId", userService.RemoveApplicationForMembership, tokenMiddleware.AuthorizeUser)

	trashService := trash.CreateService(db)
	e.GET("/trash/:id", trashService.GetTrashById)
	e.POST("/trash", trashService.CreateTrash)

	e.Logger.Fatal(e.Start(":1323"))
}

func getFirebaseAuth(opt option.ClientOption) (*auth.Client, error) {
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
