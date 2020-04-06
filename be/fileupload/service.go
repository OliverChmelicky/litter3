package main

import (
	"cloud.google.com/go/storage"
	"context"
	firebase "firebase.google.com/go"
	"fmt"
	log "github.com/sirupsen/logrus"
	"google.golang.org/api/option"
	"io"
	"os"
)

func main() {
	var r io.Reader
	f, err := os.Open("keep.jpg")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	r = f

	name := "keep-doing.jpg"

	ctx := context.Background()
	_, objAttrs, err := upload(ctx, r, "litter3-olo-gcp", "litter3-olo-gcp.appspot.com", name, true)
	if err != nil {
		switch err {
		case storage.ErrBucketNotExist:
			log.Fatal("Please create the bucket first e.g. with `gsutil mb`")
		default:
			log.Fatal(err)
		}
	}

	log.Printf("URL: %s", objectURL(objAttrs))
	log.Printf("Size: %d", objAttrs.Size)
	log.Printf("MD5: %x", objAttrs.MD5)
	log.Printf("objAttrs: %+v", objAttrs)
}

func objectURL(objAttrs *storage.ObjectAttrs) string {
	return fmt.Sprintf("https://storage.googleapis.com/%s/%s", objAttrs.Bucket, objAttrs.Name)
}

func upload(ctx context.Context, r io.Reader, projectID, bucket, name string, public bool) (*storage.ObjectHandle, *storage.ObjectAttrs, error) {
	opt := option.WithCredentialsFile("../secrets/litter3-olo-gcp-firebase-adminsdk-6ar5p-9f1130c1cc.json")
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("error initializing app: %s\n", err.Error())
	}
	st, err := app.Storage(context.Background())
	bh, err := st.Bucket(bucket)
	if err != nil {
		log.Fatalf("faking err :", err)
	}
	// Next check if the bucket exists
	if _, err = bh.Attrs(ctx); err != nil {
		return nil, nil, err
	}

	obj := bh.Object(name)
	w := obj.NewWriter(ctx)
	if _, err := io.Copy(w, r); err != nil {
		return nil, nil, err
	}
	if err := w.Close(); err != nil {
		return nil, nil, err
	}

	if public {
		if err := obj.ACL().Set(ctx, storage.AllUsers, storage.RoleReader); err != nil {
			return nil, nil, err
		}
	}

	attrs, err := obj.Attrs(ctx)
	return obj, attrs, err
}
