package repositories

import (
	"context"
	"io"

	"github.com/maciejas22/conference-manager/api/internal/config"
	"github.com/maciejas22/conference-manager/api/internal/models"
	"github.com/maciejas22/conference-manager/api/pkg/s3"
)

func createKey(conferenceId string, fileName string) string {
	return conferenceId + "/" + fileName
}

func GetFiles(ctx context.Context, s3 *s3.S3Client, path string) ([]*models.File, error) {
	files, err := s3.GetFilesFromFolder(config.AppConfig.S3BucketsConferenceFiles, path)
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

func UploadFile(ctx context.Context, s3 *s3.S3Client, path string, fileName string, fileContent io.ReadSeeker) error {
	err := s3.UploadFile(config.AppConfig.S3BucketsConferenceFiles, createKey(path, fileName), fileContent)
	if err != nil {
		return err
	}

	return nil
}

func DeleteFile(ctx context.Context, s3 *s3.S3Client, fileId string) error {
	return s3.DeleteFile(config.AppConfig.S3BucketsConferenceFiles, fileId)
}
