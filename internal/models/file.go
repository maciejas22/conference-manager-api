package models

type File struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	URL  string `json:"url"`
	Size int64  `json:"size"`
}
