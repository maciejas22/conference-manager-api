package service

import (
	"time"
)

type Subsection struct {
	Id        int       `json:"id" db:"id"`
	SectionId int       `json:"section_id" db:"section_id"`
	Title     string    `json:"title" db:"title"`
	Content   *string   `json:"content" db:"content"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type Section struct {
	Id               int       `json:"id" db:"id"`
	TermsOfServiceId int       `json:"terms_of_service_id" db:"terms_of_service_id"`
	Title            string    `json:"title" db:"title"`
	Content          *string   `json:"content" db:"content"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
	Subsections      []*Subsection
}

type TermsOfService struct {
	Id              int       `json:"id" db:"id"`
	CreatedAt       time.Time `json:"created_at"     db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"     db:"updated_at"`
	Introduction    string    `json:"introduction"  db:"introduction"`
	Acknowledgement string    `json:"acknowledgement" db:"acknowledgement"`
	Sections        []*Section
}
