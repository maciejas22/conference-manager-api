package s3

import (
	"context"
	"fmt"
	"io"
	"log/slog"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/maciejas22/conference-manager/api/internal/config"
)

type S3Client struct {
	s3     *s3.Client
	logger *slog.Logger
}

func NewS3Client(ctx context.Context, logger *slog.Logger) (*S3Client, error) {
	s3Cfg, err := NewS3Config(ctx, logger)
	if err != nil {
		return nil, err
	}

	s3Client := s3.NewFromConfig(s3Cfg, func(o *s3.Options) {
		o.UsePathStyle = true
	})

	return &S3Client{
		s3:     s3Client,
		logger: logger,
	}, nil
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
		client.logger.Error("error uploading file to s3 bucket", "error", err)
	}

	return err
}

func (client *S3Client) GetFilesFromFolder(ctx context.Context, bucket string, folder string) ([]File, error) {
	result, err := client.s3.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket: aws.String(bucket),
		Prefix: aws.String(folder),
	})
	if err != nil {
		client.logger.Error("error getting files from s3 bucket", "error", err)
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
		client.logger.Error("error deleting file from s3 bucket", "error", err)
		return err
	}

	return nil
}
