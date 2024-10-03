package s3

import (
	"context"
	"log"
	"log/slog"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	s3Config "github.com/aws/aws-sdk-go-v2/config"
)

// func NewS3Session(logger *slog.Logger) (*s3.S3, error) {
// 	sess, err := session.NewSession(&aws.Config{
// 		Region:           aws.String(config.AppConfig.S3Region),
// 		Endpoint:         aws.String(config.AppConfig.S3Endpoint),
// 		Credentials:      credentials.NewStaticCredentials(config.AppConfig.S3AccessKeyId, config.AppConfig.S3SecretAccessKey, ""),
// 		S3ForcePathStyle: aws.Bool(true),
// 	})
//
// 	if err != nil {
// 		logger.Error("failed to create new S3 session", "error", err)
// 		return nil, errors.New("failed to create new S3 session")
// 	}
//
// 	return s3.New(sess), nil
// }

func NewS3Config(ctx context.Context, logger *slog.Logger) (aws.Config, error) {
	log.Println("AWS_REGION", os.Getenv("AWS_REGION"))
	cfg, err := s3Config.LoadDefaultConfig(ctx)
	if err != nil {
		logger.Error("failed to load S3 config", "error", err)
		panic(err)
		return aws.Config{}, err
	}

	return cfg, nil
}
