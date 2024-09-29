package files

import (
	"context"
	"io"
	"strconv"

	"github.com/maciejas22/conference-manager/api/internal/config"
	"github.com/maciejas22/conference-manager/api/internal/models"
	"github.com/maciejas22/conference-manager/api/pkg/s3"
)

func createKey(conferenceId int, fileName string) string {
	return strconv.Itoa(conferenceId) + "/" + fileName
}

func GetConferenceFiles(ctx context.Context, s3 *s3.S3Client, conferenceId int) ([]*models.File, error) {
	files, err := s3.GetFilesFromFolder(config.AppConfig.S3BucketsConferenceFiles, strconv.Itoa(conferenceId))
	if err != nil {
		return nil, err
	}
	var parsedFiles []*models.File
	for _, f := range files {
		parsedFiles = append(parsedFiles, &models.File{
			Key:  f.Key,
			URL:  f.URL,
			Size: int(f.Size),
		})
	}

	return parsedFiles, nil
}

func UploadConferenceFile(ctx context.Context, s3 *s3.S3Client, conferenceId int, fileName string, fileContent io.ReadSeeker) error {
	err := s3.UploadFile(config.AppConfig.S3BucketsConferenceFiles, createKey(conferenceId, fileName), fileContent)
	if err != nil {
		return err
	}

	return nil
}

func DeleteConferenceFile(ctx context.Context, s3 *s3.S3Client, fileKey string) error {
	return s3.DeleteFile(config.AppConfig.S3BucketsConferenceFiles, fileKey)
}
