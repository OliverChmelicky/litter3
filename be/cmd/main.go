package main

import (
	"context"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"fmt"
	"github.com/go-pg/pg/v9"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/olo/litter3/event"
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
		AllowMethods: []string{http.MethodOptions, http.MethodDelete, http.MethodPut},
	}))

	userService := user.CreateService(db, firebaseAuth, fileuploadService)
	e.POST("/users/new", userService.CreateUser)
	e.GET("/users/:id", userService.GetUser)
	e.GET("/users/email/:email", userService.GetUserByEmail)
	e.GET("/users/me", userService.GetCurrentUser, tokenMiddleware.AuthorizeUser)
	e.PUT("/users/update", userService.UpdateUser, tokenMiddleware.AuthorizeUser)
	e.GET("/users/details", userService.GetUsers)

	e.GET("/users/societies", userService.GetMySocieties, tokenMiddleware.AuthorizeUser)
	e.GET("/users/societies/editable", userService.GetEditableSocieties, tokenMiddleware.AuthorizeUser)

	e.POST("/users/friend/add/id", userService.ApplyForFriendshipById, tokenMiddleware.AuthorizeUser)
	e.POST("/users/friend/add/email", userService.ApplyForFriendshipByEmail, tokenMiddleware.AuthorizeUser)
	e.GET("/users/friends", userService.GetMyFriends, tokenMiddleware.AuthorizeUser)
	e.GET("/users/friend/requests", userService.GetMyFriendRequests, tokenMiddleware.AuthorizeUser)
	e.DELETE("/users/friend/remove/:notWanted", userService.RemoveFriend, tokenMiddleware.AuthorizeUser)
	e.POST("/users/friend/accept/:wantedUser", userService.AcceptFriendship, tokenMiddleware.AuthorizeUser)
	e.DELETE("/users/friend/deny/:notWanted", userService.RemoveApplicationForFriendship, tokenMiddleware.AuthorizeUser)

	e.POST("/societies/new", userService.CreateSociety, tokenMiddleware.AuthorizeUser)
	e.GET("/societies", userService.GetSocietiesWithPaging)
	e.GET("/societies/:id", userService.GetSociety)
	e.PUT("/societies/update", userService.UpdateSociety, tokenMiddleware.AuthorizeUser)
	e.GET("/societies/admins/:societyId", userService.GetSocietyAdmins)
	e.GET("/societies/members/:societyId", userService.GetSocietyMembers)
	e.GET("/societies/requests/:societyId", userService.GetSocietyRequests)
	e.PUT("/societies/change-permission", userService.ChangeMemberRights, tokenMiddleware.AuthorizeUser)
	e.DELETE("/societies/:societyId/:removingId", userService.RemoveMember, tokenMiddleware.AuthorizeUser)

	e.POST("/membership", userService.ApplyForMembership, tokenMiddleware.AuthorizeUser)
	e.DELETE("/membership/:societyId", userService.RemoveApplicationForMembership, tokenMiddleware.AuthorizeUser)

	eventService := event.CreateService(db)
	e.POST("/events", eventService.CreateEvent, tokenMiddleware.AuthorizeUser)
	e.GET("/events", eventService.GetEventsWithPaging)
	e.GET("/events/societies/:societyId", eventService.GetSocietyEvents)
	e.GET("/events/:eventId", eventService.GetEvent)
	e.POST("/events/attend", eventService.AttendEvent, tokenMiddleware.AuthorizeUser)
	e.DELETE("/events/not-attend", eventService.CannotAttendEvent, tokenMiddleware.AuthorizeUser)

	trashService := trash.CreateService(db)
	e.GET("/trash/:id", trashService.GetTrashById)
	e.GET("/trash/range", trashService.GetTrashInRange)
	e.POST("/trash/new", trashService.CreateTrash, tokenMiddleware.FillUserContext)
	e.PUT("/trash/update", trashService.UpdateTrash, tokenMiddleware.AuthorizeUser)
	e.DELETE("/trash/delete/:trashId", trashService.DeleteTrash)

	e.POST("/fileupload/trash/:trashId", fileuploadService.UploadTrashImages)

	e.GET("/fileupload/societies/:image", fileuploadService.GetSocietyImage)
	e.POST("/fileupload/societies/:societyId", fileuploadService.UploadSocietyImage, tokenMiddleware.AuthorizeUser)

	e.Logger.Fatal(e.Start(":8081"))
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
