package repositories

import (
	"context"
	"errors"
	"log"

	"github.com/maciejas22/conference-manager/api/db"
)

type Subsection struct {
	Id        string  `json:"id" db:"id"`
	SectionId string  `json:"section_id" db:"section_id"`
	Title     string  `json:"title" db:"title"`
	Content   *string `json:"content" db:"content"`
}

func (s *Subsection) TableName() string {
	return "public.subsections"
}

type SubsectionRepository interface {
	GetSubsections(sectionId string) ([]Subsection, error)
}

type subsectionRepository struct {
	ctx context.Context
	db  *db.DB
}

func NewSubsectionRepository(ctx context.Context, db *db.DB) SubsectionRepository {
	return &subsectionRepository{
		ctx: ctx,
		db:  db,
	}
}

func (r *subsectionRepository) GetSubsections(sectionId string) ([]Subsection, error) {
	var subsections []Subsection
	s := &Subsection{}
	query := "SELECT id, section_id, title, content FROM " + s.TableName() + " WHERE section_id = $1"
	err := r.db.SqlConn.Select(
		&subsections,
		query,
		sectionId,
	)
	if err != nil {
		log.Println(err)
		return nil, errors.New("could not get subsections")
	}
	return subsections, nil
}
