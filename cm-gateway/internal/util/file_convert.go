package util

import (
	"encoding/base64"
	"errors"
	"io"
	"strings"
)

func FileConvert(base64Content string) (io.ReadSeeker, error) {
	fileParts := strings.Split(base64Content, ",")
	if len(fileParts) > 1 {
		base64Content = fileParts[1]
	}

	fileData, err := base64.StdEncoding.DecodeString(base64Content)
	if err != nil {
		return nil, errors.New("Failed to decode base64 content")
	}

	return strings.NewReader(string(fileData)), nil
}
