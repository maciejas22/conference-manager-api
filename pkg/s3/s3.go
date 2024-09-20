package s3

import (
	"fmt"
	"io"
	"log/slog"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/maciejas22/conference-manager/api/internal/config"
)

type S3Client struct {
	s3     *s3.S3
	logger *slog.Logger
}

func NewS3Client(logger *slog.Logger) (*S3Client, error) {
	s3Session, err := NewS3Session(logger)
	if err != nil {
		return nil, err
	}
	return &S3Client{s3: s3Session, logger: logger}, nil
}

type File struct {
	ID   string
	Name string
	Size int64
	URL  string
}

func (client *S3Client) UploadFile(bucket string, key string, file io.ReadSeeker) error {
	_, err := client.s3.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   file,
	})

	if err != nil {
		client.logger.Error("error uploading file to s3 bucket", "error", err)
	}

	return err
}

func (client *S3Client) GetFilesFromFolder(bucket string, folder string) ([]File, error) {
	input := &s3.ListObjectsV2Input{
		Bucket: aws.String(bucket),
		Prefix: aws.String(folder),
	}

	result, err := client.s3.ListObjectsV2(input)
	if err != nil {
		client.logger.Error("error getting files from s3 bucket", "error", err)
		return nil, err
	}

	var files []File
	for _, item := range result.Contents {
		url := fmt.Sprintf("%s/%s/%s", config.AppConfig.S3Endpoint, bucket, *item.Key)

		file := File{
			ID:   *item.Key,
			Name: *item.Key,
			Size: *item.Size,
			URL:  url,
		}
		files = append(files, file)
	}

	return files, nil
}

func (client *S3Client) DeleteFile(bucket string, key string) error {
	_, err := client.s3.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})

	if err != nil {
		client.logger.Error("error deleting file from s3 bucket", "error", err)
		return err
	}

	return client.s3.WaitUntilObjectNotExists(&s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
}
