package converter

import (
	"time"

	"github.com/maciejas22/conference-manager/api/db/repositories"
	"github.com/maciejas22/conference-manager/api/internal/models"
)

func ConvertAgendaItemRepoToSchema(a *repositories.AgendaItem) *models.AgendaItem {
	startTime, err := time.Parse(time.RFC3339, a.StartTime)
	if err != nil {
		return &models.AgendaItem{}
	}

	endTime, err := time.Parse(time.RFC3339, a.EndTime)
	if err != nil {
		return &models.AgendaItem{}
	}

	return &models.AgendaItem{
		ID:        a.Id,
		StartTime: startTime,
		EndTime:   endTime,
		Event:     a.Event,
		Speaker:   a.Speaker,
	}
}
