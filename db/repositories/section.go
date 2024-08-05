package repositories

import "github.com/jmoiron/sqlx"

type Section struct {
	Id               string  `json:"id" db:"id"`
	TermsOfServiceId string  `json:"terms_of_service_id" db:"terms_of_service_id"`
	Title            string  `json:"title" db:"title"`
	Content          *string `json:"content" db:"content"`
}

func (s *Section) TableName() string {
	return "public.sections"
}

func GetToSSections(tx *sqlx.Tx, termsOfServiceId string) ([]Section, error) {
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
