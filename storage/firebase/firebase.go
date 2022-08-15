package firebase

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"cloud.google.com/go/storage"

	firebase "firebase.google.com/go"

	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

type Config struct {
	KeyPath       string
	StorageBucket string
}

type Firebase struct {
	bucket *storage.BucketHandle
}

func NewStorage(config Config) *Firebase {
	firebaseConfig := &firebase.Config{
		StorageBucket: config.StorageBucket,
	}
	ctx := context.Background()
	opt := option.WithCredentialsFile(config.KeyPath)
	app, err := firebase.NewApp(ctx, firebaseConfig, opt)
	if err != nil {
		log.Fatalln(err.Error())
	}
	client, err := app.Storage(ctx)
	if err != nil {
		log.Fatalln(err.Error())
	}
	bucket, err := client.DefaultBucket()
	if err != nil {
		log.Fatalln(err.Error())
	}
	return &Firebase{
		bucket: bucket,
	}
}

func (fb Firebase) Upload(file *os.File) {
	contentType := "text/plain"
	ctx := context.Background()

	remoteFilename := filepath.Base(file.Name())
	writer := fb.bucket.Object(remoteFilename).NewWriter(ctx)
	writer.ObjectAttrs.ContentType = contentType
	writer.ObjectAttrs.CacheControl = "no-cache"
	writer.ObjectAttrs.ACL = []storage.ACLRule{
		{
			Entity: storage.AllUsers,
			Role:   storage.RoleReader,
		},
	}
	if _, err := io.Copy(writer, file); err != nil {
		log.Fatalln(err.Error())
	}

	if err := writer.Close(); err != nil {
		log.Fatalln(err.Error())
	}
}

func (fb Firebase) List() error {
	ctx := context.Background()
	it := fb.bucket.Objects(ctx, nil)
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}
		fmt.Println(attrs.Name)
	}
	return nil
}

func (fb Firebase) GetFileList() ([]string, error) {
	list := []string{}
	ctx := context.Background()
	it := fb.bucket.Objects(ctx, nil)
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		list = append(list, attrs.Name)
	}
	return list, nil
}

func (fb Firebase) Download(fileName string) []byte {
	ctx := context.Background()
	rc, err := fb.bucket.Object(fileName).NewReader(ctx)
	if err != nil {
		log.Fatalln(err.Error())
	}
	defer rc.Close()

	data, err := ioutil.ReadAll(rc)
	if err != nil {
		log.Fatalln(err.Error())
	}
	return data
}
