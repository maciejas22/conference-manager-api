package repository

import (
	"time"

	ss "github.com/maciejas22/conference-manager-api/cm-info/internal/repository/subsection"
)

type Section struct {
	Id               int       `json:"id" db:"id"`
	TermsOfServiceId int       `json:"terms_of_service_id" db:"terms_of_service_id"`
	Title            string    `json:"title" db:"title"`
	Content          *string   `json:"content" db:"content"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time `json:"updated_at" db:"updated_at"`
	Subsections      []*ss.Subsection
}
