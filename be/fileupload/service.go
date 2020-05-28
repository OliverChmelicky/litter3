package fileupload

import (
	"cloud.google.com/go/storage"
	"context"
	firebase "firebase.google.com/go"
	"fmt"
	"github.com/go-pg/pg/v9"
	"github.com/labstack/echo"
	custom_errors "github.com/olo/litter3/custom-errors"
	"github.com/olo/litter3/models"
	log "github.com/sirupsen/logrus"
	"google.golang.org/api/option"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type FileuploadService struct {
	db *pg.DB
	bh *storage.BucketHandle
}

func CreateService(db *pg.DB, opt option.ClientOption, bucketName string) *FileuploadService {
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("Fileupload service error initializing firebase app: %s\n", err.Error())
		panic(err.Error())
	}
	st, err := app.Storage(context.Background())
	if err != nil {
		log.Fatalf("Fileupload service initializing storage: %s\n", err.Error())
		panic(err.Error())
	}
	bh, err := st.Bucket(bucketName)
	if err != nil {
		log.Fatalf("Fileupload service error getting bucket handler: %s\n", err.Error())
		panic(err.Error())
	}
	if _, err = bh.Attrs(context.Background()); err != nil {
		log.Fatalf("Bucket does not exist: %s\n", err.Error())
		panic(err.Error())
	}

	return &FileuploadService{
		db: db,
		bh: bh,
	}
}

func (s *FileuploadService) UploadUserImage(c echo.Context) error {
	userId := c.Get("userId").(string)

	user := new(models.User)
	err := s.db.Model(user).Where("id = ?", userId).Select()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, custom_errors.WrapError(custom_errors.ErrGetUserById, err))
	}
	if user.Avatar != "" {
		_ = s.DeleteImage(user.Avatar)
	}

	objectName, err := s.UploadImage(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, custom_errors.WrapError(custom_errors.ErrUploadImage, err))
	}

	user.Avatar = objectName
	err = s.db.Update(user)
	if err != nil {
		_ = s.DeleteImage(objectName)
		return c.JSON(http.StatusInternalServerError, custom_errors.WrapError(custom_errors.ErrUpdateUser, err))
	}

	return c.NoContent(http.StatusCreated)
}

func (s *FileuploadService) UploadSocietyImage(c echo.Context) error {
	societyId := c.Param("societyId")

	objectName, err := s.UploadImage(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, custom_errors.WrapError(custom_errors.ErrUploadImage, err))
	}

	society := new(models.Society)
	_, err = s.db.Model(society).Where("id = ?", societyId).Update()
	if society.Avatar != "" {
		err = s.DeleteImage(society.Avatar)
		if err != nil {
			log.Errorf("Error delete image of society when uploading new: ", err)
		}
	}

	_, err = s.db.Model(society).Set("avatar = ?", objectName).Where("id = ?", societyId).Update()
	if err != nil {
		_ = s.DeleteImage(objectName)
		return c.JSON(http.StatusInternalServerError, custom_errors.WrapError(custom_errors.ErrUpdateUser, err))
	}

	return c.NoContent(http.StatusCreated)
}
func (s *FileuploadService) UploadTrashImages(c echo.Context) error {
	trashId := c.Param("trashId")

	objectNames, err := s.UploadImages(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, custom_errors.WrapError(custom_errors.ErrUploadImage, err))
	}

	for _, objectName := range objectNames {
		trashImage := new(models.TrashImage)
		trashImage.TrashId = trashId
		trashImage.Url = objectName
		_, err = s.db.Model(trashImage).Insert(trashImage)
		if err != nil {
			_ = s.DeleteImage(objectName)
			return c.JSON(http.StatusInternalServerError, custom_errors.WrapError(custom_errors.ErrUpdateUser, err))
		}
	}

	return c.NoContent(http.StatusCreated)
}

func (s *FileuploadService) UploadCollectionImages(c echo.Context) error {
	collectionId := c.Param("collectionId")

	objectNames, err := s.UploadImages(c)
	if err != nil {
		fmt.Println("ERROR1: ", err)
		return c.JSON(http.StatusInternalServerError, custom_errors.WrapError(custom_errors.ErrUploadImage, err))
	}

	for _, objectName := range objectNames {
		trashImage := new(models.CollectionImage)
		trashImage.CollectionId = collectionId
		trashImage.Url = objectName
		_, err = s.db.Model(trashImage).Insert(trashImage)
		if err != nil {
			fmt.Println("ERROR: ", err)
			_ = s.DeleteImage(objectName)
			return c.JSON(http.StatusInternalServerError, custom_errors.WrapError(custom_errors.ErrUpdateUser, err))
		}
	}

	return c.NoContent(http.StatusCreated)
}

