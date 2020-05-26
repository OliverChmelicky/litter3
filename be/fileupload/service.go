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

func (s *FileuploadService) UploadUserImages(c echo.Context) error {
	userId := c.Param("userId")

	objectName, err := s.UploadImage(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, custom_errors.WrapError(custom_errors.ErrUploadImage, err))
	}

	user := new(models.User)
	_, err = s.db.Model(user).Set("avatar = ?", objectName).Where("id = ?", userId).Update()
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

//func (s *FileuploadService) GetUserImage() error {
//
//}
//
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

//func (s *FileuploadService) DeleteUserImage() error {
//
//}
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
	imageName := c.Param("image")
	collectionId := c.Param("collectionId")

	tx, err := s.db.Begin()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, custom_errors.WrapError(custom_errors.ErrDeleteCollectionImage, err))
	}
	defer tx.Rollback()

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

	return c.NoContent(http.StatusOK)
}
