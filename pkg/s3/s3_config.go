package s3

import (
	"errors"
	"log/slog"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/maciejas22/conference-manager/api/internal/config"
)

func NewS3Session(logger *slog.Logger) (*s3.S3, error) {
	sess, err := session.NewSession(&aws.Config{
		Region:           aws.String(config.AppConfig.S3Region),
		Endpoint:         aws.String(config.AppConfig.S3Endpoint),
		Credentials:      credentials.NewStaticCredentials(config.AppConfig.S3AccessKeyID, config.AppConfig.S3SecretAccessKey, ""),
		S3ForcePathStyle: aws.Bool(true),
	})

	if err != nil {
		logger.Error("failed to create new S3 session", "error", err)
		return nil, errors.New("failed to create new S3 session")
	}

	return s3.New(sess), nil
}
