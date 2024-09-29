package services

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/maciejas22/conference-manager/api/db"
	"github.com/maciejas22/conference-manager/api/db/repositories"
	"github.com/maciejas22/conference-manager/api/internal/models"
)

func GetConferenceAgenda(ctx context.Context, dbClient *db.DB, conferenceId int) ([]*models.AgendaItem, error) {
	var agenda []repositories.AgendaItem
	err := db.Transaction(ctx, dbClient.Conn, func(tx *sqlx.Tx) error {
		var err error
		agenda, err = repositories.GetAgenda(tx, conferenceId)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return []*models.AgendaItem{}, err
	}

	var agendaItems []*models.AgendaItem
	for _, a := range agenda {
		startTime, err := time.Parse(time.RFC3339, a.StartTime)
		if err != nil {
			return []*models.AgendaItem{}, err
		}

		endTime, err := time.Parse(time.RFC3339, a.EndTime)
		if err != nil {
			return []*models.AgendaItem{}, err
		}

		agendaItems = append(agendaItems, &models.AgendaItem{
			ID:        a.Id,
			StartTime: startTime,
			EndTime:   endTime,
			Event:     a.Event,
			Speaker:   a.Speaker,
		})
	}

	return agendaItems, err
}

func GetAgendaItemsCount(ctx context.Context, dbClient *db.DB, conferenceId int) (int, error) {
	var agendaItemsCount int
	err := db.Transaction(ctx, dbClient.Conn, func(tx *sqlx.Tx) error {
		var err error
		agendaItemsCount, err = repositories.CountAgendaItems(tx, conferenceId)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return 0, err
	}

	return agendaItemsCount, nil
}
