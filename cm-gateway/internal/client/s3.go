package client

import (
	"context"
	"fmt"
	"io"
	"log"
	"log/slog"

	"github.com/aws/aws-sdk-go-v2/aws"
	s3Config "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/maciejas22/conference-manager-api/cm-gateway/config"
)

type S3Client struct {
	s3     *s3.Client
	logger *slog.Logger
}

type S3Service interface {
	UploadFile(ctx context.Context, bucket string, key string, file io.ReadSeeker) error
	GetFilesFromFolder(ctx context.Context, bucket string, folder string) ([]File, error)
	DeleteFile(ctx context.Context, bucket string, key string) error
}

func NewS3Client(ctx context.Context, logger *slog.Logger) S3Service {
	cfg, err := s3Config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatalf("failed to load S3 config: %v", err)
	}

	s3Client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.UsePathStyle = true
	})

	return &S3Client{
		s3:     s3Client,
		logger: logger,
	}
}

type File struct {
	Key  string
	Size int64
	URL  string
}

func (client *S3Client) UploadFile(ctx context.Context, bucket string, key string, file io.ReadSeeker) error {
	_, err := client.s3.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   file,
	})

	if err != nil {
		client.logger.Error("S3 failed to upload file to bucket", "error", err)
	}

	return err
}

func (client *S3Client) GetFilesFromFolder(ctx context.Context, bucket string, folder string) ([]File, error) {
	result, err := client.s3.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket: aws.String(bucket),
		Prefix: aws.String(folder),
	})
	if err != nil {
		client.logger.Error("S3 failed to list objects in bucket", "error", err)
		return nil, err
	}

	var files []File
	for _, item := range result.Contents {
		url := fmt.Sprintf("%s/%s/%s", config.AppConfig.AWSEndpoint, bucket, *item.Key)

		file := File{
			Key:  *item.Key,
			Size: *item.Size,
			URL:  url,
		}
		files = append(files, file)
	}

	return files, nil
}

func (client *S3Client) DeleteFile(ctx context.Context, bucket string, key string) error {
	_, err := client.s3.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})

	if err != nil {
		client.logger.Error("S3 failed to delete object from bucket", "error", err)
		return err
	}

	return nil
}
