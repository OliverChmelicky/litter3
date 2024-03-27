package main

import (
	"context"
	"fmt"
	"os"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"github.com/go-pg/pg/v9"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/api/option"

	middlewareService "github.com/olo/litter3/services/shared"
	trash "github.com/olo/litter3/services/trash/pkg"
)

func main() {
	viper.SetDefault("DB_USR", os.Getenv("DB_USR"))
	viper.SetDefault("DB_PASS", os.Getenv("DB_PASS"))
	viper.SetDefault("DB_NAME", os.Getenv("DB_NAME"))
	viper.SetDefault("DB_ADDR", os.Getenv("DB_ADDR"))
	viper.SetDefault("FIREBASE_CREDENTIALS", os.Getenv("FIREBASE_CREDENTIALS"))
	viper.SetDefault("ADDRESS", os.Getenv("ADDRESS"))
	viper.AutomaticEnv()

	prod := os.Getenv("PROD")
	var production bool
	if len(prod) == 0 {
		production = false
	} else {
		production = true
	}

	var db *pg.DB
	if production {
		db = newDBAppEngine()
	} else {
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(".")
		err := viper.ReadInConfig()
		if err != nil {
			log.Errorf("Fatal error config file: %s \n", err)
		}

		db = pg.Connect(&pg.Options{
			User:     viper.GetString("DB_USR"),
			Password: viper.GetString("DB_PASS"),
			Database: viper.GetString("DB_NAME"),
			Addr:     viper.GetString("DB_ADDR"),
		})
	}
	_, err := db.Exec("SELECT 1")
	if err != nil {
		log.Errorf("PostgreSQL is down: %s \n", err.Error())
		return
	}
	defer db.Close()

	opt := option.WithCredentialsFile(viper.GetString("FIREBASE_CREDENTIALS"))
	firebaseAuth, err := getFirebaseAuth(opt)
	if err != nil {
		log.Fatal(err)
	}

	e := echo.New()
	tokenMiddleware, err := middlewareService.NewMiddlewareService(firebaseAuth)
	if err != nil {
		log.Fatal(err)
	}
	// Middleware
	//e.Use(middleware.Logger())
	//e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		Skipper:      middleware.DefaultSkipper,
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
	}))

	trashService := trash.CreateService(db)
	e.GET("/trash/:id", trashService.GetTrashById)
	e.GET("/trash", trashService.GetTrashByIds)
	e.GET("/trash/range", trashService.GetTrashInRange)
	e.POST("/trash/new", trashService.CreateTrash, tokenMiddleware.FillUserContext)
	e.PUT("/trash/update", trashService.UpdateTrash, tokenMiddleware.AuthorizeUser)
	e.DELETE("/trash/delete/:trashId", trashService.DeleteTrash, tokenMiddleware.AuthorizeUser)

	e.POST("/trash/comment", trashService.CreateTrashComment, tokenMiddleware.AuthorizeUser)
	e.DELETE("/trash/comment/:commentId", trashService.DeleteTrashComment, tokenMiddleware.AuthorizeUser)

	if production {
		listenOn := fmt.Sprintf(":%s", os.Getenv("PORT"))
		e.Logger.Fatal(e.Start(listenOn))
	} else {
		e.Logger.Fatal(e.Start(viper.GetString("ADDRESS")))
	}

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

func newDBAppEngine() *pg.DB {
	return pg.Connect(&pg.Options{
		User:     os.Getenv("DB_USR"),
		Password: os.Getenv("DB_PASS"),
		Addr:     viper.GetString("DB_ADDR"),
		Database: os.Getenv("DB_NAME"),
		Network:  "unix",
	})
}
