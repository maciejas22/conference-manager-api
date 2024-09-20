package services

import (
	"context"

	"github.com/maciejas22/conference-manager/api/db"
	"github.com/maciejas22/conference-manager/api/db/repositories"
	"github.com/maciejas22/conference-manager/api/internal/converters"
	"github.com/maciejas22/conference-manager/api/internal/models"
)

func UpdateUser(ctx context.Context, dbClient *db.DB, userId int, updateUserInput models.UpdateUserInput) (*models.User, error) {
	user, err := repositories.UpdateUser(dbClient.QueryExecutor, userId, updateUserInput.Username, updateUserInput.Email, updateUserInput.Name, updateUserInput.Surname)
	if err != nil {
		return nil, err
	}

	return converters.ConvertUserRepoToSchema(&user), nil
}

func GetUserData(ctx context.Context, dbClient *db.DB, userId int) (*models.User, error) {
	user, err := repositories.GetUserByID(dbClient.QueryExecutor, userId)
	if err != nil {
		return nil, err
	}

	return converters.ConvertUserRepoToSchema(&user), nil
}
