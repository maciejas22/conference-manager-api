package repositories

import (
	"context"

	"github.com/maciejas22/conference-manager/api/db"
)

type TermsOfService struct {
	Id              string `json:"id" db:"id"`
	UpdatedAt       string `json:"updated_at"     db:"updated_at"`
	Introduction    string `json:"introduction"  db:"introduction"`
	Acknowledgement string `json:"acknowledgement" db:"acknowledgement"`
}

func (t *TermsOfService) TableName() string {
	return "public.terms_of_service"
}

type TermsOfServiceRepository interface {
	GetTermsOfService() (TermsOfService, error)
}

type termsOfServiceRepository struct {
	ctx context.Context
	db  *db.DB
}

func NewTermsOfServiceRepository(ctx context.Context, db *db.DB) TermsOfServiceRepository {
	return &termsOfServiceRepository{
		ctx: ctx,
		db:  db,
	}
}

func (r *termsOfServiceRepository) GetTermsOfService() (TermsOfService, error) {
	var term TermsOfService
	query := "SELECT id, updated_at, introduction, acknowledgement FROM " + term.TableName() + " LIMIT 1"

	err := r.db.SqlConn.Get(
		&term,
		query,
	)
	if err != nil {
		return TermsOfService{}, err
	}

	return term, nil
}
