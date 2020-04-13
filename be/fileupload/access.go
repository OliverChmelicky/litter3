package fileupload

import (
	"cloud.google.com/go/storage"
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"time"
)

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

func (s *fileuploadService) LoadImage(bucketName, objectName string) (string, io.Reader, error) {
	oh := s.storage.Bucket(bucketName).Object(objectName)
	attr, err := oh.Attrs(context.Background())
	if err != nil {
		log.Error("ATTR_OBJECT_ERROR", err)
		return "", nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	reader, err := oh.NewReader(ctx)
	if err != nil {
		log.Error("READER_OBJECT_ERROR", err)
		return "", nil, err
	}

	return attr.ContentType, reader, nil
}
