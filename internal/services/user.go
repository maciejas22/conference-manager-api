package services

import (
	"context"

	"github.com/maciejas22/conference-manager/api/db"
	"github.com/maciejas22/conference-manager/api/db/repositories"
	"github.com/maciejas22/conference-manager/api/internal/converter"
	"github.com/maciejas22/conference-manager/api/internal/models"
)

func UpdateUser(ctx context.Context, db *db.DB, userId string, updateUserInput models.UpdateUserInput) (*models.User, error) {
	tx, err := db.Conn.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}

	user, err := repositories.UpdateUser(tx, userId, updateUserInput.Username, updateUserInput.Email, updateUserInput.Name, updateUserInput.Surname)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return converter.ConvertUserRepoToSchema(&user), nil
}

func GetUserData(ctx context.Context, db *db.DB, userId string) (*models.User, error) {
	tx, err := db.Conn.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}

	user, err := repositories.GetUserByID(tx, userId)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return converter.ConvertUserRepoToSchema(&user), nil
}
