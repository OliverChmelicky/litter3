package fileupload

import (
	"cloud.google.com/go/storage"
	"context"
	firebase "firebase.google.com/go"
	"github.com/go-pg/pg/v9"
	log "github.com/sirupsen/logrus"
	"google.golang.org/api/option"
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

//
//func (s *FileuploadService) UploadSocietyImage() error {
//
//}
//func (s *FileuploadService) UploadTrashImages() error {
//
//}
//func (s *FileuploadService) UploadCollectionImages() error {
//
//}

//func (s *FileuploadService) GetUserImage() error {
//
//}
//
//func (s *FileuploadService) GetSocietyImage() error {
//
//}
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
