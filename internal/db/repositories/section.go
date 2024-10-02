package repositories

import (
	"github.com/jmoiron/sqlx"
)

type Section struct {
	Id               int     `json:"id" db:"id"`
	TermsOfServiceId int     `json:"terms_of_service_id" db:"terms_of_service_id"`
	Title            string  `json:"title" db:"title"`
	Content          *string `json:"content" db:"content"`
	CreatedAt        string  `json:"created_at" db:"created_at"`
	JoinedAt         string  `json:"joined_at" db:"joined_at"`
}

func (s *Section) TableName() string {
	return "public.sections"
}

func GetToSSections(tx *sqlx.Tx, termsOfServiceId int) ([]Section, error) {
	var sections []Section
	s := &Section{}
	query := "SELECT id, terms_of_service_id, title, content FROM " + s.TableName() + " WHERE terms_of_service_id = $1"
	err := tx.Select(
		&sections,
		query,
		termsOfServiceId,
	)
	if err != nil {
		return nil, err
	}
	return sections, nil
}
