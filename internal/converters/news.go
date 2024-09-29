package converters

import (
	"time"

	"github.com/maciejas22/conference-manager/api/db/repositories"
	"github.com/maciejas22/conference-manager/api/internal/models"
)

func ConvertNewsRepoToSchema(n *repositories.News) *models.News {
	parsedDate, err := time.Parse(time.RFC3339, n.CreatedAt)
	if err != nil {
		return nil
	}

	return &models.News{
		ID:      n.Id,
		Title:   n.Title,
		Content: n.Content,
		Date:    parsedDate,
	}
}
