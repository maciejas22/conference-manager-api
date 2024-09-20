package services

import (
	"errors"

	"github.com/maciejas22/conference-manager/api/db"
	"github.com/maciejas22/conference-manager/api/db/repositories"
	"github.com/maciejas22/conference-manager/api/internal/auth"
	"github.com/maciejas22/conference-manager/api/internal/models"
	"golang.org/x/net/context"
)

func RegisterUser(ctx context.Context, dbClient *db.DB, userData models.RegisterUserInput) (*string, error) {
	var userSession *string
	err := db.Transaction(ctx, dbClient.QueryExecutor, func(qe *db.QueryExecutor) error {
		hashedPassword, err := auth.HashPassword(userData.Password)
		if err != nil {
			return err
		}

		userId, err := repositories.CreateUser(qe, userData.Email, hashedPassword, userData.Role.String())
		if err != nil {
			return err
		}

		sessionId, err := auth.GenerateSessionId()
		if err != nil {
			return nil
		}

		userSession, err = repositories.CreateSession(qe, sessionId, userId)
		if err != nil {
			return err
		}

		return nil
	})

	return userSession, err
}

func LoginUser(ctx context.Context, dbClient *db.DB, userData models.LoginUserInput) (*string, error) {
	var userSession *string
	err := db.Transaction(ctx, dbClient.QueryExecutor, func(qe *db.QueryExecutor) error {
		user, err := repositories.GetUserByEmail(qe, userData.Email)
		if err != nil {
			return err
		}

		if !auth.CheckPasswordHash(userData.Password, user.Password) {
			return errors.New("invalid password")
		}

		sessionId, err := auth.GenerateSessionId()
		if err != nil {
			return errors.New("could not generate session id")
		}

		userSession, err = repositories.CreateSession(qe, sessionId, user.Id)
		if err != nil {
			return err
		}

		return nil
	})

	return userSession, err
}
