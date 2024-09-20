package models

import "time"

type SubSection struct {
	ID      int     `json:"id"`
	Title   *string `json:"title"`
	Content *string `json:"content"`
}

type Section struct {
	ID          int           `json:"id"`
	Title       *string       `json:"title"`
	Content     *string       `json:"content"`
	SubSections *[]SubSection `json:"subsections"`
}

type TermsOfService struct {
	ID              int        `json:"id"`
	UpdatedAt       time.Time  `json:"updated_at"     db:"updated_at"`
	Introduction    string     `json:"introduction"`
	Acknowledgement string     `json:"acknowledgement"`
	Sections        *[]Section `json:"sections"`
}

type News struct {
	ID      int       `json:"id"`
	Title   string    `json:"title"`
	Content string    `json:"content"`
	Date    time.Time `json:"date"    db:"created_at"`
}
