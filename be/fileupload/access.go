package fileupload

import (
	"context"
	"fmt"
	"github.com/labstack/echo"
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
	"io"
	"mime"
	"path/filepath"
	"time"
)

func (s *FileuploadService) UploadImage(c echo.Context) (string, error) {
	ctx := context.Background()

	file, err := c.FormFile("file")
	if err != nil {
		log.Error("FORM_FILE_ERROR", err)
		return "", fmt.Errorf("Error extracting image %w ", err)
	}

	src, err := file.Open()
	if err != nil {
		log.Error("FILE_OPEN_ERROR", err)
		return "", fmt.Errorf("Error opening image %w ", err)
	}
	defer src.Close()

	sufix := filepath.Ext(file.Filename)
	objectName := uuid.NewV4().String() + sufix

	obj := s.bh.Object(objectName)
	w := obj.NewWriter(ctx)
	w.ContentType = mime.TypeByExtension(sufix)
	if _, err := io.Copy(w, src); err != nil {
		return "", fmt.Errorf("Error saving image %w ", err)
	}
	if err := w.Close(); err != nil {
		_ = s.DeleteImage(objectName)
		return "", fmt.Errorf("Error closing image %w ", err)
	}

	//If images should be accessible from the internet
	//if err := obj.ACL().Set(ctx, storage.AllUsers, storage.RoleReader); err != nil {
	//	return "", fmt.Errorf("Error updating image attributes %w ", err)
	//}

	return objectName, err
}

func (s *FileuploadService) UploadImages(c echo.Context) ([]string,error) {
	form, err := c.MultipartForm()
	if err != nil {
		return []string{},err
	}
	files := form.File["files"]

	var fileNames []string
	for _,file := range files {
		src, err := file.Open()
		if err != nil {
			log.Error("FILE_OPEN_ERROR", err)
			return []string{}, fmt.Errorf("Error opening image %w ", err)
		}
		defer src.Close()

		sufix := filepath.Ext(file.Filename)
		objectName := uuid.NewV4().String() + sufix

		obj := s.bh.Object(objectName)
		w := obj.NewWriter(context.Background())
		w.ContentType = mime.TypeByExtension(sufix)
		if _, err := io.Copy(w, src); err != nil {
			return []string{}, fmt.Errorf("Error saving image %w ", err)
		}
		if err := w.Close(); err != nil {
			_ = s.DeleteImage(objectName)
			return []string{}, fmt.Errorf("Error closing image %w ", err)
		}

		//if images should be accesible from the internet
		//if err := obj.ACL().Set(context.Background(), storage.AllUsers, storage.RoleReader); err != nil {
		//	return []string{}, fmt.Errorf("Error updating image attributes %w ", err)
		//}

		fileNames = append(fileNames, objectName)
	}

	return fileNames, err
}

func (s *FileuploadService) LoadImage(objectName string) (string, io.Reader, error) {
	oh := s.bh.Object(objectName)
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

func (s *FileuploadService) DeleteImage(objectName string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	o := s.bh.Object(objectName)
	if err := o.Delete(ctx); err != nil {
		return err
	}

	return nil
}
