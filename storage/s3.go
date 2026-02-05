package storage

import (
	"context"
	"mime/multipart"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var S3 *s3.Client

func InitS3() {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic("unable to load AWS config")
	}

	S3 = s3.NewFromConfig(cfg)
}

func StoreFileInS3(fileHeader *multipart.FileHeader, storageKey string) error {
	file, err := fileHeader.Open()

	if err != nil {
		return err
	}
	defer file.Close()

	_, err = S3.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(os.Getenv("S3_BUCKET")),
		Key:         aws.String(storageKey),
		Body:        file,
		ContentType: aws.String(fileHeader.Header.Get("Content-Type")),
	})

	return err
}

func DeleteFileFromS3(storageKey string) error {
	_, err := S3.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: aws.String(os.Getenv("S3_BUCKET")),
		Key:    aws.String(storageKey),
	})

	return err
}
