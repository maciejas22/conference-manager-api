package repositories

import (
	"errors"

	"github.com/jmoiron/sqlx"
)

type TermsOfService struct {
	Id              int    `json:"id" db:"id"`
	CreatedAt       string `json:"created_at"     db:"created_at"`
	UpdatedAt       string `json:"updated_at"     db:"updated_at"`
	Introduction    string `json:"introduction"  db:"introduction"`
	Acknowledgement string `json:"acknowledgement" db:"acknowledgement"`
	Sections        []*Section
}

func (t *TermsOfService) TableName() string {
	return "public.terms_of_service"
}

func GetTermsOfService(tx *sqlx.Tx) (*TermsOfService, error) {
	query := `
    SELECT
        tos.id,
        tos.updated_at,
        tos.introduction,
        tos.acknowledgement,
        s.id,
        s.title,
        s.content,
        ss.id,
        ss.title,
        ss.content
    FROM ` + (new(TermsOfService)).TableName() + ` tos
    LEFT JOIN ` + (new(Section)).TableName() + ` s ON tos.id = s.terms_of_service_id
    LEFT JOIN ` + (new(Subsection)).TableName() + ` ss ON s.id = ss.section_id
    ORDER BY
        s.id, ss.id
  `

	rows, err := tx.Queryx(query)
	if err != nil {
		return nil, errors.New("error querying database: " + err.Error())
	}
	defer rows.Close()

	var tos *TermsOfService
	sectionsMap := make(map[int]*Section)

	for rows.Next() {
		var tosID, sectionID, subsectionID int
		var updatedAt, introduction, acknowledgement, sectionTitle, subsectionTitle string
		var sectionContent, subsectionContent *string

		err := rows.Scan(&tosID, &updatedAt, &introduction, &acknowledgement, &sectionID, &sectionTitle, &sectionContent, &subsectionID, &subsectionTitle, &subsectionContent)
		if err != nil {
			return nil, errors.New("error scanning row: " + err.Error())
		}

		if tos == nil {
			tos = &TermsOfService{
				Id:              tosID,
				UpdatedAt:       updatedAt,
				Introduction:    introduction,
				Acknowledgement: acknowledgement,
				Sections:        []*Section{},
			}
		}

		section, exists := sectionsMap[sectionID]
		if !exists {
			section = &Section{
				Id:          sectionID,
				Title:       sectionTitle,
				Content:     sectionContent,
				Subsections: []*Subsection{},
			}
			sectionsMap[sectionID] = section
			tos.Sections = append(tos.Sections, section)
		}

		subsection := &Subsection{
			Id:      subsectionID,
			Title:   subsectionTitle,
			Content: subsectionContent,
		}
		section.Subsections = append(section.Subsections, subsection)
	}

	return tos, nil

}
