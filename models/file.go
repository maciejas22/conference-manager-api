package models

import "github.com/maciejas22/conference-manager/api/pkg/s3"

type File struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	URL  string `json:"url"`
	Size int64  `json:"size"`
}

func (f *File) ToRepo() *s3.File {
	return &s3.File{
		ID:   f.ID,
		Name: f.Name,
		URL:  f.URL,
		Size: f.Size,
	}
}
