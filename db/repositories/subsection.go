package repositories

import (
	"errors"
	"log"

	"github.com/jmoiron/sqlx"
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

func GetToSSubsections(tx *sqlx.Tx, sectionId string) ([]Subsection, error) {
	var subsections []Subsection
	s := &Subsection{}
	query := "SELECT id, section_id, title, content FROM " + s.TableName() + " WHERE section_id = $1"
	err := tx.Select(
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
