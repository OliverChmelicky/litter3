package fileupload

import (
	"cloud.google.com/go/storage"
	"context"
	firebase "firebase.google.com/go"
	"fmt"
	log "github.com/sirupsen/logrus"
	"google.golang.org/api/option"
	"io"
)

type fileuploadService struct {
	bh *storage.BucketHandle
}

func CreateService(firebsaeCredentials string, bucketName string) *fileuploadService {
	//opt := option.WithCredentialsFile("../secrets/litter3-olo-gcp-firebase-adminsdk-6ar5p-9f1130c1cc.json")
	opt := option.WithCredentialsFile(firebsaeCredentials)
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

func (s *fileuploadService) Upload(ctx context.Context, r io.Reader, name string) (string, error) {
	obj := s.bh.Object(name)
	w := obj.NewWriter(ctx)
	if _, err := io.Copy(w, r); err != nil {
		return "", err
	}
	if err := w.Close(); err != nil {
		return "", err
	}

	if err := obj.ACL().Set(ctx, storage.AllUsers, storage.RoleReader); err != nil {
		return "", err
	}

	attrs, err := obj.Attrs(ctx)
	return objectURL(attrs), err
}

func objectURL(objAttrs *storage.ObjectAttrs) string {
	return fmt.Sprintf("https://storage.googleapis.com/%s/%s", objAttrs.Bucket, objAttrs.Name)
}
