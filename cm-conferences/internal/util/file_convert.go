package util

import (
	"encoding/base64"
	"errors"
	"io"
	"strings"
)

func FileConvert(base64Content string) (io.ReadSeeker, error) {
	fileData, err := base64.StdEncoding.DecodeString(base64Content)
	if err != nil {
		return nil, errors.New("failed to decode base64 content")
	}

	return strings.NewReader(string(fileData)), nil
}
