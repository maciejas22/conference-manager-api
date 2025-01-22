package repository

import (
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
	table "github.com/maciejas22/conference-manager-api/cm-info/internal/repository"
	s "github.com/maciejas22/conference-manager-api/cm-info/internal/repository/section"
	ss "github.com/maciejas22/conference-manager-api/cm-info/internal/repository/subsection"
)

type ToSRepo struct {
	Db *sqlx.DB
}

func NewToSRepo(db *sqlx.DB) ToSRepoInterface {
	return &ToSRepo{Db: db}
}

func (r *ToSRepo) GetTermsOfService() (*TermsOfService, error) {
	query := `
    SELECT
        tos.id,
        tos.created_at,
        tos.updated_at,
        tos.introduction,
        tos.acknowledgement,
        s.id,
        s.title,
        s.content,
        s.created_at,
        s.updated_at,
        ss.id,
        ss.title,
        ss.content,
        ss.created_at,
        ss.updated_at
    FROM ` + table.GetTableName(table.TermsOfService) + ` tos
    LEFT JOIN ` + table.GetTableName(table.Section) + ` s ON tos.id = s.terms_of_service_id
    LEFT JOIN ` + table.GetTableName(table.Subsection) + ` ss ON s.id = ss.section_id
    ORDER BY
        s.id, ss.id
  `

	rows, err := r.Db.Queryx(query)
	if err != nil {
		return nil, errors.New("Could not get terms of service")
	}
	defer rows.Close()

	var tos *TermsOfService
	sectionsMap := make(map[int]*s.Section)

	for rows.Next() {
		var tosID, sectionID, subsectionID int
		var updatedAt, createdAt, sectionUpdatedAt, sectionCreatedAt, subsectionUpdatedAt, subsectionCreatedAt time.Time
		var introduction, acknowledgement, sectionTitle, subsectionTitle string
		var sectionContent, subsectionContent *string

		err := rows.Scan(&tosID, &updatedAt, &createdAt, &introduction, &acknowledgement, &sectionID, &sectionTitle, &sectionContent, &sectionCreatedAt, &sectionUpdatedAt, &subsectionID, &subsectionTitle, &subsectionContent, &subsectionCreatedAt, &subsectionUpdatedAt)
		if err != nil {
			return nil, errors.New("Could not scan terms of service")
		}

		if tos == nil {
			tos = &TermsOfService{
				Id:              tosID,
				CreatedAt:       createdAt,
				UpdatedAt:       updatedAt,
				Introduction:    introduction,
				Acknowledgement: acknowledgement,
				Sections:        []*s.Section{},
			}
		}

		section, exists := sectionsMap[sectionID]
		if !exists {
			section = &s.Section{
				Id:          sectionID,
				Title:       sectionTitle,
				Content:     sectionContent,
				CreatedAt:   sectionCreatedAt,
				UpdatedAt:   sectionUpdatedAt,
				Subsections: []*ss.Subsection{},
			}
			sectionsMap[sectionID] = section
			tos.Sections = append(tos.Sections, section)
		}

		subsection := &ss.Subsection{
			Id:        subsectionID,
			Title:     subsectionTitle,
			Content:   subsectionContent,
			CreatedAt: subsectionCreatedAt,
			UpdatedAt: subsectionUpdatedAt,
		}
		section.Subsections = append(section.Subsections, subsection)
	}

	return tos, nil

}
