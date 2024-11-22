package repository

import (
	"time"

	s "github.com/maciejas22/conference-manager-api/cm-info/internal/repository/section"
)

type TermsOfService struct {
	Id              int       `json:"id" db:"id"`
	CreatedAt       time.Time `json:"created_at"     db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"     db:"updated_at"`
	Introduction    string    `json:"introduction"  db:"introduction"`
	Acknowledgement string    `json:"acknowledgement" db:"acknowledgement"`
	Sections        []*s.Section
}
