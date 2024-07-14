package repositories

import (
	"context"

	"github.com/maciejas22/conference-manager/api/db"
)

type Section struct {
	Id               string  `json:"id" db:"id"`
	TermsOfServiceId string  `json:"terms_of_service_id" db:"terms_of_service_id"`
	Title            string  `json:"title" db:"title"`
	Content          *string `json:"content" db:"content"`
}

func (s *Section) TableName() string {
	return "public.sections"
}

type SectionRepository interface {
	GetSections(termsOfServiceId string) ([]Section, error)
}

type sectionRepository struct {
	ctx context.Context
	db  *db.DB
}

func NewSectionRepository(ctx context.Context, db *db.DB) SectionRepository {
	return &sectionRepository{
		ctx: ctx,
		db:  db,
	}
}

func (r *sectionRepository) GetSections(termsOfServiceId string) ([]Section, error) {
	var sections []Section
	s := &Section{}
	query := "SELECT id, terms_of_service_id, title, content FROM " + s.TableName() + " WHERE terms_of_service_id = $1"
	err := r.db.SqlConn.Select(
		&sections,
		query,
		termsOfServiceId,
	)
	if err != nil {
		return nil, err
	}
	return sections, nil
}
