package repositories

import (
	"errors"

	"github.com/jmoiron/sqlx"
)

type Subsection struct {
	Id        int     `json:"id" db:"id"`
	SectionId int     `json:"section_id" db:"section_id"`
	Title     string  `json:"title" db:"title"`
	Content   *string `json:"content" db:"content"`
	CreatedAt string  `json:"created_at" db:"created_at"`
	JoinedAt  string  `json:"joined_at" db:"joined_at"`
}

func (s *Subsection) TableName() string {
	return "subsections"
}

func GetToSSubsections(tx *sqlx.Tx, sectionId int) ([]Subsection, error) {
	var subsections []Subsection
	s := &Subsection{}
	query := "SELECT id, section_id, title, content FROM " + s.TableName() + " WHERE section_id = $1"
	err := tx.Select(
		&subsections,
		query,
		sectionId,
	)
	if err != nil {
		return nil, errors.New("could not get subsections")
	}
	return subsections, nil
}