func (s *FileuploadService) GetUserImage(c echo.Context) error {
	imageName := c.Param("image")

	oh := s.bh.Object(imageName)

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	attr, err := oh.Attrs(ctx)
	if err != nil {
		log.Error("ATTR_OBJECT_ERROR", err)
		return c.NoContent(http.StatusInternalServerError)
	}

	rc, err := oh.NewReader(ctx)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}
	defer rc.Close()

	data, err := ioutil.ReadAll(rc)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.Blob(http.StatusOK, attr.ContentType, data)
}

func (s *FileuploadService) GetSocietyImage(c echo.Context) error {
	imageName := c.Param("image")

	oh := s.bh.Object(imageName)

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	attr, err := oh.Attrs(ctx)
	if err != nil {
		log.Error("ATTR_OBJECT_ERROR", err)
		return c.NoContent(http.StatusInternalServerError)
	}

	rc, err := oh.NewReader(ctx)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}
	defer rc.Close()

	data, err := ioutil.ReadAll(rc)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.Blob(http.StatusOK, attr.ContentType, data)
}

func (s *FileuploadService) GetTrashImage(c echo.Context) error {
	imageName := c.Param("image")

	if imageName == "" {
		return c.NoContent(http.StatusConflict)
	}

	oh := s.bh.Object(imageName)

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	attr, err := oh.Attrs(ctx)
	if err != nil {
		log.Error("ATTR_OBJECT_ERROR", err)
		return c.NoContent(http.StatusInternalServerError)
	}

	rc, err := oh.NewReader(ctx)
	if err != nil {
		fmt.Println("Err new reader: ", err)
		return c.NoContent(http.StatusInternalServerError)
	}
	defer rc.Close()

	data, err := ioutil.ReadAll(rc)
	if err != nil {
		fmt.Println("Error read all: ", err)
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.Blob(http.StatusOK, attr.ContentType, data)
}
func (s *FileuploadService) GetCollectionImages(c echo.Context) error {
	imageName := c.Param("image")

	if imageName == "" {
		return c.NoContent(http.StatusConflict)
	}

	oh := s.bh.Object(imageName)

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	attr, err := oh.Attrs(ctx)
	if err != nil {
		log.Error("ATTR_OBJECT_ERROR", err)
		return c.NoContent(http.StatusInternalServerError)
	}

	rc, err := oh.NewReader(ctx)
	if err != nil {
		fmt.Println("Err new reader: ", err)
		return c.NoContent(http.StatusInternalServerError)
	}
	defer rc.Close()

	data, err := ioutil.ReadAll(rc)
	if err != nil {
		fmt.Println("Error read all: ", err)
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.Blob(http.StatusOK, attr.ContentType, data)
}

func (s *FileuploadService) DeleteUserImage(c echo.Context) error {
	userId := c.Get("userId").(string)

	tx, err := s.db.Begin()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, custom_errors.WrapError(custom_errors.ErrDeleteImage, err))
	}
	defer tx.Rollback()

	user := new(models.User)
	err = tx.Model(user).Where("id = ?", userId).Select()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, custom_errors.WrapError(custom_errors.ErrGetUserById, err))
	}
	user.Avatar = ""
	err = tx.Update(user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, custom_errors.WrapError(custom_errors.ErrUpdateUser, err))
	}

	if user.Avatar != "" {
		err = s.DeleteImage(user.Avatar)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, custom_errors.WrapError(custom_errors.ErrDeleteImage, err))
		}
	}

	err = tx.Commit()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, custom_errors.WrapError(custom_errors.ErrDeleteImage, err))
	}

	return c.NoContent(http.StatusOK)
}

