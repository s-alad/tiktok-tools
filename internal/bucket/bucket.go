package bucket

import (
	"context"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	log "github.com/sirupsen/logrus"
)

type Bucket struct {
	id       string
	name     string
	client   *s3.Client
	uploader *manager.Uploader
}

func (b *Bucket) ID() string {
	return b.id
}

func (b *Bucket) init() {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	b.client = s3.NewFromConfig(cfg)
	b.uploader = manager.NewUploader(b.client)
	b.name = "libtok-of-alexandria"

	log.WithFields(log.Fields{
		"id":   b.id,
		"name": b.name,
	}).Info("bucket initialized")
}

func (b *Bucket) Path(key string) string {
	return fmt.Sprintf("https://%s.s3.amazonaws.com/%s", b.name, key)
}

func (b *Bucket) Exists(key string) (bool, string, error) {
	_, err := b.client.HeadObject(context.TODO(), &s3.HeadObjectInput{
		Bucket: aws.String(b.name),
		Key:    aws.String(key),
	})
	if err != nil {
		return false, "", err
	}

	return true, fmt.Sprintf("https://%s.s3.amazonaws.com/%s", b.name, key), nil
}

func (b *Bucket) Upload(key string, content io.Reader, contentType string) (string, error) {
	if exists, loc, _ := b.Exists(key); exists {
		log.WithFields(log.Fields{
			"key":    key,
			"bucket": b.name,
		}).Warn("file already exists in bucket")
		return loc, nil
	}

	result, err := b.uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(b.name),
		Key:         aws.String(key),
		Body:        content,
		ContentType: &contentType,
		/* ACL:    types.ObjectCannedACLPublicRead, */
	})

	if err != nil {
		log.WithField("key", key).Errorf("failed to upload file: %v", err)
		return "", err
	} else {
		log.WithFields(log.Fields{
			"key":      key,
			"bucket":   b.name,
			"location": result.Location,
		}).Info("file successfully uploaded to bucket")
	}

	return result.Location, nil
}

func Create(id string) *Bucket {
	b := &Bucket{id: id}
	b.init()
	return b
}
