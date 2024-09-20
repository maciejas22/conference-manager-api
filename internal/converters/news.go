package converters

import (
	"log"
	"time"

	"github.com/maciejas22/conference-manager/api/db/repositories"
	"github.com/maciejas22/conference-manager/api/internal/models"
)

func ConvertNewsRepoToSchema(n *repositories.News) *models.News {
	parsedDate, err := time.Parse("2006-01-02 15:04:05", n.CreatedAt)
	if err != nil {
		log.Println("Error while parsing date: ", err)
		return nil
	}

	return &models.News{
		ID:      n.Id,
		Title:   n.Title,
		Content: n.Content,
		Date:    parsedDate,
	}
}
