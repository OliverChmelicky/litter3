package fileupload

import (
	"cloud.google.com/go/storage"
	"context"
	firebase "firebase.google.com/go"
	log "github.com/sirupsen/logrus"
	"google.golang.org/api/option"
)

type fileuploadService struct {
	bh *storage.BucketHandle
}

func CreateService(opt option.ClientOption, bucketName string) *fileuploadService {
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

	return &fileuploadService{bh}
}

//func (s *fileuploadService) UploadUserImage() error {
//
//}
//
//func (s *fileuploadService) UploadSocietyImage() error {
//
//}
//func (s *fileuploadService) UploadTrashImage() error {
//
//}
//func (s *fileuploadService) UploadCollectionImage() error {
//
//}

//func (s *fileuploadService) GetUserImage() error {
//
//}
//
//func (s *fileuploadService) GetSocietyImage() error {
//
//}
//func (s *fileuploadService) GetTrashImages() error {
//
//}
//func (s *fileuploadService) GetCollectionImages() error {
//
//}

//func (s *fileuploadService) DeleteUserImage() error {
//
//}
//
//func (s *fileuploadService) DeleteSocietyImage() error {
//
//}
//func (s *fileuploadService) DeleteTrashImages() error {
//
//}
//func (s *fileuploadService) DeleteCollectionImages() error {
//
//}
