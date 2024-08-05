package converters

import (
	"github.com/maciejas22/conference-manager/api/db/repositories"
	"github.com/maciejas22/conference-manager/api/internal/models"
)

func ConvertSectionRepoToSchema(s *repositories.Section) *models.Section {
	return &models.Section{
		ID:      s.Id,
		Title:   &s.Title,
		Content: s.Content,
	}
}