//
//func (s *FileuploadService) DeleteSocietyImage() error {
//
//}
func (s *FileuploadService) DeleteTrashImage(c echo.Context) error {
	imageName := c.Param("image")
	trashId := c.Param("trashId")

	tx, err := s.db.Begin()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, custom_errors.WrapError(custom_errors.ErrDeleteImage, err))
	}
	defer tx.Rollback()

	imageDb := new(models.TrashImage)
	_, err = tx.Model(imageDb).Where("url = ? and trash_id = ?", imageName, trashId).Delete()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, custom_errors.WrapError(custom_errors.ErrDeleteImage, err))
	}

	err = s.DeleteImage(imageName)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, custom_errors.WrapError(custom_errors.ErrDeleteImage, err))
	}

	err = tx.Commit()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, custom_errors.WrapError(custom_errors.ErrDeleteImage, err))
	}

	return c.NoContent(http.StatusOK)
}
func (s *FileuploadService) DeleteCollectionImages(c echo.Context) error {
	userId := c.Get("userId").(string)
	collectionId := c.Param("collectionId")
	idsString := c.QueryParam("ids")
	ids := strings.Split(idsString, ",")

	err := s.db.Model(&models.UserCollection{}).Where("user_id = ? and collection_id = ?", userId, collectionId).First()
	if err == pg.ErrNoRows {
		return c.JSON(http.StatusInternalServerError, custom_errors.WrapError(custom_errors.ErrDeleteCollectionImage, err))
	} else if err != nil {
		return c.JSON(http.StatusInternalServerError, custom_errors.WrapError(custom_errors.ErrDeleteCollectionImage, err))
	}

	tx, err := s.db.Begin()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, custom_errors.WrapError(custom_errors.ErrDeleteCollectionImage, err))
	}
	defer tx.Rollback()

	for _, imageName := range ids {
		imageDb := new(models.CollectionImage)
		_, err = tx.Model(imageDb).Where("url = ? and collection_id = ?", imageName, collectionId).Delete()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, custom_errors.WrapError(custom_errors.ErrDeleteCollectionImage, err))
		}

		err = s.DeleteImage(imageName)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, custom_errors.WrapError(custom_errors.ErrDeleteCollectionImage, err))
		}

		err = tx.Commit()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, custom_errors.WrapError(custom_errors.ErrDeleteCollectionImage, err))
		}
	}

	return c.NoContent(http.StatusOK)
}

func (s *FileuploadService) DeleteEventsCollectionsImages(c echo.Context) error {
	userId := c.Get("userId").(string)
	collectionId := c.Param("collectionId")

	eventId := c.QueryParam("eventId")
	pickerId := c.QueryParam("pickerId")
	requestAsSociety := c.QueryParam("asSociety")
	idsString := c.QueryParam("ids")
	ids := strings.Split(idsString, ",")

	fmt.Println("EventId: ", eventId)
	fmt.Println("AS soc: ", requestAsSociety)
	fmt.Println("pickerId: ", pickerId)
	fmt.Println("idsImages: ", ids)
	fmt.Println("Collection ID", collectionId)

	asSociety, err := strconv.ParseBool(requestAsSociety)
	if err != nil {
		return c.JSON(http.StatusBadRequest, custom_errors.WrapError(custom_errors.ErrBindingRequest, err))
	}

	if asSociety {
		socMembership := new(models.Member)
		err := s.db.Model(socMembership).Where("user_id = ? and society_id = ?", userId, pickerId).Select()
		if err != nil {
			log.Error(err)
			return c.JSON(http.StatusInternalServerError, custom_errors.WrapError(custom_errors.ErrDeleteCollectionImage, err))
		}
		if socMembership.Permission == "member" {
			return c.JSON(http.StatusForbidden, custom_errors.WrapError(custom_errors.ErrDeleteCollectionImage, fmt.Errorf("You have no rights to act as society editor! ")))
		}

		permission := new(models.EventSociety)
		err = s.db.Model(permission).Where("society_id = ? and event_id = ?", pickerId, eventId).Select()
		if err != nil {
			fmt.Println("TU")
			log.Error(err)
			return c.JSON(http.StatusInternalServerError, custom_errors.WrapError(custom_errors.ErrDeleteCollectionImage, err))
		}
		if permission.Permission == "viewer" {
			return c.JSON(http.StatusInternalServerError, custom_errors.WrapError(custom_errors.ErrDeleteCollectionImage, fmt.Errorf("You have no rights to edit this collection! ")))
		}
	} else {
		permission := new(models.EventUser)
		err := s.db.Model(permission).Where("user_id = ? and event_id = ?", userId, eventId).Select()
		if err != nil {
			log.Error(err)
			return c.JSON(http.StatusInternalServerError, custom_errors.WrapError(custom_errors.ErrDeleteCollectionImage, err))
		}
		if permission.Permission == "viewer" {
			return c.JSON(http.StatusInternalServerError, custom_errors.WrapError(custom_errors.ErrDeleteCollectionImage, fmt.Errorf("You have no rights to edit this collection! ")))
		}
	}

	for _, imageName := range ids {
		imageDb := new(models.CollectionImage)
		res, err := s.db.Model(imageDb).Where("url = ? and collection_id = ?", imageName, collectionId).Delete()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, custom_errors.WrapError(custom_errors.ErrDeleteCollectionImage, err))
		}
		fmt.Println("In reality I deleted: ", res)

		err = s.DeleteImage(imageName)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, custom_errors.WrapError(custom_errors.ErrDeleteCollectionImage, err))
		}
	}

	return c.NoContent(http.StatusOK)
}
