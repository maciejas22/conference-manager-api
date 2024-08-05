package converter

import (
	"github.com/maciejas22/conference-manager/api/db/repositories"
	"github.com/maciejas22/conference-manager/api/internal/models"
)

func ConvertNewsRepoToSchema(n *repositories.News) *models.News {

	return &models.News{
		ID:      n.Id,
		Title:   n.Title,
		Content: n.Content,
		Date:    n.CreatedAt,
	}
}
