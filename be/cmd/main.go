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
	"os"
)

func main() {
	viper.SetDefault("DB_USR", os.Getenv("DB_USR"))
	viper.SetDefault("DB_PASS", os.Getenv("DB_PASS"))
	viper.SetDefault("DB_NAME", os.Getenv("DB_NAME"))
	viper.SetDefault("DB_ADDR", os.Getenv("DB_ADDR"))
	viper.SetDefault("FIREBASE_CREDENTIALS", os.Getenv("FIREBASE_CREDENTIALS"))
	viper.SetDefault("GCP_BUCKET_NAME", os.Getenv("GCP_BUCKET_NAME"))
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

	fileuploadService := fileupload.CreateService(db, opt, viper.GetString("GCP_BUCKET_NAME"))

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

	userService := user.CreateService(db, firebaseAuth)
	e.POST("/users/new", userService.CreateUser)
	e.GET("/users/:id", userService.GetUser)
	e.GET("/users/email/:email", userService.GetUserByEmail)
	e.GET("/users/me", userService.GetCurrentUser, tokenMiddleware.AuthorizeUser)
	e.PUT("/users/update", userService.UpdateUser, tokenMiddleware.AuthorizeUser)
	e.GET("/users/details", userService.GetUsers)
	e.DELETE("/users/delete", userService.DeleteUser)

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
	e.GET("/societies/more", userService.GetSocietiesByIds)
	e.GET("/societies/:id", userService.GetSociety)
	e.PUT("/societies/update", userService.UpdateSociety, tokenMiddleware.AuthorizeUser)
	e.GET("/societies/admins/:societyId", userService.GetSocietyAdmins)
	e.GET("/societies/members/:societyId", userService.GetSocietyMembers)
	e.GET("/societies/requests/:societyId", userService.GetSocietyRequests, tokenMiddleware.AuthorizeUser)
	e.PUT("/societies/change/permission", userService.ChangeMemberRights, tokenMiddleware.AuthorizeUser)
	e.DELETE("/societies/:societyId/:removingId", userService.RemoveMember, tokenMiddleware.AuthorizeUser)
	e.DELETE("/societies/delete/:societyId", userService.DeleteSociety, tokenMiddleware.AuthorizeUser)

	e.POST("/membership", userService.ApplyForMembership, tokenMiddleware.AuthorizeUser)
	e.POST("/membership/accept/:societyId/:userId", userService.AcceptApplicant, tokenMiddleware.AuthorizeUser)
	e.DELETE("/membership/deny/:societyId/:userId", userService.DismissApplicant, tokenMiddleware.AuthorizeUser)
	e.DELETE("/membership/:societyId", userService.RemoveApplicationForMembership, tokenMiddleware.AuthorizeUser)

	eventService := event.CreateService(db)
	e.POST("/events", eventService.CreateEvent, tokenMiddleware.AuthorizeUser)
	e.GET("/events", eventService.GetEventsWithPaging)
	e.GET("/events/societies/:societyId", eventService.GetSocietyEvents)
	e.GET("/events/:eventId", eventService.GetEvent)
	e.POST("/events/attend", eventService.AttendEvent, tokenMiddleware.AuthorizeUser)
	e.DELETE("/events/not-attend", eventService.CannotAttendEvent, tokenMiddleware.AuthorizeUser)
	e.PUT("/events/update", eventService.UpdateEvent, tokenMiddleware.AuthorizeUser)
	e.PUT("/events/members/update", eventService.EditEventRights, tokenMiddleware.AuthorizeUser)
	e.DELETE("/events/delete", eventService.DeleteEvent, tokenMiddleware.AuthorizeUser)

	trashService := trash.CreateService(db)
	e.GET("/trash/:id", trashService.GetTrashById)
	e.GET("/trash", trashService.GetTrashByIds)
	e.GET("/trash/range", trashService.GetTrashInRange)
	e.POST("/trash/new", trashService.CreateTrash, tokenMiddleware.FillUserContext)
	e.PUT("/trash/update", trashService.UpdateTrash, tokenMiddleware.AuthorizeUser)
	e.DELETE("/trash/delete/:trashId", trashService.DeleteTrash, tokenMiddleware.AuthorizeUser)

	e.POST("/trash/comment", trashService.CreateTrashComment, tokenMiddleware.AuthorizeUser)
	e.DELETE("/trash/comment/:commentId", trashService.DeleteTrashComment, tokenMiddleware.AuthorizeUser)

	e.GET("/collections/personal", trashService.GetCollectionIdsOfUser, tokenMiddleware.AuthorizeUser)
	e.POST("/collections/organized", eventService.CreateCollectionsOrganized, tokenMiddleware.AuthorizeUser)
	e.POST("/collections/add-picker", trashService.AddPickerToCollection, tokenMiddleware.AuthorizeUser)
	e.POST("/collections/random", trashService.CreateCollection, tokenMiddleware.AuthorizeUser)
	e.GET("/collections/:collectionId", trashService.GetCollection)
	e.PUT("/collections/update/col-organized", eventService.UpdateCollectionOrganized, tokenMiddleware.AuthorizeUser)
	e.PUT("/collections/update/col-random", trashService.UpdateCollectionRandom, tokenMiddleware.AuthorizeUser)
	e.DELETE("/collections/delete/:collectionId", trashService.DeleteCollectionFromUser, tokenMiddleware.AuthorizeUser)
	e.DELETE("/collections/delete/organized", eventService.DeleteCollectionOrganized, tokenMiddleware.AuthorizeUser) //query params

	e.POST("/fileupload/societies/:societyId", fileuploadService.UploadSocietyImage, tokenMiddleware.AuthorizeUser)
	e.GET("/fileupload/societies/load/:image", fileuploadService.GetSocietyImage)
	e.POST("/fileupload/trash/:trashId", fileuploadService.UploadTrashImages)
	e.GET("/fileupload/trash/load/:image", fileuploadService.GetTrashImage)
	e.DELETE("/fileupload/trash/delete/:trashId/:image", fileuploadService.DeleteTrashImage)
	e.POST("/fileupload/collections/:collectionId", fileuploadService.UploadCollectionImages)
	e.GET("/fileupload/collections/load/:image", fileuploadService.GetCollectionImages)
	e.DELETE("/fileupload/collections/delete/:collectionId", fileuploadService.DeleteCollectionImages, tokenMiddleware.AuthorizeUser)

	//e.POST("/fileupload/societies/:societyId", fileuploadService.UploadUserImage, tokenMiddleware.AuthorizeUser)
	//e.GET("/fileupload/societies/load/:image", fileuploadService.GetUserImage)
	//e.DELETE("/fileupload/collections/delete/:collectionId", fileuploadService.DeleteSocietyImage, tokenMiddleware.AuthorizeUser)
	//e.DELETE("/fileupload/collections/delete/:collectionId", fileuploadService.DeleteUserImage, tokenMiddleware.AuthorizeUser)

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
