package converters

import (
	"encoding/base64"
	"io"
	"strings"

	"github.com/maciejas22/conference-manager/api/internal/models"
)

type File struct {
	Name    string
	Content io.ReadSeeker
}

func ConvertUploadFileToReader(upload *models.UploadFile) (*File, error) {
	parts := strings.Split(upload.Base64Content, ",")
	if len(parts) > 1 {
		upload.Base64Content = parts[1]
	}

	data, err := base64.StdEncoding.DecodeString(upload.Base64Content)
	if err != nil {
		return nil, err
	}

	fileReader := strings.NewReader(string(data))

	return &File{
		Name:    upload.FileName,
		Content: fileReader,
	}, nil
}
