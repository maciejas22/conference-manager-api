package services

import (
	"context"

	"github.com/maciejas22/conference-manager/api/db"
	"github.com/maciejas22/conference-manager/api/db/repositories"
	"github.com/maciejas22/conference-manager/api/internal/converters"
	"github.com/maciejas22/conference-manager/api/internal/models"
)

func GetConferenceAgenda(ctx context.Context, db *db.DB, conferenceId string) ([]*models.AgendaItem, error) {
	tx, err := db.Conn.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}

	agenda, err := repositories.GetAgenda(tx, conferenceId)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	var agendaItems []*models.AgendaItem
	for _, a := range agenda {
		agendaItems = append(agendaItems, converters.ConvertAgendaItemRepoToSchema(&a))
	}

	return agendaItems, nil
}

func GetAgendaItemsCount(ctx context.Context, db *db.DB, conferenceId string) (int, error) {
	tx, err := db.Conn.BeginTxx(ctx, nil)
	if err != nil {
		return 0, err
	}

	agendaItemsCount, err := repositories.CountAgendaItems(tx, conferenceId)
	if err != nil {
		return 0, err
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return agendaItemsCount, nil
}
