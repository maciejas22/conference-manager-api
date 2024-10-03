package s3

import (
	"context"
	"log"
	"log/slog"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	s3Config "github.com/aws/aws-sdk-go-v2/config"
)

func NewS3Config(ctx context.Context, logger *slog.Logger) (aws.Config, error) {
	log.Println("AWS_REGION", os.Getenv("AWS_REGION"))
	cfg, err := s3Config.LoadDefaultConfig(ctx)
	if err != nil {
		logger.Error("failed to load S3 config", "error", err)
		return aws.Config{}, err
	}

	return cfg, nil
}
