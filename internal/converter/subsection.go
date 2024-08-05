package converter

import (
	"github.com/maciejas22/conference-manager/api/db/repositories"
	"github.com/maciejas22/conference-manager/api/internal/models"
)

func ConvertSubsectionRepoToSchema(s *repositories.Subsection) *models.SubSection {
	return &models.SubSection{
		ID:      s.Id,
		Title:   &s.Title,
		Content: s.Content,
	}
}
