package converters

import (
	"time"

	"github.com/maciejas22/conference-manager/api/db/repositories"
	"github.com/maciejas22/conference-manager/api/internal/models"
)

func ConvertTosRepoToSchema(t *repositories.TermsOfService) *models.TermsOfService {
	updatedAt, err := time.Parse(time.RFC3339, t.UpdatedAt)
	if err != nil {
		return &models.TermsOfService{}
	}

	return &models.TermsOfService{
		ID:              t.Id,
		UpdatedAt:       updatedAt,
		Introduction:    t.Introduction,
		Acknowledgement: t.Acknowledgement,
	}
}
