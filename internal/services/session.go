package services

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/maciejas22/conference-manager/api/internal/db"
	"github.com/maciejas22/conference-manager/api/internal/db/repositories"
)

func DestroyUserSession(ctx context.Context, dbClient *db.DB, userId int) (bool, error) {
	err := db.Transaction(ctx, dbClient.Conn, func(tx *sqlx.Tx) error {
		err := repositories.DeleteUserSession(tx, userId)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return false, err
	}

	return true, nil
}
