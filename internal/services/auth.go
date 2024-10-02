package services

import (
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/maciejas22/conference-manager/api/internal/auth"
	"github.com/maciejas22/conference-manager/api/internal/db"
	"github.com/maciejas22/conference-manager/api/internal/db/repositories"
	"github.com/maciejas22/conference-manager/api/internal/models"
	"golang.org/x/net/context"
)

func RegisterUser(ctx context.Context, dbClient *db.DB, userData models.RegisterUserInput) (*string, error) {
	var userSession *string
	hashedPassword, err := auth.HashPassword(userData.Password)
	if err != nil {
		return nil, err
	}

	err = db.Transaction(ctx, dbClient.Conn, func(tx *sqlx.Tx) error {
		var err error
		userId, err := repositories.CreateUser(tx, userData.Email, hashedPassword, userData.Role.String())
		if err != nil {
			return err
		}

		sessionId, err := auth.GenerateSessionId()
		if err != nil {
			return err
		}

		userSession, err = repositories.CreateSession(tx, sessionId, userId)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return userSession, err
}

func LoginUser(ctx context.Context, dbClient *db.DB, userData models.LoginUserInput) (*string, error) {
	var userSession *string
	err := db.Transaction(ctx, dbClient.Conn, func(tx *sqlx.Tx) error {
		var err error
		user, err := repositories.GetUserByEmail(tx, userData.Email)
		if err != nil {
			return errors.New("user not found")
		}

		if !auth.CheckPasswordHash(userData.Password, user.Password) {
			return errors.New("invalid password")
		}

		sessionId, err := auth.GenerateSessionId()
		if err != nil {
			return errors.New("could not generate session id")
		}

		userSession, err = repositories.CreateSession(tx, sessionId, user.Id)
		if err != nil {
			return errors.New("could not create session")
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return userSession, err
}
