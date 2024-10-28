package services

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/maciejas22/conference-manager/api/internal/db"
	"github.com/maciejas22/conference-manager/api/internal/db/repositories"
)

func GetConferenceOrganizerId(ctx context.Context, dbClient *db.DB, conferenceId int) (int, error) {
	var organizerId int
	err := db.Transaction(ctx, dbClient.Conn, func(tx *sqlx.Tx) error {
		var err error
		organizerId, err = repositories.GetConferenceOrganizerId(tx, conferenceId)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return 0, err
	}

	return organizerId, nil
}
