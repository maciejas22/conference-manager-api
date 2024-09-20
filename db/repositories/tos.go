package repositories

import (
	"github.com/jmoiron/sqlx"
	"github.com/maciejas22/conference-manager/api/db"
)

type TermsOfService struct {
	Id              int    `json:"id" db:"id"`
	CreatedAt       string `json:"created_at"     db:"created_at"`
	UpdatedAt       string `json:"updated_at"     db:"updated_at"`
	Introduction    string `json:"introduction"  db:"introduction"`
	Acknowledgement string `json:"acknowledgement" db:"acknowledgement"`
}

func (t *TermsOfService) TableName() string {
	return "terms_of_service"
}

func GetTermsOfService(qe *db.QueryExecutor) (TermsOfService, error) {
	var term TermsOfService
	query := "SELECT id, updated_at, introduction, acknowledgement FROM " + term.TableName() + " LIMIT 1"

	err := sqlx.Get(
		qe,
		&term,
		query,
	)
	if err != nil {
		return TermsOfService{}, err
	}

	return term, nil
}
