package services

import (
	"context"

	"github.com/maciejas22/conference-manager/api/db"
	"github.com/maciejas22/conference-manager/api/db/repositories"
	"github.com/maciejas22/conference-manager/api/internal/converters"
	"github.com/maciejas22/conference-manager/api/internal/models"
)

func GetConferenceAgenda(ctx context.Context, dbClient *db.DB, conferenceId int) ([]*models.AgendaItem, error) {
	agenda, err := repositories.GetAgenda(dbClient.QueryExecutor, conferenceId)

	var agendaItems []*models.AgendaItem
	for _, a := range agenda {
		agendaItems = append(agendaItems, converters.ConvertAgendaItemRepoToSchema(&a))
	}

	return agendaItems, err
}

func GetAgendaItemsCount(ctx context.Context, dbClient *db.DB, conferenceId int) (int, error) {
	agendaItemsCount, err := repositories.CountAgendaItems(dbClient.QueryExecutor, conferenceId)
	if err != nil {
		return 0, err
	}

	return agendaItemsCount, nil
}
