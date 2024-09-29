package services

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/maciejas22/conference-manager/api/db"
	"github.com/maciejas22/conference-manager/api/db/repositories"
	"github.com/maciejas22/conference-manager/api/internal/models"
)

func GetUserData(ctx context.Context, dbClient *db.DB, userId int) (*models.User, error) {
	var u repositories.User
	err := db.Transaction(ctx, dbClient.Conn, func(tx *sqlx.Tx) error {
		var err error
		u, err = repositories.GetUserByID(tx, userId)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &models.User{
		ID:       u.Id,
		Name:     u.Name,
		Surname:  u.Surname,
		Username: u.Username,
		Email:    *u.Email,
		Role:     models.Role(u.Role),
	}, nil
}
