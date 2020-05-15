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

//func (s *FileuploadService) UploadCollectionImages() error {
//
//}

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

	fmt.Println("dlzka dat je:")
	fmt.Println(len(data))

	return c.Blob(http.StatusOK, attr.ContentType, data)
}



//func (s *FileuploadService) GetTrashImages() error {
//
//}
//func (s *FileuploadService) GetCollectionImages() error {
//
//}

//func (s *FileuploadService) DeleteUserImage() error {
//
//}
//
//func (s *FileuploadService) DeleteSocietyImage() error {
//
//}
//func (s *FileuploadService) DeleteTrashImages() error {
//
//}
//func (s *FileuploadService) DeleteCollectionImages() error {
//
//}
