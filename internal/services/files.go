package services

import (
	"context"

	"github.com/maciejas22/conference-manager/api/internal/config"
	"github.com/maciejas22/conference-manager/api/internal/models"
	"github.com/maciejas22/conference-manager/api/pkg/s3"
)

func GetConferenceFiles(ctx context.Context, s3 *s3.S3Client, conferenceId string) ([]*models.File, error) {
	files, err := s3.GetFilesFromFolder(config.AppConfig.S3BucketsConferenceFiles, conferenceId)
	if err != nil {
		return nil, err
	}
	var parsedFiles []*models.File
	for _, f := range files {
		parsedFiles = append(parsedFiles, &models.File{
			ID:   f.ID,
			Name: f.Name,
			URL:  f.URL,
			Size: f.Size,
		})
	}

	return parsedFiles, nil
}
